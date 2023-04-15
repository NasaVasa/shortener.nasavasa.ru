package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
)

func (pgContext *PgContext) IncrementUrlClicks(secretKey pgtype.UUID) error {
	query := pgContext.Psql.Update("url_storage").
		Set("url_clicks", sq.Expr("url_clicks + 1")).
		Where(sq.Eq{"secret_key": secretKey})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pgContext.DB.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}
