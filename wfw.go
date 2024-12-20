package wfw

import "net/http"

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type Engine struct {
	routes methodRoutes
}

func New() *Engine {
	return &Engine{}
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	route := engine.routes.get(method)
	if route == nil {
		route = &methodRoute{method: method}
		engine.routes = append(engine.routes, *route)
	}
	route.addRoute(path, handlers)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.handleHTTPRequest(c)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	for _, route := range engine.routes {
		if route.method != httpMethod {
			continue
		}
		c.handlers, c.Params = route.getValue(rPath)
		c.Next()
		return
	}
}
