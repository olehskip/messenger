package types

type UserId int64
type Token string
type CredentialsType string
type Credentials struct {
	Password string
	Phone    string
	Email    string
}
type SearchUserQuery string
type AuthService interface {
	Register(CredentialsType, Credentials)
	Login(CredentialsType, Credentials) Token
	Logout(Token)
	GetIdentity(Token) UserId
	UpdateCredentials(Token, CredentialsType, Credentials)
}

type UserService interface {
	GetUser(UserId)
	SearchUsers(SearchUserQuery)
	CreateUser()
	EditUser(UserId)
}
type AuthServiceImpl struct {
	userService UserService
}
type UserServiceImpl struct {
}
