package wfw

import (
	"fmt"
	"log"
	"net/http"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", message)
				c.Status(http.StatusInternalServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
