package db

import sq "github.com/Masterminds/squirrel"

func (pgContext *PgContext) CheckShortUrl(shortUrl string) bool {
	query := pgContext.Psql.Select("short_url").
		From("url_storage").
		Where(sq.Eq{"short_url": shortUrl})

	sql, args, err := query.ToSql()
	if err != nil {
		return false
	}

	var shortUrlFromDB string
	err = pgContext.DB.Get(&shortUrlFromDB, sql, args...)
	if err != nil {
		return false
	}

	return true
}
