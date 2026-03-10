package main

import (
	"time"

	"github.com/AntonPaus/GolangAdvanced/internal/metrics"
)

func main() {
	m := metrics.NewMetrics()
	go m.Poll(time.Duration(1) * time.Second)
	go m.Report(time.Duration(2)*time.Second, "http://localhost:8080")
	select {}
}
