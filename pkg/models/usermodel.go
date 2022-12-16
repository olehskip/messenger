package usermodel

import (
	"time"
)

type UserModel struct {
	Bio string
	Id string
	Name string
	RegTimestamp time.Time
	Username string
}
