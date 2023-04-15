package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
)

func (pgContext *PgContext) GetUrlsList(token string) ([]string, error) {
	query, args, err := pgContext.Psql.Select("id").
		From("users").
		Where(sq.Eq{"token": token}).ToSql()
	if err != nil {
		return nil, err
	}

	var id pgtype.UUID
	err = pgContext.DB.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return nil, &NotFoundUserError{}
	}

	query, args, err = pgContext.Psql.Select("secret_key").
		From("url_storage").
		Where(sq.Eq{"user_id": id}).ToSql()

	if err != nil {
		return nil, err
	}

	var secretKeys []string
	err = pgContext.DB.Select(&secretKeys, query, args...)
	if err != nil {
		return nil, err
	}

	return secretKeys, nil
}
