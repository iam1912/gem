package gem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	word := ""
	r := New()
	r.Use(func(c *Context) {
		word += "HELLO"
		c.Next()
		word += "WORLD"
	})
	r.Use(func(c *Context) {
		word += "HELLO"
	})
	r.GET("/word", func(c *Context) {
		word += "GOLANG"
	})

	PerformRequest(r, "GET", "/word", nil)
	assert.Equal(t, "HELLOHELLOGOLANGWORLD", word)
}
