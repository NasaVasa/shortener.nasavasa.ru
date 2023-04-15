package utils

import (
	"math/rand"
	"shortener/db"
)

// GenerateShortUrl gereates a short url with a length of 6 and characters from A-Z, 0-9
// check is done to ensure that the generated short url is not already present in the database
// and returns it as a string with hostname as prefix
func GenerateShortUrl(pgContext db.PgCaller) string {
	// generate a random string of length 6
	shortUrl := generateRandomString(6)

	// check if the generated short url is already present in the database
	// if yes, generate a new one
	for pgContext.CheckShortUrl(shortUrl) {
		shortUrl = generateRandomString(6)
	}

	// return the short url with hostname as prefix
	return shortUrl
}

// generateRandomString generates a random string of length n and returns it as a string
func generateRandomString(n int) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
