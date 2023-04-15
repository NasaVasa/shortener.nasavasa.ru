package db

import sq "github.com/Masterminds/squirrel"

func (pgContext *PgContext) DeleteUrl(deleteBy string, value string) error {
	query := pgContext.Psql.Delete("url_storage").
		Where(sq.Eq{deleteBy: value})

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
