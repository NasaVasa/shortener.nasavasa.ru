package utils

import (
	"github.com/jackc/pgtype"
	"time"
)

// ParseLiveTime get string LiveTimeUnit that may be SECONDS, MINUTES, HOURS, DAYS and int LiveTime
// and return pgtype.Timestamp with UTC timezone and 3 hours offset
func ParseLiveTime(LiveTimeUnit string, LiveTime int) (pgtype.Timestamp, error) {
	var timeToLive pgtype.Timestamp
	switch LiveTimeUnit {
	case "SECONDS":
		timeToLive.Time = time.Now().UTC().Add(time.Second * time.Duration(LiveTime))
	case "MINUTES":
		timeToLive.Time = time.Now().UTC().Add(time.Minute * time.Duration(LiveTime))
	case "HOURS":
		timeToLive.Time = time.Now().UTC().Add(time.Hour * time.Duration(LiveTime))
	case "DAYS":
		timeToLive.Time = time.Now().UTC().Add(time.Hour * 24 * time.Duration(LiveTime))
	}
	timeToLive.Status = pgtype.Present
	return timeToLive, nil
}
