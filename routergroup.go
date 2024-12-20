package wfw

import "net/http"

type IRoutes interface {
	Use(...HandlerFunc) IRoutes

	Handle(string, string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
}

type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

func (group *RouterGroup) Group(relativePath string, middleware ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(middleware),
		basePath: group.calculateAbsolutePath(relativePath),
		engine:   group.engine,
	}
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
}

func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(httpMethod, relativePath, handlers)
}

func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPost, relativePath, handlers)
}

func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodGet, relativePath, handlers)
}

func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPut, relativePath, handlers)
}

func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodDelete, relativePath, handlers)
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return group.basePath + relativePath
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}
