package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
)

// Url is the struct for the url_storage table
type Url struct {
	Id            pgtype.UUID
	LongUrl       string
	ShortUrl      string
	SecretKey     pgtype.UUID
	UrlClicks     int
	UrlCreatedAt  pgtype.Timestamp
	UrlWillDelete pgtype.Timestamp
	IsVip         bool
}

// GetUrl returns all info from db for a given url, getBy can be "long_url" or "short_url"
func (pgContext *PgContext) GetUrl(getBy string, url string, isVip bool) (Url, bool, error) {
	basicSq := pgContext.Psql.Select("id", "long_url",
		"short_url",
		"secret_key",
		"url_clicks",
		"url_created_at",
		"url_will_delete",
		"is_vip").
		From("url_storage").
		Where(sq.Eq{getBy: url}).
		Where(sq.Eq{"is_vip": isVip})

	query, args, err := basicSq.ToSql()
	if err != nil {
		return Url{}, false, err
	}

	row := pgContext.DB.QueryRow(query, args...)
	var urlStruct Url
	err = row.Scan(&urlStruct.Id, &urlStruct.LongUrl,
		&urlStruct.ShortUrl,
		&urlStruct.SecretKey,
		&urlStruct.UrlClicks,
		&urlStruct.UrlCreatedAt,
		&urlStruct.UrlWillDelete,
		&urlStruct.IsVip)
	if err != nil {
		return Url{}, false, nil
	}

	return urlStruct, true, nil
}
