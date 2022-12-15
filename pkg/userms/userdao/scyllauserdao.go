package userdao

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/olegskip/messenger/pkg/models"
)

type ScyllaUserDao struct {
	session gocqlx.Session
}

func (s *ScyllaUserDao) Connect() bool {
	cluster := gocql.NewCluster("172.17.0.2")
	cluster.Keyspace = "messenger"

	var err error
	s.session, err = gocqlx.WrapSession(cluster.CreateSession())
	return err == nil
}

func (s *ScyllaUserDao) Disconnect() {
	s.session.Close()
}

// Finds a user by his id and returns his UserModel*, if a user was not found nil is returned
func (s *ScyllaUserDao) GetUserById(userId string) *usermodel.UserModel {
	var res usermodel.UserModel
	q := qb.Select("users").Where(qb.Eq("id")).Query(s.session).BindMap(qb.M {
		"id": userId,
	})
	
	err := q.GetRelease(&res)
	if err != nil {
		return nil 
	} else {
		return &res
	}
}

// Finds a user by his username and returns his UserModel*, if a user was not found nil is returned
func (s *ScyllaUserDao) GetUserByUsername(username string) *usermodel.UserModel {
	var res usermodel.UserModel
	q := qb.Select("users").Where(qb.Eq("username")).Query(s.session).BindMap(qb.M {
		"username": username,
	})

	err := q.GetRelease(&res)
	if err != nil {
		return nil 
	} else {
		return &res
	}
}

// Find users by name using %like% and returns a slice of their UserModel
func (s *ScyllaUserDao) GetUsersByName(name string) []usermodel.UserModel {
	// Cassandra can't do filtering usinglike by default
	// Additional cassandra configuration may be neededon

	var res []usermodel.UserModel
	q := qb.Select("users").Where(qb.Like("username")).Query(s.session).BindMap(qb.M {
		"name": "%" + name + "%",
	})

	err := q.GetRelease(&res)
	if err != nil {
		return nil 
	} else {
		return res
	}	
}

// Creates a user using the info in userModel, if in userModel id is not provided then it is generated
// If id is provided but is not unique then false is returned
func (s *ScyllaUserDao) CreateUser(userModel usermodel.UserModel) bool {
	// TODO: stop manually list all columns
	q := qb.Insert("users").Columns("username", "name", "bio")
	if userModel.Id == "" {
		// if id is not in the model then generate it by using cql function uuid()
		q = q.FuncColumn("id", qb.Fn("uuid"))

	} else {
		q = q.Columns("id")
	}
	

	// Cassandra doesn't have values by default so we need to pass timestamp manually
	// though gocqlx can handle functions(qb.Func) it can't handle nested functions
	// so we use LitColumn function because it doesn't wrap the column name with quotes
	q = q.LitColumn("reg_timestamp", "toTimestamp(now())")

	qFinal := q.Query(s.session).BindStruct(userModel)
	
	err := qFinal.ExecRelease()
	return err == nil
}

// Updates the inforamtion about the user with given id
// TODO: stop setting fields to empty string
func (s *ScyllaUserDao) UpdateUser(userModel usermodel.UserModel) bool {
	// TODO: stop manually list all columns
	q := qb.Update("users").Set("username", "name", "bio").Where(qb.Eq("id")).Query(s.session).BindStruct(userModel)
	
	err := q.ExecRelease()
	return err == nil
}

// Delete the user with given id, even if the user doesn't exist(but no error occurs) true is returned
func (s *ScyllaUserDao) DeleteUser(userId string) bool {	
	q := qb.Delete("users").Where(qb.Eq("id")).Query(s.session).Bind(userId)
	
	err := q.ExecRelease()
	return err == nil
}

