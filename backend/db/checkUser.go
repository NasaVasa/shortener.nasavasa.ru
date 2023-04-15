package db

import sq "github.com/Masterminds/squirrel"

func (pgContext *PgContext) CheckUser(getBy string, value string) bool {
	query := pgContext.Psql.Select(getBy).
		From("users").
		Where(sq.Eq{getBy: value})

	sql, args, err := query.ToSql()
	if err != nil {
		return false
	}

	var valueFromDB string
	err = pgContext.DB.Get(&valueFromDB, sql, args...)
	if err != nil {
		return false
	}

	return true
}
