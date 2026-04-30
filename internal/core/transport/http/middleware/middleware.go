package core_http_middleware

import "net/http"

// Создаем новый тип Middleware, под который будет подходить любая функция, принимающая и возвращающая хэндлер
// делается это просто для краткости, понятности и чистоты
type Middleware func(http.Handler) http.Handler

// функция, которая повзоляет правильно определять порядок выполнения middleware, 
// обернуть поступивший хэндлер во все middleware
func ChainMiddleware(
	h http.Handler,
	m ...Middleware,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	// идем по списку middlewares в обратном порядке, потому что если идти прямо,
	// то последовательность слоев нарушится
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

