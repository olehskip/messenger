package dal

import (
	"time"
)

type RefreshTokenModel struct {
	Token string
	ExpireTimestamp time.Time
}

