package handlers

import (
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqLoadPeronsalData struct {
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ResLoadPersonalData struct {
	Success  bool `json:"success"`
	Contents User `json:"contents"`
}

func LoadPersonalData(c echo.Context) error {
	resLoadPersonalData := ResLoadPersonalData{}
	reqLoadPeronsalData := ReqLoadPeronsalData{}
	if err := c.Bind(&reqLoadPeronsalData); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: 0, Message: ERR_REQUEST_BINDING})
	}

	if isValid, err := util.IsValidAccessToken(reqLoadPeronsalData.AccessToken, reqLoadPeronsalData.Username); !isValid || err != nil {
		return c.JSON(http.StatusUnauthorized, ResFail{ErrCode: 0, Message: ERR_ACCESSTOKEN})
	}

	if personalData, err := getPersonalData(reqLoadPeronsalData); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: 0, Message: ERR_ACCESSTOKEN})
	} else {
		resLoadPersonalData.Contents = personalData
	}

	resLoadPersonalData.Success = true
	return c.JSON(http.StatusOK, resLoadPersonalData)
}

func getPersonalData(user ReqLoadPeronsalData) (User, error) {
	personalData := User{}
	query := fmt.Sprintf("SELECT username, email FROM USER WHERE username = \"%s\"", user.Username)
	err := util.GetDB().QueryRow(query).Scan(&personalData.Username, &personalData.Email)
	return personalData, err
}
