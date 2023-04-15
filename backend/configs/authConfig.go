package configs

import "time"

const ttl = 30 * 24 * 60 * 60

type AuthConfig struct {
	TokenLiveTime time.Duration
}

func NewAuthConfig() AuthConfig {
	return AuthConfig{
		TokenLiveTime: time.Duration(ttl) * time.Second,
	}
}
