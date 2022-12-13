package tokenfactory

import (
	"crypto/sha256"
	"fmt"
	"reflect"

	. "github.com/olegskip/messenger/pkg/authms/types"
)

type TokenProperties struct {
	UserIdBruh UserId
	IssuedAt   float32
	ExpiresAt  float32
	TokenType  string
	Device     string
}

type TokenCreator interface {
	CreateToken(TokenProperties) Token
}
type TokenCreatorImpl struct {
}

func (tpr *TokenProperties) String() string {
	v := reflect.ValueOf(*tpr)

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	res := ""
	fmt.Println(values)
	for i := 0; i < v.NumField(); i++ {
		res += fmt.Sprintf("%v", values[i])
	}
	return res
}

func (tck TokenCreatorImpl) CreateToken(tps TokenProperties) Token {
	sum := sha256.Sum256([]byte(tps.String()))
	fmt.Printf("%x", sum)
	return Token(fmt.Sprintf("%x", sum))
}
