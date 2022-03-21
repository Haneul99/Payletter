package handlers

import (
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LogoutInfo struct {
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}

type ResLogout struct {
	Success bool `json:"success"`
}

func Logout(c echo.Context) error {
	logoutInfo := LogoutInfo{}
	resLogout := ResLogout{}
	if err := c.Bind(&logoutInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_REQUEST_BINDING})
	}

	// 해당 accessToken이 유효한지 검사
	// 해당 accessToken이 DB에 저장된 것과 동일한지 검사
	if isValid, err := util.IsValidAccessToken(logoutInfo.AccessToken, logoutInfo.Username); !isValid || err != nil {
		return c.JSON(http.StatusBadRequest, ResFail{ErrCode: false, Message: ERR_ACCESSTOKEN})
	}

	if err := deleteUserAccessToken(logoutInfo.Username); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_DELETE_DB})
	}

	resLogout.Success = true
	return c.JSON(http.StatusOK, resLogout)
}

// DB에서 AccessToken 삭제
func deleteUserAccessToken(username string) error {
	query := fmt.Sprintf("UPDATE USER SET accessToken = \"\" WHERE username = \"%s\"", username)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		return err
	}
	return nil
}
