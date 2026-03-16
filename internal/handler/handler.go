package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AntonPaus/GolangAdvanced/internal/interfaces"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Storage interfaces.MetricsStorage
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strings.Join(h.Storage.GetAll(), "\n")))
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	mType, mName := chi.URLParam(r, "type"), chi.URLParam(r, "name")
	v, err := h.Storage.Get(mType, mName)
	if err != nil {
		http.Error(w, "error getting metric", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v", v)))
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	mType, mName, mValue := chi.URLParam(r, "type"), chi.URLParam(r, "name"), chi.URLParam(r, "value")
	if err := r.ParseForm(); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	switch mType {
	case interfaces.MetricTypeGauge:
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			http.Error(w, "invalid gauge value", http.StatusBadRequest)
			return
		}
		if err := h.Storage.Set(mType, mName, value); err != nil {
			http.Error(w, "error setting gauge metric", http.StatusBadRequest)
			return
		}
	case interfaces.MetricTypeCounter:
		value, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			return
		}
		if err := h.Storage.Set(mType, mName, value); err != nil {
			http.Error(w, "failed to save counter metric", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "metric type not found", http.StatusBadRequest)
		return
	}
	v, err := h.Storage.Get(mType, mName)
	if err != nil {
		http.Error(w, "error getting metric", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("metric updated: %v\n", v)))
}
