package handlers

import (
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResLogin struct {
	Success     bool   `json:"success"`
	Username    string `json:"username"`
	AccessToken string `json:"AccessToken"`
}

func Login(c echo.Context) error {
	resLogin := ResLogin{}
	loginInfo := LoginInfo{}
	if err := c.Bind(&loginInfo); err != nil {
		return err
	}

	// 비밀번호 오류
	if pwd, err := selectUserPwd(loginInfo); pwd != loginInfo.Password || err != nil {
		resLogin.Success = false
		return c.JSON(http.StatusBadRequest, resLogin)
	}

	accessToken, err := util.CreateJWTAccessToken(loginInfo.Username)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, resLogin)
	}
	resLogin.AccessToken = accessToken
	resLogin.Success = true
	resLogin.Username = loginInfo.Username

	// DB에 AccessToken 삽입
	err = insertUserAccessToken(resLogin.AccessToken, resLogin.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, resLogin)
	}

	return c.JSON(http.StatusOK, resLogin)
}

func selectUserPwd(loginInfo LoginInfo) (string, error) {
	query := fmt.Sprintf("SELECT password FROM USER WHERE username = \"%s\"", loginInfo.Username)
	var password = ""
	err := util.GetDB().QueryRow(query).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func insertUserAccessToken(token, username string) error {
	query := fmt.Sprintf("UPDATE USER SET accessToken = \"%s\" WHERE username = \"%s\"", token, username)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
