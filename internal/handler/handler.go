package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AntonPaus/GolangAdvanced/internal/repository"
)

type Handler struct {
	Storage repository.MetricsStorage
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Nothing happens. It is main page."))
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 {
		http.Error(w, "wrong path", http.StatusNotFound)
		return
	}
	metricType := parts[1]
	metricName := parts[2]
	metricValue := parts[3]

	switch metricType {
	case repository.MetricTypeGauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "invalid gauge value", http.StatusBadRequest)
			return
		}
		if err := h.Storage.Set(metricType, metricName, value); err != nil {
			http.Error(w, "error setting gauge metric", http.StatusBadRequest)
			return
		}
	case repository.MetricTypeCounter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			return
		}
		if err := h.Storage.Set(metricType, metricName, value); err != nil {
			http.Error(w, "failed to save counter metric", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "metric type not found", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	v, err := h.Storage.Get(metricType, metricName)
	if err != nil {
		http.Error(w, "error getting metric", http.StatusBadRequest)
		return
	}
	w.Write([]byte(fmt.Sprintf("metric updated: %v\n", v)))
}
