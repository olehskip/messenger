package main

import (
	"fmt"
	"github.com/olegskip/messenger/pkg/userms/userdao"
	// "github.com/olegskip/messenger/pkg/models"
)

func main() {
	testv := userdao.ScyllaUserDao{}
	testv.Connect()
	// fmt.Println(testv.GetUserById("3447ff60-ef45-4fb4-ac56-699752ed642c"))
	// fmt.Println(testv.GetUserByUsername("skip"))
	// fmt.Println(testv.GetUsersByName("skip"))
	// fmt.Println(testv.CreateUser(usermodel.UserModel{Username: "Skela", Name: "Sasha", Bio: "chel1"}));
	// fmt.Println(testv.UpdateUser(usermodel.UserModel{Id: "4c090a00-09ce-4360-947a-248d5cb2aa3a", Bio: "chel"}))
	fmt.Println(testv.DeleteUser("4c090a00-09ce-4360-947a-248d5cb2aa3a"))
}

