package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the password using bcrypt
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}
