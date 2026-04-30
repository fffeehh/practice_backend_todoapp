package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/fffeehh/practice_backend_todoapp/internal/core/logger"
	core_http_response "github.com/fffeehh/practice_backend_todoapp/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Зачастую в go middleware - это просто какая-то функция, которая оборачивает обработку http запросов. То есть это обертка, которая расширяет обработчик, без нужды писать в него лишний повторяющийся код.

// данный middleware производит работу с requestID. Проверяет, что он есть, а если нет, то задает и вставляет в хеддер
func RequestID() Middleware {
	const requestIDHeader = "X-Request-ID"

	// возвращаем хэндлер, как и описали выше (middleware - функция, которая принимает и возвращает хэндлер)
	// при описании middleware принято называть обработчик не h, как обычно, а next
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// задача данного middleware - во входящие http запросы зашить логгер.
// На момент попадания в этот middleware, уже вызовется предыдущий и будет проставлен requestID
func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")

			// преконфигурирование логгера на автоматическое логгирование определенных полей
			// метод With принимает набор полей, которые нужно автоматически логгировать при каждом вызове
			// этого логгера
			// Для правильного использования, его необходимо переопределить в файле логгера 
			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
				)
			// Теперь при логгирование информации на любом уровне автоматически будут логгироваться эти два поля
			
			// передаем наш логгер в хэндлер через контекст.
			// здесь мы создаем дочерний контекст родительского контекста хэндлера, но в этом дочернем 
			// контексте мы создаем еще и значение, в которое кладем наш логгер
			// ctx := context.WithValue(r.Context(), core_logger.LoggerContextKey, l)


			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
 

// задача данного middleware - отлавливать возможные паники
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// получаем логгер
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func(){
				if p := recover(); p != nil {
					// очень много всего необходимо сделать при панике, поэтому делаем отдельный обработчик в
					// internal/core/transport/http/response/handler.go
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
						)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}


func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", time.Now().UTC()),
				)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
				)
		})
	}
}
