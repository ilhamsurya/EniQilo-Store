package database

import (
	"context"
	"fmt"
	"projectsphere/eniqlo-store/config"

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

func NewDatabase() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		config.GetString("DB_HOST"),
		config.GetString("DB_PORT"),
		config.GetString("DB_USERNAME"),
		config.GetString("DB_PASSWORD"),
		config.GetString("DB_NAME"),
	)

	env := config.GetString("ENV")

	if env == "production" {
		dsn += " sslmode=verify-full rootcert=ap-southeast-1-bundle.pem"
	} else {
		dsn += " sslmode=disable"
	}

	db, err := sqlx.Connect("pgx", dsn)

	return db, err
}
