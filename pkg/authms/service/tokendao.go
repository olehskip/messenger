package service

import (
	"github.com/olegskip/messenger/pkg/models"

	"errors"
)

type ITokenDAO interface {
	AddNewRT(rt models.RTModel)
	DeleteAllRTsByUsername(username string)

	GetRTModel(rt string) (models.RTModel, error)
}

type InMemoryDAO struct {	
	rts []models.RTModel
}

func (i *InMemoryDAO) AddNewRT(rt models.RTModel) {
	i.rts = append(i.rts, rt)
}

func (i *InMemoryDAO) DeleteAllRTsByUsername(username string) {
	for j := len(i.rts); j >= 0; j-- {
		if i.rts[j].Username == username {
			// swap the current element with the last one and make i.rts smaller by 1
			i.rts[j], i.rts[len(i.rts) - 1] = i.rts[len(i.rts) - 1], i.rts[j]
			i.rts = i.rts[:len(i.rts) - 1]
		}
	}
}

func (i *InMemoryDAO) GetRTModel(rt string) (models.RTModel, error) {
	for _, tokenModel := range(i.rts) {
		if tokenModel.Token == rt {
			return tokenModel, nil
		}
	}

	return models.RTModel{}, errors.New("Not found")
}

