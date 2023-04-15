package db

import sq "github.com/Masterminds/squirrel"

func (pgContext *PgContext) UpdateUserTokenTime(token string) error {
	query, args, err := pgContext.Psql.Update("users").
		Set("token_created_at", "NOW()").
		Where(sq.Eq{"token": token}).ToSql()
	if err != nil {
		return err
	}
	_, err = pgContext.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
