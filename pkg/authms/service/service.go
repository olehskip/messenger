package service

import (
	"errors"
	"time"

	"github.com/olegskip/messenger/pkg/authms/dal"
)

type IAuthService interface {
	GetNewRt(credentials CredentialsDto) (rt RtDto, jwt JwtDto, err error)
	ExchangeRt(oldRt RtDto) (newRt RtDto, err error)
}

type AuthService struct {
	tokenDao dal.ITokenDao
}

func (a *AuthService) GetNewRt(credentials CredentialsDto) (rt RtDto, jwt JwtDto, err error) {
	// Just for testing, if login equsl to password then return fake token	
	if credentials.Username == credentials.Password {
		rt = a.generateRt()
		jwt = a.generateJwt(rt)
		a.tokenDao.AddNewRt(dal.RtModel{Token: rt.Token})
		err = nil
	} else {
		rt.Token = ""
		jwt.Token = ""
		err = errors.New("Invalid username or password")
	}

	return rt, jwt, err
}

func (a *AuthService) ExchangeRt(oldRt RtDto) (newRt RtDto, err error) {
	oldRtModel, oldTokenFindingErr := a.tokenDao.FindRTByToken(oldRt.Token)
	if oldTokenFindingErr != nil {
		return RtDto{}, errors.New("Rt is not found") 
	}

	if oldRtModel.IsRevoked {
		return RtDto{}, errors.New("Rt is revoked") 
	}

	if time.Now().Before(oldRtModel.ExpireTimestamp) {	
		return RtDto{}, errors.New("Rt is expired") 
	}

	a.tokenDao.RevokeRt(oldRtModel)
	newRt = a.generateRt()
	a.tokenDao.AddNewRt(dal.RtModel{Token: newRt.Token})

	return newRt, errors.New("Rt is not valid") 
}

func (a *AuthService) generateRt() (rt RtDto) {
	rt.Token = time.Now().String()

	return rt
}

func (a *AuthService) generateJwt(rt RtDto) (jwt JwtDto) {
	jwt.Token = time.Now().String() + "." + rt.Token

	return jwt
}

func NewAuthService(tokenDAO dal.ITokenDao) *AuthService {
	return &AuthService{ 
		tokenDao: tokenDAO,
	}
}

