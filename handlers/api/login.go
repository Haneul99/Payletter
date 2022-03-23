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

	// Bind
	if err := c.Bind(&reqLogin); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_LOGIN_REQUEST_BINDING)
	}

	// CheckParam
	if isCorrect, status, errCode, err := loginCheckParamPassword(reqLogin); !isCorrect || err != nil {
		return handleError.ReturnResFail(c, status, err, errCode)
	}

	accessToken, errCode, err := util.CreateJWTAccessToken(reqLogin.Username)
	if err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// Process
	// DB에 AccessToken 삽입
	if errCode, err := insertUserAccessToken(accessToken, reqLogin.Username); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// Return
	resLogin.AccessToken = accessToken
	resLogin.ErrCode = 0
	resLogin.Username = reqLogin.Username
	return c.JSON(http.StatusOK, resLogin)
}

func loginCheckParamPassword(reqLogin ReqLogin) (bool, int, int, error) {
	query := fmt.Sprintf("SELECT password FROM USER WHERE username = \"%s\"", reqLogin.Username)
	var password = ""
	err := util.GetDB().QueryRow(query).Scan(&password)
	if err != nil {
		return false, http.StatusInternalServerError, handleError.ERR_JWT_GET_DB, err
	}
	if password != reqLogin.Password {
		return false, http.StatusBadRequest, handleError.ERR_LOGIN_INCORRECT_PASSWORD, errors.New("ERR_LOGIN_INCORRECT_PASSWORD")
	}
	return true, http.StatusOK, 0, nil
}

func insertUserAccessToken(token, username string) (int, error) {
	query := fmt.Sprintf("UPDATE USER SET accessToken = \"%s\" WHERE username = \"%s\"", token, username)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		return handleError.ERR_JWT_GET_DB, err
	}
	return 0, nil
}
