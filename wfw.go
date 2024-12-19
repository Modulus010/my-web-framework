package wfw

import "net/http"

type HandlerFunc func(*Context)

type Engine struct {
}

func New() *Engine {
	return &Engine{}
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
