package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(connectionString string) (*DBStorage, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}
	return &DBStorage{db: db}, nil
}

func (s *DBStorage) Set(mType string, mKey string, mValue any) error {
	return nil
}

func (s *DBStorage) Get(mType string, mKey string) (any, error) {
	return nil, nil
}

func (s *DBStorage) GetAll() []string {
	return nil
}

func (s *DBStorage) Ping() error {
	return s.db.Ping()
}

func (s *DBStorage) Close() error {
	return s.db.Close()
}
