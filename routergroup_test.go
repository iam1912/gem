package gem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterGroupBasic(t *testing.T) {
	r := New()
	v1 := r.Group("/hello")
	v1.Use(func(c *Context) {})

	//logger and recover middleware
	assert.Len(t, v1.middlewares, 1)
	assert.Equal(t, "/hello", v1.basePath)
	assert.Equal(t, r, v1.engie)

	v2 := v1.Group("/world")
	v2.Use(func(c *Context) {})

	assert.Len(t, v2.middlewares, 2)
	assert.Equal(t, "/hello/world", v2.basePath)
	assert.Equal(t, r, v2.engie)
}
