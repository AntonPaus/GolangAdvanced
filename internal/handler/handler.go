package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AntonPaus/GolangAdvanced/internal/compression"
	"github.com/AntonPaus/GolangAdvanced/internal/logger"
	"github.com/AntonPaus/GolangAdvanced/internal/model"
	"github.com/AntonPaus/GolangAdvanced/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	Storage storage.Storage
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	body := []byte(strings.Join(h.Storage.GetAll(ctx), "\n"))
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
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	var m storage.Metrics
	fmt.Println("r.Header.Get(\"Content-Type\")", r.Header.Get("Content-Type"))
	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("m", m)
		v, err := h.Storage.Get(ctx, m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		switch m.MType {
		case model.MetricTypeGauge:
			g, ok := v.(float64)
			if !ok {
				http.Error(w, "invalid gauge value", http.StatusInternalServerError)
				return
			}
			m.Value = &g
		case model.MetricTypeCounter:
			c, ok := v.(int64)
			if !ok {
				http.Error(w, "invalid counter value", http.StatusInternalServerError)
				return
			}
			m.Delta = &c
		default:
			http.Error(w, "metric type not found", http.StatusBadRequest)
			return
		}
		body, err := json.Marshal(m)
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
	default:
		m.MType, m.ID = chi.URLParam(r, "type"), chi.URLParam(r, "name")
		v, err := h.Storage.Get(ctx, m)
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
}

// func (h *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
// 	defer cancel()
// 	var metrics []storage.Metrics = make([]storage.Metrics, 1)
// 	if err := json.NewDecoder(r.Body).Decode(&metrics[0]); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	v, err := h.Storage.Get(ctx, metrics)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	switch metrics[0].MType {
// 	case model.MetricTypeGauge:
// 		g, ok := v.(float64)
// 		if !ok {
// 			http.Error(w, "invalid gauge value", http.StatusInternalServerError)
// 			return
// 		}
// 		metrics[0].Value = &g
// 	case model.MetricTypeCounter:
// 		c, ok := v.(int64)
// 		if !ok {
// 			http.Error(w, "invalid counter value", http.StatusInternalServerError)
// 			return
// 		}
// 		metrics[0].Delta = &c
// 	default:
// 		http.Error(w, "metric type not found", http.StatusBadRequest)
// 		return
// 	}
// 	body, err := json.Marshal(metrics)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if r.Header.Get("Accept-Encoding") == "gzip" {
// 		body, err = compression.CompressGzip(body)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Encoding", "gzip")
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(body)
// }

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	var metrics storage.Metrics
	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			logger.Log.Error("error decoding json", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		if err := r.ParseForm(); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		mType, mName, mValue := chi.URLParam(r, "type"), chi.URLParam(r, "name"), chi.URLParam(r, "value")
		metrics.MType = mType
		metrics.ID = mName
		switch mType {
		case model.MetricTypeGauge:
			value, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				logger.Log.Error("invalid gauge value", zap.Error(err))
				http.Error(w, "invalid gauge value", http.StatusBadRequest)
				return
			}
			metrics.Value = &value
		case model.MetricTypeCounter:
			delta, err := strconv.ParseInt(mValue, 10, 64)
			if err != nil {
				logger.Log.Error("invalid counter value", zap.Error(err))
				http.Error(w, "invalid counter value", http.StatusBadRequest)
				return
			}
			metrics.Delta = &delta
		default:
			http.Error(w, "metric type not found", http.StatusBadRequest)
			return
		}
	}
	if err := h.Storage.Set(ctx, []storage.Metrics{metrics}); err != nil {
		logger.Log.Error("error setting metric", zap.Error(err))
		http.Error(w, "error setting metric", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
