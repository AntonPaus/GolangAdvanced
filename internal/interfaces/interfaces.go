package interfaces

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

type MetricsStorage interface {
	Set(mType string, mKey string, mValue any) error
	Get(mType string, mKey string) (any, error)
	GetAll() []string
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
