package memory

import (
	"context"
	"errors"
	"fmt"

	"github.com/AntonPaus/GolangAdvanced/internal/model"
	"github.com/AntonPaus/GolangAdvanced/internal/storage"
)

type MemoryStorage struct {
	// mu             sync.Mutex
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

func (s *MemoryStorage) Set(_ context.Context, metrics []storage.Metrics) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	for _, m := range metrics {
		switch m.MType {
		case model.MetricTypeGauge:
			s.metricsFloat64[m.ID] = *m.Value
		case model.MetricTypeCounter:
			if _, ok := s.metricsInt64[m.ID]; !ok {
				s.metricsInt64[m.ID] = 0
			}
			s.metricsInt64[m.ID] += *m.Delta
		default:
			return errors.New("invalid metric type")
		}
	}
	return nil
}

func (s *MemoryStorage) Get(_ context.Context, m storage.Metrics) (any, error) {
	switch m.MType {
	case model.MetricTypeGauge:
		if _, ok := s.metricsFloat64[m.ID]; !ok {
			return nil, errors.New("metric not found")
		}
		return s.metricsFloat64[m.ID], nil
	case model.MetricTypeCounter:
		if _, ok := s.metricsInt64[m.ID]; !ok {
			return nil, errors.New("metric not found")
		}
		return s.metricsInt64[m.ID], nil
	default:
		return nil, errors.New("invalid metric type")
	}
}

func (s *MemoryStorage) GetAll(_ context.Context) []string {
	result := []string{}
	for k, v := range s.metricsFloat64 {
		result = append(result, fmt.Sprintf("%s: %f", k, v))
	}
	for k, v := range s.metricsInt64 {
		result = append(result, fmt.Sprintf("%s: %d", k, v))
	}
	return result
}

func (s *MemoryStorage) Ping() error {
	return nil
}

func (s *MemoryStorage) Close() error {
	return nil
}
