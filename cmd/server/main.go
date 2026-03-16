package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AntonPaus/GolangAdvanced/internal/config"
	"github.com/AntonPaus/GolangAdvanced/internal/handler"
	"github.com/AntonPaus/GolangAdvanced/internal/storage/memory"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Router   *chi.Mux
	Handlers handler.Handler
	// Logger          *log.Logger
}

func NewApp(cfg *config.Config) (*App, error) {
	storage, err := memory.NewMemoryStorage()
	if err != nil {
		return nil, fmt.Errorf("cannot initiate storage: %w", err)
	}
	app := &App{
		Router: chi.NewRouter(),
		Handlers: handler.Handler{
			Storage: storage,
		},
	}
	app.setupRoutes()
	return app, nil
}

func (a *App) setupRoutes() {
	a.Router.Get("/", a.Handlers.MainPage)
	a.Router.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", a.Handlers.UpdateMetric)
	})
	a.Router.Get("/value/{type}/{name}", a.Handlers.GetMetric)
}

func (a *App) Run() {
	// fmt.Println("Current config:")
	// fmt.Println("\tRestore values:", a.Config.Restore)
	// fmt.Println("\tServer address:", a.Config.Address)
	// fmt.Println("\tFile storage path:", a.Config.FileStoragePath)
	// fmt.Println("\tStore interval:", a.Config.StoreInterval)
	// fmt.Printf("\nStarting server on %s", a.Config.Address)
	http.ListenAndServe(":8080", a.Router)
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed load config: %s", err)
	}
	app, err := NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed run app: %s", err)
	}
	app.Run()
}
