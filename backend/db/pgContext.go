package db

import (
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"os"
)

// PgContext struct with the sqlx database connection
type PgContext struct {
	DB   *sqlx.DB
	Psql sq.StatementBuilderType
}

// NewDB creates a new database sqlx connection with standard settings
func NewDB() (PgCaller, error) {
	db, err := sqlx.Open("pgx", os.Getenv("CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	pgContext := &PgContext{
		DB:   db,
		Psql: psql,
	}
	return pgContext, nil
}
