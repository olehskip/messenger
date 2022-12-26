package dal

import (
	"time"
)

type RtModel struct {
	Token string
	UserId string
	ExpireTimestamp time.Time
	IsRevoked bool
}

