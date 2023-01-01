package dal

import (
	"time"
)

type TokensPairModel struct {
	IsRevoked bool
	HRefreshToken string
	RefreshTokenExpiryTimestamp time.Time
	HAccessToken string
	AccessTokenExpiryTimestamp time.Time
	UserUuid string
}

