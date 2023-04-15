package db

import (
	"github.com/jackc/pgtype"
	"time"
)

// PgCaller is an interface for PgContext
type PgCaller interface {
	CheckShortUrl(shortUrl string) bool
	CheckToken(token string) bool
	CheckUser(getBy string, value string) bool
	CreateUser(login string, password string, email string, token string) error
	DeleteUrl(deleteBy string, value string) error
	GetAllRedirectTimes(id pgtype.UUID) ([]time.Time, error)
	GetUrl(getBy string, url string, isVip bool) (Url, bool, error)
	GetUrlsList(token string) ([]string, error)
	GetUser(getBy string, value string) (User, error)
	IncrementUrlClicks(secretKey pgtype.UUID) error
	InsertClickInfo(urlId pgtype.UUID) error
	InsertUrl(longUrl string, shortUrl string, isVip bool, urlWillDelete pgtype.Timestamp, userId pgtype.UUID) error
	UpdateUser(user User) error
	UpdateUserTokenTime(token string) error
}
