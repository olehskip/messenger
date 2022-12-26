package dal

import (
	// "github.com/olegskip/messenger/pkg/models"

	"errors"
)

type ITokenDao interface {
	AddNewRt(rt RtModel)
	FindRTByToken(rtToken string) (RtModel, error)
	DeleteAllRTsByUserId(userId string)
	RevokeRt(rt RtModel)

	GetRTModel(rt string) (RtModel, error)
}

type InMemoryDao struct {	
	rts []RtModel
}

func (i *InMemoryDao) AddNewRt(rt RtModel) {
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

func (i *InMemoryDao) RevokeRt(rt RtModel) {
	for j := 0; j < len(i.rts); j++ {
		if i.rts[j] == rt {
			i.rts[j].IsRevoked = true
		}
	}
}

func (i *InMemoryDao) DeleteAllRTsByUserId(userId string) {
	for j := len(i.rts); j >= 0; j-- {
		if i.rts[j].UserId == userId {
			// swap the current element with the last one and make i.rts smaller by 1
			i.rts[j], i.rts[len(i.rts) - 1] = i.rts[len(i.rts) - 1], i.rts[j]
			i.rts = i.rts[:len(i.rts) - 1]
		}
	}
}

func (i *InMemoryDao) GetRTModel(rt string) (RtModel, error) {
	for _, tokenModel := range(i.rts) {
		if tokenModel.Token == rt {
			return tokenModel, nil
		}
	}

	return RtModel{}, errors.New("Not found")
}

