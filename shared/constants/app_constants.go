package constants

import (
	os "os"
	"time"
)

const (
	AccessTokenExpiration  = 15 * time.Minute
	RefreshTokenExpiration = 7 * 24 * time.Hour // 7 days
)

var (
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
)