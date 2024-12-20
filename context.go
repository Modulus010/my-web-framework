package wfw

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	Params   Params
	handlers HandlersChain
	index    int8

	Keys map[string]any

	queryCache url.Values

	formCache url.Values
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		index:   -1,
	}
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

func (c *Context) Abort() {
	c.index = int8(len(c.handlers))
}

func (c *Context) Set(key string, value any) {
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}
	c.Keys[key] = value
}

func (c *Context) Get(key string) (value any, exists bool) {
	value, exists = c.Keys[key]
	return
}

func (c *Context) Param(key string) string {
	return c.Params.ByName(key)
}

func (c *Context) AddParam(key, value string) {
	c.Params = append(c.Params, Param{Key: key, Value: value})
}

func (c *Context) Query(key string) string {
	if c.queryCache == nil {
		c.queryCache = c.Request.URL.Query()
	}
	return c.queryCache.Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

func (c *Context) HTML(code int, html string) {
	c.Header("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Header("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.Header("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
