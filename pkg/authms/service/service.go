package service

import (
	"errors"

	"github.com/olegskip/messenger/pkg/authms/dal"
)

type IAuthService interface {
	GetNewRefreshToken(credentials CredentialsDto) (hTokens HTokensPairDto, err error)
	ExchangeRefreshToken(oldHRefreshToken string) (hTokens HTokensPairDto, err error)
	RevokeRefreshToken(oldHRefreshToken string) (err error)
	GetUserUuid(hRefreshToken string) (uuid string, isRevoked bool, err error)
}

type AuthService struct {
	tokensPairDao dal.ITokensPairDao
	refreshTokensGenerator ITokensGenerator
	accessTokensGenerator ITokensGenerator
}

func (a *AuthService) GetNewRefreshToken(credentials CredentialsDto) (hTokens HTokensPairDto, err error){
	if credentials.Password == "pass" {
		hTokens.HashedRefreshToken, hTokens.HashedAccessToken = a.generateAndSaveTokensPair(credentials.UserUuid)
		err = nil
	} else {
		err = errors.New("invalid username or password")
	}
	
	return hTokens, err
}

func (a *AuthService) ExchangeRefreshToken(oldHRefreshToken string) (hTokens HTokensPairDto, err error){
	userUuid, isRevoked, findingErr := a.tokensPairDao.GetUserUuidByHRefreshToken(oldHRefreshToken)

	if isRevoked {
		return HTokensPairDto{}, errors.New("old refresh token is revoked")
	}

	if findingErr != nil {
		return HTokensPairDto{}, findingErr
	}

	// revoke the old refresh token to prevent creating multiple access tokens
	// if oldRefreshToken is already revoked we can't give an access token
	if revokingErr := a.tokensPairDao.RevokeTokensPair(oldHRefreshToken); revokingErr != nil {
		return HTokensPairDto{}, revokingErr
	}

	hTokens.HashedRefreshToken, hTokens.HashedAccessToken = a.generateAndSaveTokensPair(userUuid)
	
	return hTokens, nil 
}

func (a *AuthService) RevokeRefreshToken(oldHRefreshToken string) (err error) {
	if revokingErr := a.tokensPairDao.RevokeTokensPair(oldHRefreshToken); revokingErr != nil {
		return revokingErr 
	}
	
	return nil
}

func (a *AuthService) GetUserUuid(hAccessToken string) (uuid string, isRevoked bool, err error) {
	return a.tokensPairDao.GetUserUuidByHAcessToken(hAccessToken)
}

func (a *AuthService) generateAndSaveTokensPair(userUuid string) (hRefreshToken string, hAccessToken string) {
	refreshToken := a.refreshTokensGenerator.GenerateToken(userUuid)
	hRefreshToken = a.refreshTokensGenerator.GenerateHashedToken(refreshToken)
	
	accessToken := a.accessTokensGenerator.GenerateToken(userUuid)
	hAccessToken = a.accessTokensGenerator.GenerateHashedToken(accessToken)

	a.tokensPairDao.AddTokensPair(toTokensPairModel(refreshToken, hRefreshToken, accessToken, hAccessToken))

	return hRefreshToken, hAccessToken
}

func NewAuthService(tokensPairDao dal.ITokensPairDao, refreshTokensGenerator ITokensGenerator, accessTokensGenerator ITokensGenerator) *AuthService {
	return &AuthService{
		tokensPairDao: tokensPairDao,
		refreshTokensGenerator: refreshTokensGenerator,
		accessTokensGenerator: accessTokensGenerator,
	}
}

func toTokensPairModel(refreshToken Token, hRefreshToken string, accessToken Token, hAccessToken string) dal.TokensPairModel {
	return dal.TokensPairModel {
		IsRevoked: false,
		HRefreshToken: hRefreshToken,
		RefreshTokenExpiryTimestamp: refreshToken.ExpiryTimestamp,
		HAccessToken: hAccessToken,
		AccessTokenExpiryTimestamp: accessToken.ExpiryTimestamp,
		UserUuid: refreshToken.UserUuid,
	}
}

