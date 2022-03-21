package handlers

import (
	"Haneul99/Payletter/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqRequestPay struct {
	OTTserviceId int    `json:"OTTserviceId"`
	Platform     string `json:"platform"`
	Membership   string `json:"membership"`
	Price        int    `json:"price"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
}

type ResRequestPay struct {
	Success bool `json:"success"`
}

func RequestPay(c echo.Context) error {
	reqRequestPay := ReqRequestPay{}
	if err := c.Bind(&reqRequestPay); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_REQUEST_BINDING})
	}

	if isValid, err := util.IsValidAccessToken(reqRequestPay.AccessToken, reqRequestPay.Username); !isValid || err != nil {
		return c.JSON(http.StatusUnauthorized, ResFail{ErrCode: false, Message: ERR_ACCESSTOKEN})
	}

	requestPayPayletter()

	resRequestPay := ResRequestPay{}
	resRequestPay.Success = true

	return c.JSON(http.StatusOK, resRequestPay)
}

// Payletter 결제요청 api 호출
func requestPayPayletter() {

}
