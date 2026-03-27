package interfaces

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

type Storage interface {
	Set(mType string, mKey string, mValue any) error
	Get(mType string, mKey string) (any, error)
	GetAll() []string
	Ping() error
	Close() error
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
