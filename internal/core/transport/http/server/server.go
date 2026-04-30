package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/fffeehh/practice_backend_todoapp/internal/core/logger"
	core_http_middleware "github.com/fffeehh/practice_backend_todoapp/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	// создаем мультиплексор
	// мультиплексор - сущность может по входящему http запросу распознать, через какие middleware пропустить этот запрос и в какой http обработчик направить 
	mux *http.ServeMux
	config Config
	log *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
		config: config,
		log: log,
		middleware: middleware,
	}
}

// регестрируем апи роутер
func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)	

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
			)
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	// перед регистрацией мультиплексора в сервере, навешиваем на него middlewares
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr: s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	// запускаем сервер
	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			// передаем возможные ошибки через канал ошибок
			ch <- err
		}
	}()

	select {
	// слушаем возможные ошибки
	case err := <- ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	// если завершается контекст, который был передан в функцию Run(), пробуем совершить graceful shutdown
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		// создаем независимый контекст, который отменится через время, указанное в конфиге (ShutdownTimeout)
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		// метод Shutdown() пытается корректно завершит работу сервера, если не получается за ShutdownTimeout,
		// тогда вернет ошибку. Если так вышло, насильно завершаем сервер
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
