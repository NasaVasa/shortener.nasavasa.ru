package configs

import "os"

var urlPrefix = os.Getenv("URL_PREFIX")

// UrlConfig is the struct for the url config
type UrlConfig struct {
	UrlPrefix string
}

// NewUrlConfig creates a new url config
func NewUrlConfig() UrlConfig {
	return UrlConfig{
		UrlPrefix: urlPrefix,
	}
}
