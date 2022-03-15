package gem

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	buffer := new(bytes.Buffer)
	r := New()
	r.Use(RecoveryWithWriter(RecoveryOption{
		OutPut: buffer,
	}))
	r.GET("/panic", func(c *Context) {
		arr := []int{1}
		c.String(200, fmt.Sprintf("%v", arr[2]))
	})

	w := PerformRequest(r, "GET", "/panic", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
