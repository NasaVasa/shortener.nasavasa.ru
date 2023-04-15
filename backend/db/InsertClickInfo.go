package db

import "github.com/jackc/pgtype"

func (pgContext *PgContext) InsertClickInfo(urlId pgtype.UUID) error {
	query := pgContext.Psql.Insert("click_info").
		Columns("url_id").
		Values(urlId)
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
