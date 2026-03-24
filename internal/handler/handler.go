package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AntonPaus/GolangAdvanced/internal/compression"
	"github.com/AntonPaus/GolangAdvanced/internal/interfaces"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Storage interfaces.MetricsStorage
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	body := []byte(strings.Join(h.Storage.GetAll(), "\n"))
	err := error(nil)
	if r.Header.Get("Accept-Encoding") == "gzip" {
		body, err = compression.CompressGzip(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	mType, mName := chi.URLParam(r, "type"), chi.URLParam(r, "name")
	v, err := h.Storage.Get(mType, mName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	body := []byte(fmt.Sprintf("%v", v))
	if r.Header.Get("Accept-Encoding") == "gzip" {
		body, err = compression.CompressGzip(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (h *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metrics interfaces.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	v, err := h.Storage.Get(metrics.MType, metrics.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	switch metrics.MType {
	case interfaces.MetricTypeGauge:
		g, ok := v.(float64)
		if !ok {
			http.Error(w, "invalid gauge value", http.StatusInternalServerError)
			return
		}
		metrics.Value = &g
	case interfaces.MetricTypeCounter:
		c, ok := v.(int64)
		if !ok {
			http.Error(w, "invalid counter value", http.StatusInternalServerError)
			return
		}
		metrics.Delta = &c
	default:
		http.Error(w, "metric type not found", http.StatusBadRequest)
		return
	}
	body, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Header.Get("Accept-Encoding") == "gzip" {
		body, err = compression.CompressGzip(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (h *Handler) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metrics interfaces.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch metrics.MType {
	case interfaces.MetricTypeGauge:
		if err := h.Storage.Set(metrics.MType, metrics.ID, *metrics.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case interfaces.MetricTypeCounter:
		if err := h.Storage.Set(metrics.MType, metrics.ID, *metrics.Delta); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "metric type not found", http.StatusBadRequest)
		return
	}
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
