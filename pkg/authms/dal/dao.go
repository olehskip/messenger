package dal

import (
	"errors"
)

type ITokenDao interface {
	AddRt(rt RtModel)
	FindRTByToken(rtToken string) (RtModel, error)
}

type InMemoryDao struct {	
	rts []RtModel
}

func (i *InMemoryDao) AddRt(rt RtModel) {
	i.rts = append(i.rts, rt)
}

func (i *InMemoryDao) FindRTByToken(token string) (RtModel, error) {
	for _, tokenModel := range(i.rts) {
		if tokenModel.Token == token {
			return tokenModel, nil
		}
	}
	
	return RtModel{}, errors.New("RT wasn't found")
}

