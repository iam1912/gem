package gem

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Method     string
	Path       string
	Params     map[string]string
	StatusCode int
	index      int
	handlers   []HandlerFunc
	engine     *Engine
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Method: r.Method,
		Path:   r.URL.Path,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index += 1
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}

func (c *Context) AbortWithStatus(code int) {
	c.SetStatusCode(code)
	c.Abort()
}

func (c *Context) AbortStatusJSON(code int, err interface{}) {
	c.Abort()
	c.JSON(code, err)
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	value, ok := c.Params[key]
	if !ok {
		return ""
	}
	return value
}

func (c *Context) FormFile(name string, maxMemory int64) (*multipart.FileHeader, error) {
	if c.Req.MultipartForm == nil {
		if err := c.Req.ParseMultipartForm(maxMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Req.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, nil
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) SetStatusCode(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatusCode(code)
	data, err := json.Marshal(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
	c.Writer.Write([]byte(data))
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatusCode(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatusCode(code)
	if err := c.engine.htmlRender.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.AbortStatusJSON(500, err.Error())
	}
}

func (c *Context) Data(code int, contentType string, data []byte) {
	c.SetHeader("Content-Type", contentType)
	c.SetStatusCode(code)
	c.Writer.Write(data)
}

func (c *Context) ClientIP() string {
	xRealIP := c.Req.Header.Get("X-Real-Ip")
	xForwardedFor := c.Req.Header.Get("X-Forwarded-For")
	if xRealIP == "" && xForwardedFor == "" {
		ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Req.RemoteAddr))
		if err != nil {
			return ""
		} else {
			remoteIP := net.ParseIP(ip)
			return remoteIP.String()
		}
	}
	address := strings.Split(xForwardedFor, ",")
	if len(address) != 0 {
		return strings.TrimSpace(address[0])
	}
	return xRealIP
}
