package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rgurov/bookkeeper/internal/config"
	"github.com/rgurov/bookkeeper/internal/handler/httprouter"
	"github.com/rgurov/bookkeeper/internal/handler/httpserver"
	"github.com/rgurov/bookkeeper/internal/repository"
	"github.com/rgurov/bookkeeper/internal/service"
	"github.com/rgurov/bookkeeper/pkg/postgres"
	"golang.org/x/exp/slog"
)

func Run() {
	cfg := config.Read()
	logger := slog.Default()

	pqClient, err := postgres.Connect(
		context.Background(),
		&postgres.Config{
			Host:     cfg.PostgresHost,
			Port:     cfg.PostgresPort,
			User:     cfg.PostgresUser,
			Password: cfg.PostgresPassword,
			Database: cfg.PostgresDatabase,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(pqClient)
	serv := service.NewService(
		repo.User,
	)

	router := httprouter.New(
		logger,
		serv.Auth,
		cfg.JwtSecret,
	)

	go func() {
		httpserver.Start(
			logger,
			router,
			&httpserver.Config{
				Host:        cfg.Host,
				Port:        cfg.Port,
				Timeout:     cfg.Timeout,
				IdleTimeout: cfg.IdleTimeout,
			},
		)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	logger.Info("shuting down...")
}
