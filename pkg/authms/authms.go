package main

import (
	"errors"	
)

type IAuthMs interface {
	Auth(username string, password string) (string, error);
}

type AuthMs struct {
}

func (a *AuthMs) Auth(username string, password string) (string, error) {
	// Just for testing, if login equsl to password then return fake token
	if username == password {
		return "fake token", nil
	} else {
		return "", errors.New("Bad username or password")
	}
}

