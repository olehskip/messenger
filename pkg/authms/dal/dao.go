package dal

import (
	"errors"
)

type IRevokedRefreshTokenDao interface {
	AddRefreshToken(rt RefreshTokenModel)
	FindRefreshTokenByToken(rtToken string) (RefreshTokenModel, error)
}

type InMemoryDao struct {	
	rts []RefreshTokenModel
}

func (i *InMemoryDao) AddRefreshToken(rt RefreshTokenModel) {
	i.rts = append(i.rts, rt)
}

func (i *InMemoryDao) FindRefreshTokenByToken(token string) (RefreshTokenModel, error) {
	for _, tokenModel := range(i.rts) {
		if tokenModel.Token == token {
			return tokenModel, nil
		}
	}
	
	return RefreshTokenModel{}, errors.New("RT wasn't found")
}

