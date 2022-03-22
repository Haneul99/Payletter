package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
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
	ErrCode  int  `json:"errCode"`
	Contents User `json:"contents"`
}

func LoadPersonalData(c echo.Context) error {
	resLoadPersonalData := ResLoadPersonalData{}
	reqLoadPeronsalData := ReqLoadPeronsalData{}
	if err := c.Bind(&reqLoadPeronsalData); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_LOAD_PERSONAL_DATA_REQUEST_BINDING)
	}

	if isValid, errCode, err := util.IsValidAccessToken(reqLoadPeronsalData.AccessToken, reqLoadPeronsalData.Username); !isValid || err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	if personalData, errCode, err := getPersonalData(reqLoadPeronsalData); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	} else {
		resLoadPersonalData.Contents = personalData
	}

	resLoadPersonalData.ErrCode = 0
	return c.JSON(http.StatusOK, resLoadPersonalData)
}

func getPersonalData(user ReqLoadPeronsalData) (User, int, error) {
	personalData := User{}
	query := fmt.Sprintf("SELECT username, email FROM USER WHERE username = \"%s\"", user.Username)
	if err := util.GetDB().QueryRow(query).Scan(&personalData.Username, &personalData.Email); err != nil {
		return personalData, handleError.ERR_JWT_GET_DB, err
	}
	return personalData, 0, nil
}
