package wfw

import (
	"html/template"
	"net/http"
)

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

type Engine struct {
	RouterGroup

	HTMLTemplate *template.Template
	FuncMap template.FuncMap
	routes  methodRoutes
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

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine	
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.HTMLTemplate = template.Must(template.New("").Funcs(engine.FuncMap).ParseGlob(pattern))
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.FuncMap = funcMap
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
	c.engine = engine
	engine.handleHTTPRequest(c)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	route := engine.routes.get(httpMethod)
	if route != nil {
		c.handlers, c.Params = route.getValue(rPath)
		if c.handlers != nil {
			c.Next()
			return
		}
	}
	http.NotFound(c.Writer, c.Request)
}
