package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/AntonPaus/GolangAdvanced/internal/metrics"
)

func main() {
	address := flag.String("a", "localhost:8080", "server endpoint")
	reportInterval := flag.Int("r", 10, "reportInterval")
	pollInterval := flag.Int("p", 2, "pollInterval")
	flag.Parse()
	osEP := os.Getenv("ADDRESS")
	osRI := os.Getenv("REPORT_INTERVAL")
	osPI := os.Getenv("POLL_INTERVAL")
	if osEP != "" {
		*address = osEP
	}
	if osRI != "" {
		*reportInterval, _ = strconv.Atoi(osRI)
	}
	if osPI != "" {
		*pollInterval, _ = strconv.Atoi(osPI)
	}
	m := metrics.NewMetrics()
	go m.Poll(time.Duration(*pollInterval) * time.Second)
	go m.Report(time.Duration(*reportInterval)*time.Second, *address)
	select {}
}
