package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

type Config struct {
	Host        string
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func Start(logger *slog.Logger, handler http.Handler, cfg *Config) {
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      handler,
	}

	logger.Info("listening at " + cfg.Host + ":" + cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server error", err)
	}
	logger.Info("server stopped")
}
