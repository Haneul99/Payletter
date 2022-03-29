package util

import (
	handleError "Haneul99/Payletter/handlers/error"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTtoken struct {
	Token string `json:"token"`
}

type AccessTokenClaims struct {
	Username string
	jwt.StandardClaims
}

var signKey = []byte(ServerConfig.GetStringData("SECRET_KEY"))

// JWT Token 생성
func CreateJWTAccessToken(username string) (string, int, error) {
	claims := &AccessTokenClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := accessToken.SignedString(signKey)
	if err != nil {
		return "", handleError.ERR_JWT_CREATE_ACCESSTOKEN, err
	}
	return tk, handleError.SUCCESS, nil
}

// 유효한 accessToken인지 검사
func IsValidAccessToken(accessToken, username string) (int, error) {
	if errCode, err := decodeJWT(accessToken); err != nil {
		return errCode, err
	}

	if errCode, err := isStoredAccessToken(accessToken, username); err != nil {
		return errCode, err
	}

	return handleError.SUCCESS, nil
}

// JWT Token 검증
func decodeJWT(accessToken string) (int, error) {
	if accessToken == "" {
		return handleError.ERR_JWT_NULL_ACCESSTOKEN, errors.New("ERR_JWT_NULL_ACCESSTOKEN")
	}

	claims := &AccessTokenClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(tk *jwt.Token) (interface{}, error) {
		return signKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return handleError.ERR_JWT_INVALID_ACCESSTOKEN, err
		}
		return handleError.ERR_JWT_ACCESSTOKEN_EXPIRED, err
	}
	return handleError.SUCCESS, nil
}

// 저장되어 있는 accessToken이 일치하는지 검사
func isStoredAccessToken(accessToken, username string) (int, error) {
	query := fmt.Sprintf("SELECT accessToken From USER WHERE username = \"%s\"", username)
	var storedTK = ""
	err := GetDB().QueryRow(query).Scan(&storedTK)
	if err != nil {
		return handleError.ERR_JWT_GET_DB, err
	}
	if storedTK != accessToken {
		return handleError.ERR_JWT_INCORRECT_ACCESSTOKEN, errors.New("ERR_JWT_INCORRECT_ACCESSTOKEN")
	}
	return handleError.SUCCESS, nil
}
