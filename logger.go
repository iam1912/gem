package gem

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func Logger() HandlerFunc {
	return logger()
}

func logger() HandlerFunc {
	return func(c *Context) {
		starTime := time.Now()
		params, _ := json.Marshal(c.Params)
		defer func() {
			log.Printf("[GEM] %s [status=%d] duration=%s ip=%s method=%s path=%v params=%s",
				starTime.Format("2006-01-02 15:04:05"),
				c.StatusCode,
				fmt.Sprintf("%v", time.Since(starTime)),
				c.ClientIP(),
				c.Method,
				c.Path,
				string(params),
			)
		}()
		c.Next()
	}
}
