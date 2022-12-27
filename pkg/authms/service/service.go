package service

import (
	"errors"
	"time"

	"github.com/olegskip/messenger/pkg/authms/dal"
)

type IAuthService interface {
	GetNewRefreshToken(credentials CredentialsDto) (tokens TokensDto, err error)
	ExchangeRefreshToken(oldRefreshToken RefreshTokenDto) (tokens TokensDto, err error)
	RevokeRefreshToken(oldRefreshToken RefreshTokenDto) (err error)
	ValidateRefreshToken(refreshToken RefreshTokenDto) (err error)
}

type AuthService struct {
	revokedRefreshTokenDao dal.IRevokedRefreshTokenDao
	refreshTokenTtl time.Duration
	jwtTtl time.Duration
}

func (a *AuthService) GetNewRefreshToken(credentials CredentialsDto) (tokens TokensDto, err error) {
	// Just for testing, if login equsl to password then return fake token		
	if credentials.Username == credentials.Password {
		tokens.RefreshToken = a.generateRefreshToken()
		tokens.AccessToken = a.generateAccessToken(tokens.RefreshToken)
		err = nil
	} else {
		err = errors.New("Invalid username or password")
	}
	
	return tokens, err
}

func (a *AuthService) ExchangeRefreshToken(oldRefreshToken RefreshTokenDto) (tokens TokensDto, err error) {
	if validationErr := a.ValidateRefreshToken(oldRefreshToken); validationErr != nil {
		return TokensDto{}, errors.New("RefreshToken is invalid") 
	}
	// revoke the old refresh token to prevent creating multiple access tokens
	// if oldRefreshToken is already revoked we can't give an access token
	if revokingErr := a.RevokeRefreshToken(oldRefreshToken); revokingErr != nil {
		return TokensDto{}, errors.New("RefreshToken is revoked")
	}

	tokens.RefreshToken = a.generateRefreshToken()
	tokens.AccessToken = a.generateAccessToken(tokens.RefreshToken)

	return tokens, nil 
}

func (a *AuthService) RevokeRefreshToken(oldRefreshToken RefreshTokenDto) (err error) {
	if validationErr := a.ValidateRefreshToken(oldRefreshToken); validationErr != nil {
		return errors.New("RefreshToken is invalid") 
	}
	
	_, findingErr := a.revokedRefreshTokenDao.FindRefreshTokenByToken(oldRefreshToken.Token)
	if findingErr == nil {
		return errors.New("RefreshToken is already revoked") 
	}

	a.revokedRefreshTokenDao.AddRefreshToken(dtoToModelRefreshToken(oldRefreshToken))

	return nil
}

func (a *AuthService) ValidateRefreshToken(refreshToken RefreshTokenDto) (err error) {
	if time.Now().After(refreshToken.ExpireTimestamp) {	
		return errors.New("RefreshToken is expired") 
	}
	
	// TODO: actual validating

	return nil
}

func (a *AuthService) generateRefreshToken() (refreshToken RefreshTokenDto) {
	refreshToken.Token = time.Now().String()
	refreshToken.ExpireTimestamp = time.Now().Add(a.refreshTokenTtl)

	return refreshToken
}

func (a *AuthService) generateAccessToken(refreshToken RefreshTokenDto) (jwt AccessTokenDto) {
	jwt.Token = time.Now().String() + "." + refreshToken.Token
	jwt.ExpireTimestamp = time.Now().Add(a.jwtTtl)

	return jwt
}

func NewAuthService(revokedRefreshTokenDao dal.IRevokedRefreshTokenDao, refreshTokenTtl time.Duration) *AuthService {
	return &AuthService{ 
		revokedRefreshTokenDao: revokedRefreshTokenDao,
		refreshTokenTtl: refreshTokenTtl,
	}
}

func dtoToModelRefreshToken(refreshTokenDto RefreshTokenDto) (refreshTokenModel dal.RefreshTokenModel) {
	refreshTokenModel.Token = refreshTokenDto.Token
	refreshTokenModel.ExpireTimestamp = refreshTokenDto.ExpireTimestamp

	return refreshTokenModel
}

