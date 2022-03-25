package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqLogout struct {
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}

type ResLogout struct {
	ErrCode int `json:"errCode"`
}

func Logout(c echo.Context) error {
	reqLogout := ReqLogout{}
	resLogout := ResLogout{}

	// Bind
	if err := c.Bind(&reqLogout); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_LOGOUT_REQEUST_BINDING)
	}

	// CheckParam
	// 해당 accessToken이 유효한지 검사
	// 해당 accessToken이 DB에 저장된 것과 동일한지 검사
	if errCode, err := util.IsValidAccessToken(reqLogout.AccessToken, reqLogout.Username); err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	// Process
	if errCode, err := deleteUserAccessToken(reqLogout.Username); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// Return
	resLogout.ErrCode = 0
	return c.JSON(http.StatusOK, resLogout)
}

// DB에서 AccessToken 삭제
func deleteUserAccessToken(username string) (int, error) {
	query := fmt.Sprintf("UPDATE USER SET accessToken = \"\" WHERE username = \"%s\"", username)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		return handleError.ERR_JWT_GET_DB, err
	}
	return 0, nil
}
