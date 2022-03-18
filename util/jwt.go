package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTtoken struct {
	Token string `json: "token"`
}

type AccessTokenClaims struct {
	Username string
	jwt.StandardClaims
}

var signKey = []byte(ServerConfig.GetStringData("SECRET_KEY"))

// JWT Token 생성
func CreateJWTAccessToken(username string) (string, error) {
	claims := &AccessTokenClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := accessToken.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tk, nil
}

// JWT Token 검증
func DecodeJWT(accessToken string) (bool, error) {
	if accessToken == "" {
		fmt.Println("token is null")
		return false, nil
	}

	claims := &AccessTokenClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(tk *jwt.Token) (interface{}, error) {
		return signKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println("jwt token is invalid")
			return false, err
		}
	}
	fmt.Println("decode succes", claims.Username)
	return true, nil
}

// 저장되어 있는 accessToken이 일치하는지 검사
func IsStoredAccessToken(accessToken, username string) (bool, error) {
	query := fmt.Sprintf("SELECT accessToken From USER WHERE username = \"%s\"", username)
	var storedTK = ""
	err := GetDB().QueryRow(query).Scan(&storedTK)
	if err != nil {
		return false, err
	}
	if storedTK != accessToken {
		return false, nil
	}
	return true, nil
}

// 유효한 accessToken인지 검사
func IsValidAccessToken(accessToken, username string) (bool, error) {
	isValid, err := DecodeJWT(accessToken)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, nil
	}

	isStored, err := IsStoredAccessToken(accessToken, username)
	if err != nil {
		return false, err
	}
	if !isStored {
		return false, nil
	}

	return true, nil
}
