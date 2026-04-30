package core_http_response

import "net/http"

var (
	StatusCodeUninitialized = -1
)

// создаем кастомную структуру, которую мы можем использовать с интерфейсом там, где надо вычленить статус код
// используется в middleware логгирования, чтобы иметь возможность залоггировать статус код
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter (w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode: StatusCodeUninitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) GetStatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		return http.StatusOK
	}

	return rw.statusCode
} 
