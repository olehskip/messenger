package service

import (
	"time"
)

type CredentialsDto struct {
	Username string
	Password string
}

type RtDto struct {
	Token string
	ExpireTimestamp time.Time
}

type JwtDto struct {
	Token string
	ExpireTimestamp time.Time
}

type TokensDto struct {
	Rt RtDto
	Jwt JwtDto
}
