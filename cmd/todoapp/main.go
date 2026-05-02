package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/fffeehh/practice_backend_todoapp/internal/core/config"
	core_logger "github.com/fffeehh/practice_backend_todoapp/internal/core/logger"
	core_postgres_pool "github.com/fffeehh/practice_backend_todoapp/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/fffeehh/practice_backend_todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/fffeehh/practice_backend_todoapp/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/fffeehh/practice_backend_todoapp/internal/features/statistics/repository/postgres"
	statistics_service "github.com/fffeehh/practice_backend_todoapp/internal/features/statistics/service"
	statistics_transport_http "github.com/fffeehh/practice_backend_todoapp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/fffeehh/practice_backend_todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/fffeehh/practice_backend_todoapp/internal/features/tasks/service"
	tasks_transport "github.com/fffeehh/practice_backend_todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/fffeehh/practice_backend_todoapp/internal/features/users/repository/postgres"
	users_service "github.com/fffeehh/practice_backend_todoapp/internal/features/users/service"
	users_transport_http "github.com/fffeehh/practice_backend_todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main(){
	// создаем глобальный конфиг приложения
	cfg := core_config.NewConfigMust()
	// выставляем зону из конфига (UTC по умолчанию), чтобы не хардкодить ее в код
	time.Local = cfg.TimeZone

	// создаем родительский контекст для передачи в сервер, который будет завязан на системных сигналах
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM, syscall.SIGINT,
		)
	// необязательно, тк мы и не собираемся вручную отменять контекст, но это является хорошей практикой в go
	defer cancel()

	// 	создаем логгер
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		// завершаем выполнение программы со статус кодом 1
		// статус код 1 при завершении программы означает, что во время исполнения программы возникла ошибка 
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	// создаем пул
logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
		)
	if err != nil {
		// метод fatal() логгирует сообщение на уровне fatal и автоматически вызывает os.Exit(1)
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool) // создаем репозиторий
	usersService := users_service.NewUsersService(usersRepository) // создаем сервис
	usersTransoportHTTP := users_transport_http.NewUsersHTTPHandler(usersService) //уровень транспорта для фичи users

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport.NewTasksHTTPService(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)


	logger.Debug("initializing HTTP server")

	// создаем http сервер
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	
	// создаем апи роутер первой версии
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	// регистрируем полученые из usersTransportHTTP роуты в apiVersionRouter
	apiVersionRouter.RegisterRoutes(usersTransoportHTTP.Routes()...) // передаем роуты
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)


	/* 
	usersRoutes := usersTransoportHTTP.Routes() // получаем роуты фичи users
	*/

	// готовый apiVersionRouter регестрируем в нашем http сервере
	httpServer.RegisterAPIRouters(apiVersionRouter)


	// запуск http сервера
	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}



}
