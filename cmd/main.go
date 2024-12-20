package main

import (
	wfw "github.com/Modulus010/my-web-framework"
)

func main() {
	r := wfw.New()
	r.GET("/ping", func(c *wfw.Context) {
		c.JSON(200, wfw.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}
