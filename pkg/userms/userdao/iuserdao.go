package userdao

import (
	"github.com/olegskip/messenger/pkg/models"
)

type IUserDao interface {
	Connect() bool
	Disconnect() bool
	GetUserById(userId string) *usermodel.UserModel
	GetUserByUsername(username string) *usermodel.UserModel
	GetUsersByName(name string) []usermodel.UserModel
	CreateUser(userModel usermodel.UserModel) bool 
}

