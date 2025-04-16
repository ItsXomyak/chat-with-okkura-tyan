package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
  *sqlx.DB
}

func New(connStr string) (*DB, error) {
  db, err := sqlx.Connect("postgres", connStr)
  if err != nil {
    return nil, err
  }
  return &DB{db}, nil
}