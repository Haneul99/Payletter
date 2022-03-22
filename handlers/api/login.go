package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResLogin struct {
	ErrCode     int    `json:"errCode"`
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}

func Login(c echo.Context) error {
	resLogin := ResLogin{}
	reqLogin := ReqLogin{}
	if err := c.Bind(&reqLogin); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_LOGIN_REQUEST_BINDING)
	}

	// 비밀번호 오류
	pwd, errCode, err := getUserPwd(reqLogin)
	if err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}
	if pwd != reqLogin.Password {
		return handleError.ReturnResFail(c, http.StatusBadRequest, errors.New("ERR_LOGIN_INCORRECT_PASSWORD"), errCode)
	}

	accessToken, errCode, err := util.CreateJWTAccessToken(reqLogin.Username)
	if err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// DB에 AccessToken 삽입
	if errCode, err := insertUserAccessToken(accessToken, reqLogin.Username); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	resLogin.AccessToken = accessToken
	resLogin.ErrCode = 0
	resLogin.Username = reqLogin.Username

	return c.JSON(http.StatusOK, resLogin)
}

func getUserPwd(reqLogin ReqLogin) (string, int, error) {
	query := fmt.Sprintf("SELECT password FROM USER WHERE username = \"%s\"", reqLogin.Username)
	var password = ""
	err := util.GetDB().QueryRow(query).Scan(&password)
	if err != nil {
		return "", handleError.ERR_JWT_GET_DB, err
	}
	return password, 0, nil
}

func insertUserAccessToken(token, username string) (int, error) {
	query := fmt.Sprintf("UPDATE USER SET accessToken = \"%s\" WHERE username = \"%s\"", token, username)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		return handleError.ERR_JWT_GET_DB, err
	}
	return 0, nil
}
