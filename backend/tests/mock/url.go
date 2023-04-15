package mock

import (
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"shortener/db"
	"time"
)

type Url struct {
	Id            uuid.UUID
	LongUrl       string
	ShortUrl      string
	SecretKey     uuid.UUID
	UrlClicks     int
	UrlCreatedAt  time.Time
	UrlWillDelete time.Time
	IsVip         bool
	UserId        uuid.UUID
}

func NewUrl(longUrl string, shortUrl string, urlWillDelete time.Time, isVip bool, userId uuid.UUID) Url {
	return Url{
		Id:            uuid.New(),
		LongUrl:       longUrl,
		ShortUrl:      shortUrl,
		SecretKey:     uuid.New(),
		UrlClicks:     0,
		UrlCreatedAt:  time.Now().UTC(),
		UrlWillDelete: urlWillDelete,
		IsVip:         isVip,
		UserId:        userId,
	}
}

func (pgMock *PgMock) CheckShortUrl(shortUrl string) bool {
	for _, url := range pgMock.Urls {
		if url.ShortUrl == shortUrl {
			return true
		}
	}
	return false
}

func (pgMock *PgMock) DeleteUrl(deleteBy string, value string) error {
	for i, url := range pgMock.Urls {
		if deleteBy == "Id" {
			if url.Id.String() == value {
				pgMock.Urls = append(pgMock.Urls[:i], pgMock.Urls[i+1:]...)
				return nil
			}
		}
		if deleteBy == "short_url" {
			if url.ShortUrl == value {
				pgMock.Urls = append(pgMock.Urls[:i], pgMock.Urls[i+1:]...)
				return nil
			}
		}
		if deleteBy == "secret_key" {
			if url.SecretKey.String() == value {
				pgMock.Urls = append(pgMock.Urls[:i], pgMock.Urls[i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (pgMock *PgMock) GetUrl(getBy string, url string, isVip bool) (db.Url, bool, error) {
	for _, urlStruct := range pgMock.Urls {
		if getBy == "secret_key" {
			var secretKey [16]byte
			copy(secretKey[:], urlStruct.SecretKey[:])
			var id [16]byte
			copy(id[:], urlStruct.Id[:])
			if urlStruct.SecretKey.String() == url && urlStruct.IsVip == isVip {
				// remake url struct to db.Url and return
				return db.Url{
					Id:            pgtype.UUID{Bytes: id, Status: pgtype.Present},
					LongUrl:       urlStruct.LongUrl,
					ShortUrl:      urlStruct.ShortUrl,
					SecretKey:     pgtype.UUID{Bytes: secretKey, Status: pgtype.Present},
					UrlClicks:     urlStruct.UrlClicks,
					UrlWillDelete: pgtype.Timestamp{Time: urlStruct.UrlWillDelete, Status: pgtype.Present},
					IsVip:         urlStruct.IsVip,
				}, true, nil
			}
		}
		if getBy == "long_url" {
			var secretKey [16]byte
			copy(secretKey[:], urlStruct.SecretKey[:])
			var id [16]byte
			copy(id[:], urlStruct.Id[:])
			if urlStruct.LongUrl == url && urlStruct.IsVip == isVip {
				// remake url struct to db.Url and return
				return db.Url{
					Id:            pgtype.UUID{Bytes: id, Status: pgtype.Present},
					LongUrl:       urlStruct.LongUrl,
					ShortUrl:      urlStruct.ShortUrl,
					SecretKey:     pgtype.UUID{Bytes: secretKey, Status: pgtype.Present},
					UrlClicks:     urlStruct.UrlClicks,
					UrlWillDelete: pgtype.Timestamp{Time: urlStruct.UrlWillDelete, Status: pgtype.Present},
					IsVip:         urlStruct.IsVip,
				}, true, nil
			}
		}
		if getBy == "short_url" {
			var secretKey [16]byte
			copy(secretKey[:], urlStruct.SecretKey[:])
			var id [16]byte
			copy(id[:], urlStruct.Id[:])
			if urlStruct.ShortUrl == url && urlStruct.IsVip == isVip {
				return db.Url{
					Id:            pgtype.UUID{Bytes: id, Status: pgtype.Present},
					LongUrl:       urlStruct.LongUrl,
					ShortUrl:      urlStruct.ShortUrl,
					SecretKey:     pgtype.UUID{Bytes: secretKey, Status: pgtype.Present},
					UrlClicks:     urlStruct.UrlClicks,
					UrlWillDelete: pgtype.Timestamp{Time: urlStruct.UrlWillDelete, Status: pgtype.Present},
					IsVip:         urlStruct.IsVip,
				}, true, nil
			}
		}
	}
	return db.Url{}, false, nil
}

func (pgMock *PgMock) GetUrlsList(token string) ([]string, error) {
	var userId uuid.UUID
	for _, user := range pgMock.Users {
		if user.Token == token {
			userId = user.id
		}
	}

	var secretKeys []string
	for _, url := range pgMock.Urls {
		if url.UserId == userId {
			secretKeys = append(secretKeys, url.SecretKey.String())
		}
	}

	return secretKeys, nil
}

func (pgMock *PgMock) IncrementUrlClicks(secretKey pgtype.UUID) error {
	var secretKeyString [16]byte
	copy(secretKeyString[:], secretKey.Bytes[:])
	for i, url := range pgMock.Urls {
		if url.SecretKey == secretKeyString {
			pgMock.Urls[i].UrlClicks++
			return nil
		}
	}
	return nil
}

func (pgMock *PgMock) InsertUrl(longUrl string, shortUrl string, isVip bool, urlWillDelete pgtype.Timestamp, userId pgtype.UUID) error {
	newUserId := uuid.UUID{}
	newUserId = userId.Bytes
	pgMock.Urls = append(pgMock.Urls, NewUrl(longUrl, shortUrl, urlWillDelete.Time, isVip, newUserId))
	return nil
}
