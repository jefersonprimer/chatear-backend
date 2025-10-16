package constants

import (
	"time"
)

const (
	AccessTokenExpiration  = 15 * time.Minute
	RefreshTokenExpiration = 7 * 24 * time.Hour // 7 days
)

var (
	JwtSecret = []byte("your-secret-key") // TODO: Load from environment variable
)