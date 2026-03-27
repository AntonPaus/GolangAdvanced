package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AntonPaus/GolangAdvanced/internal/compression"
	"github.com/AntonPaus/GolangAdvanced/internal/config"
	"github.com/AntonPaus/GolangAdvanced/internal/handler"
	"github.com/AntonPaus/GolangAdvanced/internal/interfaces"
	"github.com/AntonPaus/GolangAdvanced/internal/logger"
	"github.com/AntonPaus/GolangAdvanced/internal/storage/db"
	"github.com/AntonPaus/GolangAdvanced/internal/storage/file"
	"github.com/AntonPaus/GolangAdvanced/internal/storage/memory"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	Config   *config.Config
	Router   *chi.Mux
	Handlers handler.Handler
}

func NewApp() (*App, error) {
	// Initialize config
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot initiate config: %w", err)
	}

	// Initialize logger
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}
	logger.Log.Info("Logger loaded")

	// Initialize storage
	var storage interfaces.Storage
	if cfg.DatabaseDsn != "" {
		logger.Log.Info("Database DSN: ", zap.String("dsn", cfg.DatabaseDsn))
		storage, err = db.NewDBStorage(cfg.DatabaseDsn)
		if err != nil {
			logger.Log.Panic("cannot initiate database", zap.Error(err))
			return nil, fmt.Errorf("cannot initiate database: %w", err)
		}
		logger.Log.Info("Storage DB initialized")
	} else if cfg.FileStoragePath != "" {
		storage, err = file.NewFileStorage(cfg.Restore, cfg.FileStoragePath, cfg.StoreInterval)
		if err != nil {
			logger.Log.Panic("cannot initiate file storage", zap.Error(err))
			return nil, fmt.Errorf("cannot initiate file storage: %w", err)
		}
		logger.Log.Info("Storage File initialized")
	} else {
		storage, err = memory.NewMemoryStorage()
		if err != nil {
			return nil, fmt.Errorf("cannot initiate memory storage: %w", err)
		}
		logger.Log.Info("Storage Memory initialized")
	}

	// Initialize app
	app := &App{
		Config: cfg,
		Router: chi.NewRouter(),
		Handlers: handler.Handler{
			Storage: storage,
		},
	}
	app.setupRoutes()
	logger.Log.Info("Routes initialized")
	return app, nil
}

func (a *App) setupRoutes() {
	a.Router.Use(logger.RequestLogger)
	a.Router.Use(compression.UncompressHandler)
	a.Router.Get("/", a.Handlers.MainPage)
	a.Router.Get("/ping", a.Handlers.Ping)
	a.Router.Route("/update", func(r chi.Router) {
		r.Post("/", a.Handlers.UpdateMetricJSON)
		r.Post("/{type}/{name}/{value}", a.Handlers.UpdateMetric)
	})
	a.Router.Route("/value", func(r chi.Router) {
		r.Post("/", a.Handlers.GetMetricJSON)
		r.Get("/{type}/{name}", a.Handlers.GetMetric)
	})
}

func (a *App) Run() {
	logger.Log.Info("Starting server", zap.String("address", a.Config.Address))
	http.ListenAndServe(a.Config.Address, a.Router)
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed run app: %s", err)
	}
	defer app.Handlers.Storage.Close()
	app.Run()
}
