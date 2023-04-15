package mock

import (
	"github.com/google/uuid"
	"shortener/utils"
	"time"
)

func FillTestDataInMock(pgMock *PgMock) {
	pgMock.Users = []User{
		NewUser(uuid.Nil, "user1", "example1@gmail.com", utils.HashPassword("11111"), uuid.NewString(), time.Now(), time.Now()),
		NewUser(uuid.Nil, "user2", "example2@gmail.com", utils.HashPassword("22222"), uuid.NewString(), time.Now(), time.Now()),
		NewUser(uuid.Nil, "user3", "example3@gmail.com", utils.HashPassword("33333"), uuid.NewString(), time.Now(), time.Now()),
	}

	pgMock.Urls = []Url{
		NewUrl("https://www.google.com", "google", time.Now().Add(time.Hour*24*7), false, uuid.Nil),
		NewUrl("https://www.yandex.ru", "yandex", time.Now().Add(time.Hour*24*7), true, pgMock.Users[0].id),
		NewUrl("https://www.youtube.com", "youtube", time.Now().Add(time.Hour*24*7), true, pgMock.Users[0].id),
		NewUrl("https://www.facebook.com", "facebook", time.Now().Add(time.Hour*24*7), true, pgMock.Users[1].id),
	}

	pgMock.Clicks = []ClickInfo{
		NewClickInfo(pgMock.Urls[0].Id),
		NewClickInfo(pgMock.Urls[0].Id),
		NewClickInfo(pgMock.Urls[1].Id),
		NewClickInfo(pgMock.Urls[1].Id),
		NewClickInfo(pgMock.Urls[1].Id),
		NewClickInfo(pgMock.Urls[2].Id),
	}
}

func ClearMock(pgMock *PgMock) {
	pgMock.Urls = []Url{}
	pgMock.Users = []User{}
	pgMock.Clicks = []ClickInfo{}
}
