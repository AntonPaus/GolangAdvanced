package storage

import "context"

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

type Storage interface {
	Set(ctx context.Context, metrics []Metrics) error
	Get(ctx context.Context, mType string, mKey string) (any, error)
	GetAll(ctx context.Context) []string
	Ping() error
	Close() error
}
