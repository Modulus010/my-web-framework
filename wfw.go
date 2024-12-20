package wfw

import "net/http"

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type Engine struct {
	RouterGroup
	routes methodRoutes
}

func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			root: true,
		},
	}
	engine.RouterGroup.engine = engine
	return engine
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	route := engine.routes.get(method)
	if route == nil {
		engine.routes = append(engine.routes, methodRoute{method: method})
		route = &engine.routes[len(engine.routes)-1]
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
	route := engine.routes.get(httpMethod)

	c.handlers, c.Params = route.getValue(rPath)
	if c.handlers != nil {
		c.Next()
		return
	}
	http.NotFound(c.Writer, c.Request)
}
