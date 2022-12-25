package service

import (
	"errors"
	"time"
)

type IAuthService interface {
	GetNewRT(username string, password string) (rt string, jwt string, err error);
}

type AuthService struct {
	tokenDAO ITokenDAO
}

func (a *AuthService) GetNewRT(username string, password string) (rt string, jwt string, err error) {
	// Just for testing, if login equsl to password then return fake token	
	if username == password {
		rt = a.generateRT(username)
		jwt = a.generateJWT(rt, username)
		err = nil
	} else {
		rt = ""
		jwt = ""
		err = errors.New("Invalid username or password")
	}

	return rt, jwt, err
}

func (a *AuthService) generateRT(username string) (rt string) {
	rt = username + "/" + time.Now().String()

	return rt
}

func (a *AuthService) generateJWT(username string, rt string) (jwt string) {
	jwt = username + "/" + time.Now().String() + "." + rt

	return jwt
}

func NewAuthService(tokenDAO ITokenDAO) *AuthService {
	return &AuthService{ 
		tokenDAO: tokenDAO,
	}
}

