package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AntonPaus/GolangAdvanced/internal/config"
	"github.com/AntonPaus/GolangAdvanced/internal/handler"
	"github.com/AntonPaus/GolangAdvanced/internal/logger"
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
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot initiate config: %w", err)
	}
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}
	logger.Log.Info("Logger loaded")

	storage, err := memory.NewMemoryStorage()
	if err != nil {
		return nil, fmt.Errorf("cannot initiate storage: %w", err)
	}
	logger.Log.Info("Storage initialized")

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
	a.Router.Get("/", a.Handlers.MainPage)
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
	// fmt.Println("\tRestore values:", a.Config.Restore)
	// fmt.Println("\tFile storage path:", a.Config.FileStoragePath)
	// fmt.Println("\tStore interval:", a.Config.StoreInterval)
	// fmt.Printf("\nStarting server on %s", a.Config.Address)
	logger.Log.Info("Starting server", zap.String("address", a.Config.Address))
	http.ListenAndServe(a.Config.Address, a.Router)
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed run app: %s", err)
	}
	app.Run()
}
