package main

import (
	"flag"
	"log"
	"time"

	"github.com/Ropho/avito-bootcamp-assignment/api"
	"github.com/Ropho/avito-bootcamp-assignment/config"
	timeModel "github.com/Ropho/avito-bootcamp-assignment/internal/models/time"
	"github.com/Ropho/avito-bootcamp-assignment/internal/repository/postgres"
	"github.com/Ropho/avito-bootcamp-assignment/internal/router"
	"github.com/Ropho/avito-bootcamp-assignment/internal/server"
	"github.com/Ropho/avito-bootcamp-assignment/internal/service"
	"github.com/Ropho/avito-bootcamp-assignment/internal/usecases"
	"github.com/Ropho/avito-bootcamp-assignment/pkg/jwt"
	"github.com/Ropho/avito-bootcamp-assignment/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

const configPath = "./config"

var (
	configName string
)

func init() {
	flag.StringVar(&configName, "config_name", "config", "name for config file in folder ./config")
}

func main() {
	flag.Parse()

	err := config.Init(configPath, configName)
	if err != nil {
		log.Fatalf("failed to init config: [%v]", err)
	}

	cfg, err := config.GetServiceCONF()
	if err != nil {
		log.Fatalf("failed to get service configuration: [%v]", err)
	}

	logger := logger.NewLogger(cfg.LoggerConfig)

	conn, err := postgres.OpenConnection(cfg.URL)
	if err != nil {
		logger.Fatal("failed to open postgres connection: [%v]", err)
	}

	repo := postgres.NewPgRepo(&postgres.NewPgRepoParams{
		Conn: conn,
	})

	jwtService, err := jwt.NewJWTService(&jwt.NewJWTServiceParams{
		JwtConfig: cfg.JWTServiceConf,
	})
	if err != nil {
		logger.Fatal("failed to init jwt service: [%v]", err)
	}

	usecases := usecases.NewUsecases(usecases.NewUsecasesParams{
		Repo:       &repo,
		Time:       timeModel.NewTimeImpl(time.Now()),
		JWTService: jwtService,
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

	logger.Fatal("server stopped", server.ListenAndServe())
}
