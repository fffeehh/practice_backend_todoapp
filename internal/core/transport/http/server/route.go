package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/fffeehh/practice_backend_todoapp/internal/core/transport/http/middleware"
)

// Роут - набор параметров, благодаря которым мультиплексор сможет понять,
// как по входящим данным http запроса выбрать обработчик
/*
То есть каждая фича сама сможет описать набор роутов, по которым она будет предоставлять свою функциональность.
Затем мы где нибудь в main.go будем от каждой фичи получать ее набор роутов, регистрировать в определенном
APIVersionRouter, а он в свою очередь должен быть зарегестрирован в HTTPServer
*/
type Route struct {
	Method string
	Path string
	Handler http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}

/* пока не используется, поэтому не нужен
func NewRoute(
	method string,
	path string,
	handler http.HandlerFunc,
) Route {
	return Route{
		Method: method,
		Path: path,
		Handler: handler,
	}
}
*/
