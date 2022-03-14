package gem

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w *statusWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

var DefaultWriter io.Writer = os.Stdout

func Logger() HandlerFunc {
	return logger(LoggerOption{})
}

type LoggerOption struct {
	OutPut io.Writer
}

func logger(option LoggerOption) HandlerFunc {
	out := option.OutPut
	if out == nil {
		out = DefaultWriter
	}
	return func(c *Context) {
		starTime := time.Now()
		params, _ := json.Marshal(c.Params)
		defer func() {
			logMsg := fmt.Sprintf("[GEM] %s [status=%d] duration=%s ip=%s method=%s path=%v params=%s\n",
				starTime.Format("2006-01-02 15:04:05"),
				c.StatusCode,
				fmt.Sprintf("%v", time.Since(starTime)),
				c.ClientIP(),
				c.Method,
				c.Path,
				string(params),
			)
			fmt.Fprint(out, logMsg)
		}()
		c.Next()
	}
}
