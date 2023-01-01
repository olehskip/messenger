package service

import (
	"time"
	"fmt"
	"crypto/sha256"
	"reflect"
)

type ITokensGenerator interface {
	GenerateToken(userUuid string) Token
	GenerateHashedToken(token Token) string
}

type Token struct {
	UserUuid string
	IssueTimestamp time.Time
	ExpiryTimestamp time.Time
	Secret string
	TokenId int
	GeneratorId string
}

type TokensGenerator struct {
	secret string
	lastTokenId int
	generatorId string
	tokenDuration time.Duration
}

func (t *TokensGenerator) GenerateToken(userUuid string) Token {
	return Token {
		UserUuid: userUuid,
		IssueTimestamp: time.Now(),
		ExpiryTimestamp: time.Now().Add(t.tokenDuration),
		Secret: t.secret,
		TokenId: t.lastTokenId,
		GeneratorId: t.generatorId,
	}
}

func (t *TokensGenerator) GenerateHashedToken(token Token) string {
	v := reflect.ValueOf(token)
	tokenStr := ""
	for i := 0; i < v.NumField(); i++ {
		tokenStr += fmt.Sprintf("%v", v.Field(i).Interface())
	}

	t.lastTokenId++
	return fmt.Sprintf("%x", sha256.Sum256([]byte(tokenStr)))
}

func NewTokensGenerator(secret string, tokenDuration time.Duration) *TokensGenerator {
	return &TokensGenerator {
		secret: secret,
		lastTokenId: 0,
		generatorId: time.Now().String(),
		tokenDuration: tokenDuration,
	}
}

