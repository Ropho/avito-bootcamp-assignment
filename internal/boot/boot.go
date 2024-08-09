package boot

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/Ropho/avito-bootcamp-assignment/api"
	"github.com/Ropho/avito-bootcamp-assignment/config"
	timeModel "github.com/Ropho/avito-bootcamp-assignment/internal/models/time"
	"github.com/Ropho/avito-bootcamp-assignment/internal/worker"

	"github.com/Ropho/avito-bootcamp-assignment/internal/repository/postgres"
	"github.com/Ropho/avito-bootcamp-assignment/internal/router"
	"github.com/Ropho/avito-bootcamp-assignment/internal/server"
	"github.com/Ropho/avito-bootcamp-assignment/internal/service"
	"github.com/Ropho/avito-bootcamp-assignment/internal/usecases"
	"github.com/Ropho/avito-bootcamp-assignment/pkg/jwt"
	"github.com/Ropho/avito-bootcamp-assignment/pkg/logger"
)

func App(configPath, configName string) error {
	err := config.Init(configPath, configName)
	if err != nil {
		return fmt.Errorf("failed to init config: [%w]", err)
	}

	cfg, err := config.GetServiceCONF()
	if err != nil {
		return fmt.Errorf("failed to get service configuration: [%w]", err)
	}

	logger := logger.NewLogger(cfg.LoggerConfig)

	conn, err := postgres.OpenConnection(cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to open postgres connection: [%w]", err)
	}

	repo := postgres.NewPgRepo(&postgres.NewPgRepoParams{
		Conn: conn,
	})

	jwtService, err := jwt.NewJWTService(&jwt.NewJWTServiceParams{
		JwtConfig: cfg.JWTServiceConf,
	})
	if err != nil {
		return fmt.Errorf("failed to init jwt service: [%w]", err)
	}

	emailsChan := make(chan []string, 1000)

	emailWorker := worker.NewWEmailorker(worker.NewWorkerParams{
		EmailsChan: emailsChan,
	})
	go func() {
		logger.Fatal("worker stopped with error: ", emailWorker.Work())
	}()

	usecases := usecases.NewUsecases(usecases.NewUsecasesParams{
		Repo:       &repo,
		Time:       timeModel.NewTimeImpl(time.Now()),
		JWTService: jwtService,
		Logger:     logger,
		EmailChan:  emailsChan,
	})

	serv := service.NewService(service.NewServiceParams{
		Usecases: &usecases,
		Logger:   logger,
	})

	mux := router.NewRouter()

	manager := router.NewInterceptorsManager(jwtService, logger)

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(manager.Authentication)

	h := api.HandlerFromMux(&serv, mux)
	server := server.NewServer(h, cfg.ServerPort)

	return server.ListenAndServe()
}
