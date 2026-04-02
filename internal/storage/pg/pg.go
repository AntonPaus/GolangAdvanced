package pg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AntonPaus/GolangAdvanced/internal/logger"
	"github.com/AntonPaus/GolangAdvanced/internal/model"
	"github.com/AntonPaus/GolangAdvanced/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type Storage struct {
	conn *sql.DB
}

func NewStorage(connectionString string) (*Storage, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		logger.Log.Error("cannot create tables", zap.Error(err))
		return nil, err
	}
	logger.Log.Info("Database connected")
	return &Storage{conn: db}, nil
}

func (s *Storage) Bootstrap() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS counters (
            id TEXT PRIMARY KEY,
            value BIGINT
        )
    `)
	tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS gauges (
            id TEXT PRIMARY KEY,
            value DOUBLE PRECISION
        )
    `)
	logger.Log.Info("Tables created")
	return tx.Commit()
}

// func (s *Storage) Set(ctx context.Context, mType string, mKey string, mValue any) error {
// 	switch mType {
// 	case model.MetricTypeGauge:
// 		_, err := s.conn.ExecContext(ctx, "INSERT INTO gauges (id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = excluded.value;", mKey, mValue)
// 		if err != nil {
// 			logger.Log.Error("cannot set gauge", zap.Error(err))
// 			return err
// 		}
// 	case model.MetricTypeCounter:
// 		var currentValue int64
// 		err := s.conn.QueryRowContext(ctx, "SELECT value FROM counters WHERE id = $1", mKey).Scan(&currentValue)
// 		if err == nil {
// 			_, err = s.conn.ExecContext(ctx, "UPDATE counters SET value = value + $1 WHERE id = $2", mValue, mKey)
// 			if err != nil {
// 				logger.Log.Error("cannot increment counter", zap.Error(err))
// 				return err
// 			}
// 			return nil
// 		} else if err != sql.ErrNoRows {
// 			_, err := s.conn.ExecContext(ctx, "INSERT INTO counters (id, value) VALUES ($1, $2);", mKey, mValue)
// 			if err != nil {
// 				logger.Log.Error("cannot create counter", zap.Error(err))
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

func (s *Storage) Set(ctx context.Context, metrics []storage.Metrics) error {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, m := range metrics {
		switch m.MType {
		case model.MetricTypeGauge:
			_, err := tx.ExecContext(ctx, `
				INSERT INTO gauges (id, value) VALUES ($1, $2) 
				ON CONFLICT (id) 
				DO UPDATE SET value = excluded.value;`,
				m.ID, *m.Value)
			if err != nil {
				logger.Log.Error("cannot set gauge", zap.Error(err))
				return err
			}
		case model.MetricTypeCounter:
			_, err = tx.ExecContext(ctx, `
				INSERT INTO counters (id, value) 
				VALUES ($1, $2) 
				ON CONFLICT (id) DO 
				UPDATE SET value = counters.value + EXCLUDED.value;`,
				m.ID, *m.Delta)
			if err != nil {
				logger.Log.Error("cannot increment counter", zap.Error(err))
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *Storage) Get(ctx context.Context, mType string, mKey string) (any, error) {
	switch mType {
	case model.MetricTypeGauge:
		var v float64
		err := s.conn.QueryRowContext(ctx, "SELECT value FROM gauges WHERE id = $1", mKey).Scan(&v)
		if err != nil {
			logger.Log.Error("cannot get gauge", zap.Error(err))
			return nil, err
		}
		return v, nil
	case model.MetricTypeCounter:
		var v int64
		err := s.conn.QueryRowContext(ctx, "SELECT value FROM counters WHERE id = $1", mKey).Scan(&v)
		if err != nil {
			logger.Log.Error("cannot get gauge", zap.Error(err))
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", mType)
	}
}

func (s *Storage) GetAll(ctx context.Context) []string {
	return nil
}

func (s *Storage) Ping() error {
	return s.conn.Ping()
}

func (s *Storage) Close() error {
	return s.conn.Close()
}
