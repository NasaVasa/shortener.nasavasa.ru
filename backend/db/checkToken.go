package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
)

func (pgContext *PgContext) CheckToken(token string) bool {
	query := pgContext.Psql.Select("id").
		From("users").
		Where(sq.Eq{"token": token})

	sql, args, err := query.ToSql()
	if err != nil {
		return false
	}

	var idFromDB pgtype.UUID
	err = pgContext.DB.Get(&idFromDB, sql, args...)
	if err != nil {
		return false
	}

	return true
}
