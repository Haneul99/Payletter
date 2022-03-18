package handlers

import (
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LogoutInfo struct {
	Username    string `json: "username"`
	AccessToken string `json: "accessToken"`
}

type ResLogout struct {
	Success bool `json: "success"`
}

func Logout(c echo.Context) error {
	logoutInfo := LogoutInfo{}
	resLogout := ResLogout{}
	if err := c.Bind(&logoutInfo); err != nil {
		fmt.Println(err)
		return err
	}

	// 해당 accessToken이 유효한지 검사
	// 해당 accessToken이 DB에 저장된 것과 동일한지 검사
	isValid, err := util.IsValidAccessToken(logoutInfo.AccessToken, logoutInfo.Username)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !isValid {
		return c.JSON(http.StatusOK, "invalid accessToken")
	}

	err = deleteUserAccessToken(logoutInfo.Username)
	if err != nil {
		fmt.Println(err)
		return err
	}

	resLogout.Success = true
	return c.JSON(http.StatusOK, resLogout)
}

// DB에서 AccessToken 삭제
func deleteUserAccessToken(username string) error {
	query := fmt.Sprintf("UPDATE USER SET accessToken = \"\" WHERE username = \"%s\"", username)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
