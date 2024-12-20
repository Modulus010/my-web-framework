package wfw

import "net/http"

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	Params   Params
	handlers HandlersChain
	index    int8
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		if c.handlers[c.index] != nil {
			c.handlers[c.index](c)
		}
		c.index++
	}
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		index:   -1,
	}
}
