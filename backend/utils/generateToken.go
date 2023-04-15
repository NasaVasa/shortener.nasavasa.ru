package utils

import (
	uuid "github.com/google/uuid"
	"shortener/db"
)

// GenerateToken generates a random uuid token
func GenerateToken(pgContext db.PgCaller) uuid.UUID {
	token, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	for pgContext.CheckToken(token.String()) {
		token, err = uuid.NewUUID()
		if err != nil {
			panic(err)
		}
	}

	return token
}
