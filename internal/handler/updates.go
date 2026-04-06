package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AntonPaus/GolangAdvanced/internal/storage"
)

func (h *Handler) Updates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	var metrics []storage.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(metrics) == 0 {
		http.Error(w, "No metrics to update", http.StatusBadRequest)
		return
	}
	if err := h.Storage.Set(ctx, metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
