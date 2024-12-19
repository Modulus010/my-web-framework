package wfw

import "net/http"

type Context struct {
	Requese *http.Request
	Writer  http.ResponseWriter
}
