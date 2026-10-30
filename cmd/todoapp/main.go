package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/krxzy37/todoapp-go/internal/core/logger"
	core_pgx_pool "github.com/krxzy37/todoapp-go/internal/core/repository/postgres/conn/pgx"
	core_http_middleware "github.com/krxzy37/todoapp-go/internal/core/transport/http/middleware"
	core_http_server "github.com/krxzy37/todoapp-go/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/krxzy37/todoapp-go/internal/features/tasks/repository/postgres"
	tasks_service "github.com/krxzy37/todoapp-go/internal/features/tasks/service"
	tasks_transport_http "github.com/krxzy37/todoapp-go/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/krxzy37/todoapp-go/internal/features/users/repository/postgres"
	users_service "github.com/krxzy37/todoapp-go/internal/features/users/service"
	user_transport_http "github.com/krxzy37/todoapp-go/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {

	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("applocation time zone", zap.Any("zone", timeZone))

	logger.Debug("Initializing postgres connection pool")
	//pool, err := core_postgres_pool.NewConnectionPool(
	//	ctx,
	//	core_postgres_pool.NewConfigMust(),
	//)

	pool, err := core_pgx_pool.NewPool(ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("Initializing feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)

	userTransportHTTP := user_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("Initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTaskService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	ApiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	ApiVersionRouterV1.RegisterRoutes(userTransportHTTP.Routes()...)

	ApiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(ApiVersionRouterV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}

}
