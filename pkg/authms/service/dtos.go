package service

type CredentialsDto struct {
	Username string
	Password string
}

type RtDto struct {
	Token string
	// UserId string
}

type JwtDto struct {
	Token string
	// UserId string
}

type TokensDto struct {
	rt RtDto
	jwt JwtDto
}
