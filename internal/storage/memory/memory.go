package memory

import (
	"errors"
	"fmt"

	"github.com/AntonPaus/GolangAdvanced/internal/interfaces"
)

type MemoryStorage struct {
	metricsFloat64 map[string]float64
	metricsInt64   map[string]int64
}

func NewMemoryStorage() (*MemoryStorage, error) {
	m := &MemoryStorage{
		metricsFloat64: make(map[string]float64),
		metricsInt64:   make(map[string]int64),
	}
	return m, nil
}

func (m *MemoryStorage) Set(mType string, mKey string, mValue any) error {
	switch mType {
	case interfaces.MetricTypeGauge:
		g, ok := mValue.(float64)
		if !ok {
			return errors.New("invalid metric value, it is not gauge")
		}
		m.metricsFloat64[mKey] = g
	case interfaces.MetricTypeCounter:
		c, ok := mValue.(int64)
		if !ok {
			return errors.New("invalid metric value, it is not counter")
		}
		if _, ok := m.metricsInt64[mKey]; !ok {
			m.metricsInt64[mKey] = 0
		}
		m.metricsInt64[mKey] += c
	default:
		return errors.New("invalid metric type")
	}
	return nil
}

func (m *MemoryStorage) Get(mType string, mKey string) (any, error) {
	switch mType {
	case interfaces.MetricTypeGauge:
		if _, ok := m.metricsFloat64[mKey]; !ok {
			return nil, errors.New("metric not found")
		}
		return m.metricsFloat64[mKey], nil
	case interfaces.MetricTypeCounter:
		if _, ok := m.metricsInt64[mKey]; !ok {
			return nil, errors.New("metric not found")
		}
		return m.metricsInt64[mKey], nil
	default:
		return nil, errors.New("invalid metric type")
	}
}

func (m *MemoryStorage) GetAll() []string {
	result := []string{}
	for k, v := range m.metricsFloat64 {
		result = append(result, fmt.Sprintf("%s: %f", k, v))
	}
	for k, v := range m.metricsInt64 {
		result = append(result, fmt.Sprintf("%s: %d", k, v))
	}
	return result
}

// func (m *MemoryStorage) Delete(type string, key string) {
// 	switch type {
// 	case "gauge":
// 		delete(m.metricsFloat64, key)
// 	case "counter":
// 		delete(m.metricsInt64, key)
// 	}
// }
