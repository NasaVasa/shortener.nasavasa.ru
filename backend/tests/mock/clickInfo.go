package mock

import (
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"time"
)

type ClickInfo struct {
	id        uuid.UUID
	urlId     uuid.UUID
	clickTime time.Time
}

func NewClickInfo(urlId uuid.UUID) ClickInfo {
	return ClickInfo{
		id:        uuid.New(),
		urlId:     urlId,
		clickTime: time.Now().UTC(),
	}
}

func (pgMock *PgMock) InsertClickInfo(urlId pgtype.UUID) error {
	newUrlId := uuid.UUID{}
	newUrlId = urlId.Bytes

	pgMock.Clicks = append(pgMock.Clicks, NewClickInfo(newUrlId))
	return nil
}

func (pgMock *PgMock) GetAllRedirectTimes(id pgtype.UUID) ([]time.Time, error) {
	var times []time.Time
	for _, click := range pgMock.Clicks {
		if click.urlId == id.Bytes {
			times = append(times, click.clickTime)
		}
	}
	return times, nil
}
