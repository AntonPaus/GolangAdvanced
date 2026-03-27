package interfaces

import "context"

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

type Storage interface {
	Set(ctx context.Context, mType string, mKey string, mValue any) error
	Get(ctx context.Context, mType string, mKey string) (any, error)
	GetAll(ctx context.Context) []string
	Ping() error
	Close() error
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
