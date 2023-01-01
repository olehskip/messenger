package service

type CredentialsDto struct {
	UserUuid string
	Password string
}

type HTokensPairDto struct {
	HashedRefreshToken string
	HashedAccessToken string
}
