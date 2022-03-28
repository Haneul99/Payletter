package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"database/sql"
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
	if status, errCode, err := loginCheckParamPassword(reqLogin); err != nil {
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
	resLogin.ErrCode = handleError.SUCCESS
	resLogin.Username = reqLogin.Username
	return c.JSON(http.StatusOK, resLogin)
}

func loginCheckParamPassword(reqLogin ReqLogin) (int, int, error) {
	query := fmt.Sprintf("SELECT password FROM USER WHERE username = \"%s\"", reqLogin.Username)
	var password = ""
	err := util.GetDB().QueryRow(query).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusBadRequest, handleError.ERR_LOGIN_SQL_NO_RESULT, err
		}
		return http.StatusInternalServerError, handleError.ERR_JWT_GET_DB, err
	}
	if password != reqLogin.Password {
		return http.StatusBadRequest, handleError.ERR_LOGIN_INCORRECT_PASSWORD, errors.New("ERR_LOGIN_INCORRECT_PASSWORD")
	}
	return http.StatusOK, handleError.SUCCESS, nil
}

func insertUserAccessToken(token, username string) (int, error) {
	query := fmt.Sprintf("UPDATE USER SET accessToken = \"%s\" WHERE username = \"%s\"", token, username)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		return handleError.ERR_JWT_GET_DB, err
	}
	return handleError.SUCCESS, nil
}
