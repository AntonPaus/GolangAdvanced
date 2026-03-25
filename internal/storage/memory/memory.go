package memory

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AntonPaus/GolangAdvanced/internal/interfaces"
)

type MemoryStorage struct {
	// mu             sync.Mutex
	metricsFloat64 map[string]float64
	metricsInt64   map[string]int64
	dumpInterval   uint
	dumpFile       *os.File
	scanner        *bufio.Scanner
}

func NewMemoryStorage(restore bool, fileStoragePath string, dumpInterval uint) (*MemoryStorage, error) {
	m := &MemoryStorage{
		dumpInterval:   dumpInterval,
		metricsFloat64: make(map[string]float64),
		metricsInt64:   make(map[string]int64),
	}
	err := error(nil)
	m.dumpFile, err = os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	m.scanner = bufio.NewScanner(m.dumpFile)
	if restore {
		err := m.restoreFromFile()
		if err != nil {
			fmt.Println("No storage file found. Continue")
		}
	}
	if dumpInterval > 0 {
		go m.tickerDump()
	}
	return m, nil
}

func (m *MemoryStorage) Set(mType string, mKey string, mValue any) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
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

func (m *MemoryStorage) dump() {
	dataGauge, err := json.Marshal(m.metricsFloat64)
	if err != nil {
		fmt.Println("cannot marshal data: %w", err)
	}
	dataGauge = append(dataGauge, '\n')
	dataCounter, err := json.Marshal(m.metricsInt64)
	if err != nil {
		fmt.Println("cannot marshal data: %w", err)
	}
	os.WriteFile(m.dumpFile.Name(), append(dataGauge, dataCounter...), 0666)
}

func (m *MemoryStorage) tickerDump() {
	ticker := time.NewTicker(time.Duration(m.dumpInterval) * time.Second)
	for {
		<-ticker.C
		// fmt.Println(int(t.Second()))
		m.dump()
	}
}

func (m *MemoryStorage) restoreFromFile() error {
	if !m.scanner.Scan() {
		return m.scanner.Err()
	}
	data := m.scanner.Bytes()
	err := json.Unmarshal(data, &m.metricsFloat64)
	if err != nil {
		return err
	}
	if !m.scanner.Scan() {
		return m.scanner.Err()
	}
	data = m.scanner.Bytes()
	err = json.Unmarshal(data, &m.metricsInt64)
	if err != nil {
		return err
	}
	return nil
}
