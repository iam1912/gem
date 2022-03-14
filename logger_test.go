package gem

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	r := New()
	r.Use(logger(LoggerOption{
		OutPut: buffer,
	}))

	r.GET("/example", func(c *Context) {
		c.SetStatusCode(200)
	})
	r.POST("/example", func(c *Context) {
		c.SetStatusCode(200)
	})
	r.PUT("/example", func(c *Context) {
		c.SetStatusCode(200)
	})
	r.DELETE("/example", func(c *Context) {
		c.SetStatusCode(200)
	})

	PerformRequest(r, "GET", "/example?name=hello", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	PerformRequest(r, "POST", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	PerformRequest(r, "PUT", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	PerformRequest(r, "DELETE", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")
}

func PerformRequest(r http.Handler, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
