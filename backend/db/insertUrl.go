package db

import "github.com/jackc/pgtype"

func (pgContext *PgContext) InsertUrl(longUrl string, shortUrl string, isVip bool, urlWillDelete pgtype.Timestamp, userId pgtype.UUID) error {
	query := pgContext.Psql.Insert("url_storage").
		Columns("long_url", "short_url")

	if isVip {
		// insert url_will_delete with UTC timezone
		query = query.Columns("url_will_delete", "is_vip", "user_id").
			Values(longUrl, shortUrl, urlWillDelete, true, userId)
	} else {
		query = query.Values(longUrl, shortUrl)
	}

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
