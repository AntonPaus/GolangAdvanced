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
