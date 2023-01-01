package dal

import (
	"errors"
	"time"
)

type ITokensPairDao interface {
	AddTokensPair(tokensPair TokensPairModel)
	// IsRefreshTokenRevoked(hashedRefreshToken string) bool, err
	
	// if hRefreshToken wasn't found or it's revoked or it's expired error is returned
	UpdateHAccessToken(hRefreshToken string, newHAccessToken string) error
	// GetUserIdByHRefreshToken(hRefreshToken string) (int, error)
	
	GetUserUuidByHRefreshToken(hRefreshToken string) (string, bool, error)
	GetUserUuidByHAcessToken(hAccessToken string) (string, bool, error)
	
	IsHRefrehTokenValid(hRefreshToken string) error
	RevokeTokensPair(hRefreshToken string) error

	isTokensPairValid(tokensPair TokensPairModel) error
}

type InMemoryDao struct {	
	tokensPairs []TokensPairModel
}

func (i *InMemoryDao) AddTokensPair(tokensPair TokensPairModel) {
	i.tokensPairs = append(i.tokensPairs, tokensPair)
}

func (i *InMemoryDao) UpdateHAccessToken(hRefreshToken string, newHAccessToken string) error {
	for _, tokensPair := range i.tokensPairs {
		if tokensPair.HRefreshToken == hRefreshToken {
			if err := i.isTokensPairValid(tokensPair); err != nil {
				return err
			}
			tokensPair.HAccessToken = newHAccessToken
			return nil
		}
	}

	return errors.New("hRefreshToken wasn't found")
}

func (i *InMemoryDao) GetUserUuidByHRefreshToken(hRefreshToken string) (string, bool, error) {
	for _, tokensPair := range i.tokensPairs {
		if tokensPair.HRefreshToken == hRefreshToken {
			return tokensPair.UserUuid, tokensPair.IsRevoked, nil
		}
	}

	return "", true, errors.New("hAccessToken wasn't found")
}

func (i *InMemoryDao) GetUserUuidByHAcessToken(hAccessToken string) (string, bool, error) {
	for _, tokensPair := range i.tokensPairs {
		if tokensPair.HAccessToken == hAccessToken {

			return tokensPair.UserUuid, tokensPair.IsRevoked, nil
		}
	}

	return "", true, errors.New("hAccessToken wasn't found")
}

func (i *InMemoryDao) RevokeTokensPair(hRefreshToken string) error {
	for j := 0; j < len(i.tokensPairs); j++ {
		if i.tokensPairs[j].HRefreshToken == hRefreshToken {
			if i.tokensPairs[j].IsRevoked {
				return errors.New("tokens pair is already revoked")
			}
			i.tokensPairs[j].IsRevoked = true
			return nil
		}
	}

	return errors.New("hRefreshToken wasn't found")
}


func (i *InMemoryDao) IsHRefrehTokenValid(hRefreshToken string) error {
	for _, tokensPair := range i.tokensPairs {
		if tokensPair.HRefreshToken == hRefreshToken {
			if err := i.isTokensPairValid(tokensPair); err != nil {
				return err
			}

			return  nil
		}
	}

	return errors.New("hRefreshToken wasn't found")
}

func (i *InMemoryDao) isTokensPairValid(tokensPair TokensPairModel) error {
	if tokensPair.IsRevoked {
		return errors.New("tokens pair is revoked")
	}

	if time.Now().After(tokensPair.RefreshTokenExpiryTimestamp) {
		return errors.New("refresh token is expired")
	}

	if time.Now().After(tokensPair.AccessTokenExpiryTimestamp) {
		return errors.New("access token is expired")
	}
	
	return nil;
}
