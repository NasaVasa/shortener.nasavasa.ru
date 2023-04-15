package utils

func IsUrl(url string) bool {
	if len(url) > 7 && url[:7] == "http://" {
		return true
	}
	if len(url) > 8 && url[:8] == "https://" {
		return true
	}
	return false
}
