package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
	"time"
)

func (pgContext *PgContext) GetAllRedirectTimes(id pgtype.UUID) ([]time.Time, error) {
	query, args, err := pgContext.Psql.Select("click_time").
		From("click_info").
		Where(sq.Eq{"url_id": id}).ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := pgContext.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var times []time.Time
	for rows.Next() {
		var t time.Time
		err = rows.Scan(&t)
		if err != nil {
			return nil, err
		}
		times = append(times, t)
	}

	return times, nil
}
