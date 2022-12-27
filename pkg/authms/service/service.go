package service

import (
	"errors"
	"time"

	"github.com/olegskip/messenger/pkg/authms/dal"
)

type IAuthService interface {
	GetNewRt(credentials CredentialsDto) (tokens TokensDto, err error)
	ExchangeRt(oldRt RtDto) (tokens TokensDto, err error)
	RevokeRt(oldRt RtDto) (err error)
	ValidateRt(rt RtDto) (err error)
}

type AuthService struct {
	tokenDao dal.ITokenDao
	rtTtl time.Duration
	jwtTtl time.Duration
}

func (a *AuthService) GetNewRt(credentials CredentialsDto) (tokens TokensDto, err error) {
	// Just for testing, if login equsl to password then return fake token		
	if credentials.Username == credentials.Password {
		tokens.Rt = a.generateRt()
		tokens.Jwt = a.generateJwt(tokens.Rt)
		err = nil
	} else {
		err = errors.New("Invalid username or password")
	}
	
	return tokens, err
}

func (a *AuthService) ExchangeRt(oldRt RtDto) (tokens TokensDto, err error) {
	if validationErr := a.ValidateRt(oldRt); validationErr != nil {
		return TokensDto{}, errors.New("Rt is invalid") 
	}
	// revoke the old refresh token to prevent creating multiple access tokens
	// if oldRt is already revoked we can't give an access token
	if revokingErr := a.RevokeRt(oldRt); revokingErr != nil {
		return TokensDto{}, errors.New("Rt is revoked")
	}

	tokens.Rt = a.generateRt()
	tokens.Jwt = a.generateJwt(tokens.Rt)

	return tokens, nil 
}

func (a *AuthService) RevokeRt(oldRt RtDto) (err error) {
	if validationErr := a.ValidateRt(oldRt); validationErr != nil {
		return errors.New("Rt is invalid") 
	}
	
	_, findingErr := a.tokenDao.FindRTByToken(oldRt.Token)
	if findingErr == nil {
		return errors.New("Rt is already revoked") 
	}

	a.tokenDao.AddRt(dtoToModelRt(oldRt))

	return nil
}

func (a *AuthService) ValidateRt(rt RtDto) (err error) {
	if time.Now().After(rt.ExpireTimestamp) {	
		return errors.New("Rt is expired") 
	}
	
	// TODO: actual validating

	return nil
}

func (a *AuthService) generateRt() (rt RtDto) {
	rt.Token = time.Now().String()
	rt.ExpireTimestamp = time.Now().Add(a.rtTtl)

	return rt
}

func (a *AuthService) generateJwt(rt RtDto) (jwt JwtDto) {
	jwt.Token = time.Now().String() + "." + rt.Token
	jwt.ExpireTimestamp = time.Now().Add(a.jwtTtl)

	return jwt
}

func NewAuthService(tokenDAO dal.ITokenDao, rtTtl time.Duration) *AuthService {
	return &AuthService{ 
		tokenDao: tokenDAO,
		rtTtl: rtTtl,
	}
}

func dtoToModelRt(rtDto RtDto) (rtModel dal.RtModel) {
	rtModel.Token = rtDto.Token
	rtModel.ExpireTimestamp = rtDto.ExpireTimestamp

	return rtModel
}

