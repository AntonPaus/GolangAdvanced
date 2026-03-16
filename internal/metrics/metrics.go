package metrics

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type Metrics struct {
	values struct {
		Alloc         float64
		BuckHashSys   float64
		Frees         float64
		GCCPUFraction float64
		GCSys         float64
		HeapAlloc     float64
		HeapIdle      float64
		HeapInuse     float64
		HeapObjects   float64
		HeapReleased  float64
		HeapSys       float64
		LastGC        float64
		Lookups       float64
		MCacheInuse   float64
		MCacheSys     float64
		MSpanInuse    float64
		MSpanSys      float64
		Mallocs       float64
		NextGC        float64
		NumForcedGC   float64
		NumGC         float64
		OtherSys      float64
		PauseTotalNs  float64
		StackInuse    float64
		StackSys      float64
		Sys           float64
		TotalAlloc    float64
		RandomValue   float64
		PollCount     int64
	}
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Poll(interval time.Duration) {
	var mem runtime.MemStats
	for {
		time.Sleep(interval)
		runtime.ReadMemStats(&mem)
		m.values.Alloc = float64(mem.Alloc)
		m.values.BuckHashSys = float64(mem.BuckHashSys)
		m.values.Frees = float64(mem.Frees)
		m.values.GCCPUFraction = float64(mem.GCCPUFraction)
		m.values.GCSys = float64(mem.GCSys)
		m.values.HeapAlloc = float64(mem.HeapAlloc)
		m.values.HeapIdle = float64(mem.HeapIdle)
		m.values.HeapInuse = float64(mem.HeapInuse)
		m.values.HeapObjects = float64(mem.HeapObjects)
		m.values.HeapReleased = float64(mem.HeapReleased)
		m.values.HeapSys = float64(mem.HeapSys)
		m.values.LastGC = float64(mem.LastGC)
		m.values.Lookups = float64(mem.Lookups)
		m.values.MCacheInuse = float64(mem.MCacheInuse)
		m.values.MCacheSys = float64(mem.MCacheSys)
		m.values.MSpanInuse = float64(mem.MSpanInuse)
		m.values.MSpanSys = float64(mem.MSpanSys)
		m.values.Mallocs = float64(mem.Mallocs)
		m.values.NextGC = float64(mem.NextGC)
		m.values.NumForcedGC = float64(mem.NumForcedGC)
		m.values.NumGC = float64(mem.NumGC)
		m.values.OtherSys = float64(mem.OtherSys)
		m.values.PauseTotalNs = float64(mem.PauseTotalNs)
		m.values.StackInuse = float64(mem.StackInuse)
		m.values.StackSys = float64(mem.StackSys)
		m.values.Sys = float64(mem.Sys)
		m.values.TotalAlloc = float64(mem.TotalAlloc)
		m.values.PollCount = 1
		m.values.RandomValue = float64(rand.Float64())
		fmt.Printf("Poll completed\n")
	}
}

func (m *Metrics) Report(interval time.Duration, ep string) {
	var c int64
	var g float64
	for {
		errFound := false
		time.Sleep(interval)
		statsType := reflect.TypeOf(m.values)
		statsValue := reflect.ValueOf(m.values)
		for i := range statsType.NumField() {
			value := statsValue.Field(i)
			fieldName := statsType.Field(i).Name
			switch value.Kind() {
			case reflect.Int64:
				c = int64(value.Int())
				if err := sendMetricCounter(fieldName, c, ep); err != nil {
					fmt.Println("Error sending HTTP request:", err)
					errFound = true
				}
			case reflect.Float64:
				g = float64(value.Float())
				if err := sendMetricGauge(fieldName, g, ep); err != nil {
					fmt.Println("Error sending HTTP request:", err)
					errFound = true
				}
			default:
				fmt.Printf("Value type error\nSkipping...\n")
			}
		}
		if !errFound {
			fmt.Println("Report completed")
		}
	}
}

func sendMetricCounter(fieldName string, metric int64, ep string) error {
	s := fmt.Sprintf("http://%s/update/counter/%s/%d", ep, fieldName, metric)
	fmt.Println(s)
	req, err := http.NewRequest("POST", s, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func sendMetricGauge(fieldName string, metric float64, ep string) error {
	s := fmt.Sprintf("http://%s/update/gauge/%s/%f", ep, fieldName, metric)
	fmt.Println(s)
	req, err := http.NewRequest("POST", s, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// func (m *Metrics) Report(interval time.Duration, ep string) {
// 	var c int64
// 	var g float64
// 	for {
// 		errFound := false
// 		time.Sleep(interval)
// 		statsType := reflect.TypeOf(m.values)
// 		statsValue := reflect.ValueOf(m.values)
// 		for i := range statsType.NumField() {
// 			var h handlers.Metrics
// 			field := statsType.Field(i)
// 			value := statsValue.Field(i)
// 			fieldTypeParts := strings.Split(field.Type.String(), ".")
// 			fieldType := fieldTypeParts[len(fieldTypeParts)-1]
// 			h.ID, h.MType = field.Name, fieldType
// 			switch value.Kind() {
// 			case reflect.Int64:
// 				c = int64(value.Int())
// 				h.Delta = &c
// 			case reflect.Float64:
// 				g = float64(value.Float())
// 				h.Value = &g
// 			default:
// 				fmt.Printf("Value type error\nSkipping...\n")
// 			}
// 			jsonData, err := json.Marshal(h)
// 			if err != nil {
// 				fmt.Printf("JSON Marshaling error: %v\nSkipping...\n", err)
// 			}
// 			compressedData, err := compression.CompressGzip(jsonData)
// 			if err != nil {
// 				fmt.Printf("Compression error: %v\nSkipping...\n", err)
// 			}
// 			if err := sendMetric(compressedData, ep); err != nil {
// 				fmt.Println("Error sending HTTP request:", err)
// 				errFound = true
// 				break
// 			}
// 		}
// 		if !errFound {
// 			fmt.Println("Report completed")
// 		}
// 	}
// }

// func sendMetric(compressedData []byte, ep string) error {
// 	s := fmt.Sprintf("http://%s/update/", ep)
// 	req, err := http.NewRequest("POST", s, bytes.NewBuffer(compressedData))
// 	if err != nil {
// 		return fmt.Errorf("error creating request: %w", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Content-Encoding", "gzip")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("error making HTTP request: %w", err)
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode >= 300 {
// 		return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
// 	}

// 	return nil
// }
