package memory

import (
	"errors"

	"github.com/AntonPaus/GolangAdvanced/internal/repository"
)

type MemoryStorage struct {
	metricsFloat64 map[string]float64
	metricsInt64   map[string]int64
}

func NewMemoryStorage() *MemoryStorage {
	m := &MemoryStorage{
		metricsFloat64: make(map[string]float64),
		metricsInt64:   make(map[string]int64),
	}
	return m
}

func (m *MemoryStorage) Set(mType string, mKey string, mValue any) error {
	switch mType {
	case repository.MetricTypeGauge:
		g, ok := mValue.(float64)
		if !ok {
			return errors.New("invalid metric value, it is not gauge")
		}
		m.metricsFloat64[mKey] = g
	case repository.MetricTypeCounter:
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
	case repository.MetricTypeGauge:
		if _, ok := m.metricsFloat64[mKey]; !ok {
			return nil, errors.New("metric not found")
		}
		return m.metricsFloat64[mKey], nil
	case repository.MetricTypeCounter:
		if _, ok := m.metricsInt64[mKey]; !ok {
			return nil, errors.New("metric not found")
		}
		return m.metricsInt64[mKey], nil
	default:
		return nil, errors.New("invalid metric type")
	}
}

// func (m *MemoryStorage) Delete(type string, key string) {
// 	switch type {
// 	case "gauge":
// 		delete(m.metricsFloat64, key)
// 	case "counter":
// 		delete(m.metricsInt64, key)
// 	}
// }
