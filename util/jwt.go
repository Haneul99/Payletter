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
	return tk, 0, nil
}

// 유효한 accessToken인지 검사
func IsValidAccessToken(accessToken, username string) (bool, int, error) {
	isValid, errCode, err := decodeJWT(accessToken)
	if !isValid || err != nil {
		return false, errCode, err
	}

	isStored, errCode, err := isStoredAccessToken(accessToken, username)
	if !isStored || err != nil {
		return false, errCode, err
	}

	return true, 0, nil
}

// JWT Token 검증
func decodeJWT(accessToken string) (bool, int, error) {
	if accessToken == "" {
		return false, handleError.ERR_JWT_NULL_ACCESSTOKEN, errors.New("ERR_JWT_NULL_ACCESSTOKEN")
	}

	claims := &AccessTokenClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(tk *jwt.Token) (interface{}, error) {
		return signKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, handleError.ERR_JWT_INVALID_ACCESSTOKEN, err
		}
	}
	return true, 0, nil
}

// 저장되어 있는 accessToken이 일치하는지 검사
func isStoredAccessToken(accessToken, username string) (bool, int, error) {
	query := fmt.Sprintf("SELECT accessToken From USER WHERE username = \"%s\"", username)
	var storedTK = ""
	err := GetDB().QueryRow(query).Scan(&storedTK)
	if err != nil {
		return false, handleError.ERR_JWT_GET_DB, err
	}
	if storedTK != accessToken {
		return false, handleError.ERR_JWT_INCORRECT_ACCESSTOKEN, nil
	}
	return true, 0, nil
}
