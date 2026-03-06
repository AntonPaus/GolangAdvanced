package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type gauge float64
type counter int64
type MemStorage struct {
	metricsFloat64 map[string]gauge
	metricsInt64   map[string]counter
}

var storage = MemStorage{
	metricsFloat64: make(map[string]gauge),
	metricsInt64:   make(map[string]counter),
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Nothing happens. It is main page."))
}

func updatePage(w http.ResponseWriter, r *http.Request) {
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
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "invalid gauge value", http.StatusBadRequest)
			return
		}
		storage.metricsFloat64[metricName] = gauge(value)
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			return
		}
		storage.metricsInt64[metricName] += counter(value)
	default:
		http.Error(w, "metric type not found", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	body := ""
	for k, v := range storage.metricsFloat64 {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	for k, v := range storage.metricsInt64 {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	w.Write([]byte(body))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	mux.HandleFunc(`/update/`, updatePage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
