package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type PostgresConnector struct {
	DB *sqlx.DB
}

func NewPostgresConnector(ctx context.Context, db *sqlx.DB) PostgresConnector {
	return PostgresConnector{
		DB: db,
	}
}
