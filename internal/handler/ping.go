package handler

import (
	"net/http"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if h.Storage.Ping() != nil {
		http.Error(w, "failed to ping database", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
