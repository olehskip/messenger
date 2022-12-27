package dal

import (
	"time"
)

type RtModel struct {
	Token string
	ExpireTimestamp time.Time
}

