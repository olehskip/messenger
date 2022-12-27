package service

import (
	"time"
)

type CredentialsDto struct {
	Username string
	Password string
}

type RefreshTokenDto struct {
	Token string
	ExpireTimestamp time.Time
}

type AccessTokenDto struct {
	Token string
	ExpireTimestamp time.Time
}

type TokensDto struct {
	RefreshToken RefreshTokenDto
	AccessToken AccessTokenDto
}
