package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AntonPaus/GolangAdvanced/internal/interfaces"
)

type FileStorage struct {
	// mu             sync.Mutex
	metricsFloat64 map[string]float64
	metricsInt64   map[string]int64
	dumpInterval   uint
	dumpFile       *os.File
	scanner        *bufio.Scanner
}

func NewFileStorage(restore bool, fileStoragePath string, dumpInterval uint) (*FileStorage, error) {
	s := &FileStorage{
		dumpInterval:   dumpInterval,
		metricsFloat64: make(map[string]float64),
		metricsInt64:   make(map[string]int64),
	}
	err := error(nil)
	s.dumpFile, err = os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	s.scanner = bufio.NewScanner(s.dumpFile)
	if restore {
		err := s.restoreFromFile()
		if err != nil {
			fmt.Println("No storage file found. Continue")
		}
	}
	if dumpInterval > 0 {
		go s.tickerDump()
	}
	return s, nil
}

func (s *FileStorage) Set(mType string, mKey string, mValue any) error {
	// s.mu.Lock()
	// defer s.mu.Unlock()
	switch mType {
	case interfaces.MetricTypeGauge:
		g, ok := mValue.(float64)
		if !ok {
			return errors.New("invalid metric value, it is not gauge")
		}
		s.metricsFloat64[mKey] = g
	case interfaces.MetricTypeCounter:
		c, ok := mValue.(int64)
		if !ok {
			return errors.New("invalid metric value, it is not counter")
		}
		if _, ok := s.metricsInt64[mKey]; !ok {
			s.metricsInt64[mKey] = 0
		}
		s.metricsInt64[mKey] += c
	default:
		return errors.New("invalid metric type")
	}
	return nil
}

func (s *FileStorage) Get(mType string, mKey string) (any, error) {
	switch mType {
	case interfaces.MetricTypeGauge:
		if _, ok := s.metricsFloat64[mKey]; !ok {
			return nil, errors.New("metric not found")
		}
		return s.metricsFloat64[mKey], nil
	case interfaces.MetricTypeCounter:
		if _, ok := s.metricsInt64[mKey]; !ok {
			return nil, errors.New("metric not found")
		}
		return s.metricsInt64[mKey], nil
	default:
		return nil, errors.New("invalid metric type")
	}
}

func (s *FileStorage) GetAll() []string {
	result := []string{}
	for k, v := range s.metricsFloat64 {
		result = append(result, fmt.Sprintf("%s: %f", k, v))
	}
	for k, v := range s.metricsInt64 {
		result = append(result, fmt.Sprintf("%s: %d", k, v))
	}
	return result
}

func (s *FileStorage) dump() {
	dataGauge, err := json.Marshal(s.metricsFloat64)
	if err != nil {
		fmt.Println("cannot marshal data: %w", err)
	}
	dataGauge = append(dataGauge, '\n')
	dataCounter, err := json.Marshal(s.metricsInt64)
	if err != nil {
		fmt.Println("cannot marshal data: %w", err)
	}
	os.WriteFile(s.dumpFile.Name(), append(dataGauge, dataCounter...), 0666)
}

func (s *FileStorage) tickerDump() {
	ticker := time.NewTicker(time.Duration(s.dumpInterval) * time.Second)
	for {
		<-ticker.C
		s.dump()
	}
}

func (s *FileStorage) restoreFromFile() error {
	if !s.scanner.Scan() {
		return s.scanner.Err()
	}
	data := s.scanner.Bytes()
	err := json.Unmarshal(data, &s.metricsFloat64)
	if err != nil {
		return err
	}
	if !s.scanner.Scan() {
		return s.scanner.Err()
	}
	data = s.scanner.Bytes()
	err = json.Unmarshal(data, &s.metricsInt64)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) Ping() error {
	return nil
}

func (s *FileStorage) Close() error {
	return s.dumpFile.Close()
}
