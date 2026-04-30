package core_http_server

import "net/http"

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
}

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
