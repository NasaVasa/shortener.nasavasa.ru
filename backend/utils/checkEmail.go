package utils

import "regexp"

// CheckEmail checks email using regexp
func CheckEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$`)
	return emailRegex.MatchString(email)
}
