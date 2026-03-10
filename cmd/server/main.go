package main

import (
	"net/http"

	"github.com/AntonPaus/GolangAdvanced/internal/handler"
	"github.com/AntonPaus/GolangAdvanced/internal/storage/memory"
)

var h handler.Handler = handler.Handler{
	Storage: memory.NewMemoryStorage(),
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, h.MainPage)
	mux.HandleFunc(`/update/`, h.UpdateMetric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
