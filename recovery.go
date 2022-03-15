package gem

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var DefaultErrorWriter io.Writer = os.Stdout

type RecoveryOption struct {
	OutPut io.Writer
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return RecoveryWithWriter(RecoveryOption{
		OutPut: DefaultErrorWriter,
	})
}

func RecoveryWithWriter(option RecoveryOption) HandlerFunc {
	out := option.OutPut
	if out == nil {
		out = DefaultErrorWriter
	}
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				errorMsg := fmt.Sprintf("[Recovery] %s panic recovered: %s\n\n", time.Now().Format("2006-01-02 15:04:05"), trace(message))
				fmt.Fprint(out, errorMsg)
				c.AbortWithStatus(http.StatusBadRequest)
			}
		}()
		c.Next()
	}
}
