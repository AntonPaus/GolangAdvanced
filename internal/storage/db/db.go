package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AntonPaus/GolangAdvanced/internal/interfaces"
	"github.com/AntonPaus/GolangAdvanced/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const createTables = `
CREATE TABLE IF NOT EXISTS gauges(
	id TEXT PRIMARY KEY,
	value DOUBLE PRECISION
);
CREATE TABLE IF NOT EXISTS counters(
	id TEXT PRIMARY KEY,
	value INTEGER
);`

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(connectionString string) (*DBStorage, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		logger.Log.Error("cannot create tables", zap.Error(err))
		return nil, err
	}
	logger.Log.Info("Database connected")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, createTables)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("Tables created")
	return &DBStorage{db: db}, nil
}

func (s *DBStorage) Set(ctx context.Context, mType string, mKey string, mValue any) error {
	switch mType {
	case interfaces.MetricTypeGauge:
		_, err := s.db.ExecContext(ctx, "INSERT INTO gauges (id, value) VALUES ($1, $2)", mKey, mValue)
		if err != nil {
			logger.Log.Error("cannot set gauge", zap.Error(err))
			return err
		}
	case interfaces.MetricTypeCounter:
		_, err := s.db.ExecContext(context.Background(), "INSERT INTO counters (id, value) VALUES ($1, $2)", mKey, mValue)
		if err != nil {
			logger.Log.Error("cannot set counter", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *DBStorage) Get(ctx context.Context, mType string, mKey string) (any, error) {
	switch mType {
	case interfaces.MetricTypeGauge:
		var v float64
		err := s.db.QueryRowContext(ctx, "SELECT value FROM gauges WHERE id = $1", mKey).Scan(&v)
		if err != nil {
			logger.Log.Error("cannot get gauge", zap.Error(err))
			return nil, err
		}
		return v, nil
	case interfaces.MetricTypeCounter:
		var v int64
		err := s.db.QueryRowContext(ctx, "SELECT value FROM counters WHERE id = $1", mKey).Scan(&v)
		if err != nil {
			logger.Log.Error("cannot get gauge", zap.Error(err))
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", mType)
	}
}

func (s *DBStorage) GetAll(ctx context.Context) []string {
	return nil
}

func (s *DBStorage) Ping() error {
	return s.db.Ping()
}

func (s *DBStorage) Close() error {
	return s.db.Close()
}
