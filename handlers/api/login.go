package handlers

import (
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResLogin struct {
	Success     bool   `json:"success"`
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}

func Login(c echo.Context) error {
	resLogin := ResLogin{}
	reqLogin := ReqLogin{}
	if err := c.Bind(&reqLogin); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_REQUEST_BINDING})
	}

	// 비밀번호 오류
	if pwd, err := getUserPwd(reqLogin); pwd != reqLogin.Password || err != nil {
		resLogin.Success = false
		return c.JSON(http.StatusBadRequest, ResFail{ErrCode: false, Message: ERR_INCORRECT_PASSWORD})
	}

	accessToken, err := util.CreateJWTAccessToken(reqLogin.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_ACCESSTOKEN})
	}

	// DB에 AccessToken 삽입
	err = insertUserAccessToken(accessToken, reqLogin.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_INSERT_DB})
	}

	resLogin.AccessToken = accessToken
	resLogin.Success = true
	resLogin.Username = reqLogin.Username

	return c.JSON(http.StatusOK, resLogin)
}

func getUserPwd(reqLogin ReqLogin) (string, error) {
	query := fmt.Sprintf("SELECT password FROM USER WHERE username = \"%s\"", reqLogin.Username)
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
		return err
	}
	return nil
}
