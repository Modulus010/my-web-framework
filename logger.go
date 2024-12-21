package wfw

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		if raw != "" {
			path = path + "?" + raw
		}
		log.Printf("[%d] %s in %v", c.StatusCode, path, latency)
	}

}
