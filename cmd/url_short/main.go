package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hisnameisivan/demo_url_short/internal/config"
	"github.com/hisnameisivan/demo_url_short/internal/http-server/handlers/url/save"
	mvLogger "github.com/hisnameisivan/demo_url_short/internal/http-server/middleware/logger"
	"github.com/hisnameisivan/demo_url_short/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("start service", slog.String("env", cfg.Env))
	log.Debug("debug mode enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.String("err", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	// router.Use(middleware.Logger)
	router.Use(mvLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	// router.Get("{alias}", redirect.New(log, storage))

	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:      router,
		ReadTimeout:  cfg.RWTimeout,
		WriteTimeout: cfg.RWTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info(fmt.Sprintf("server starting at %s", srv.Addr))
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server error", slog.String("err", err.Error()))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
		))
	case envProd:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo},
		))
	}

	return log
}
