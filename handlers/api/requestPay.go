package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
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
	ErrCode int `json:"errCode"`
}

func RequestPay(c echo.Context) error {
	reqRequestPay := ReqRequestPay{}
	if err := c.Bind(&reqRequestPay); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_REQUEST_PAY_REQUEST_BINDING)
	}

	if isValid, errCode, err := util.IsValidAccessToken(reqRequestPay.AccessToken, reqRequestPay.Username); !isValid || err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	requestPayPayletter()

	resRequestPay := ResRequestPay{}
	resRequestPay.ErrCode = 0

	return c.JSON(http.StatusOK, resRequestPay)
}

// Payletter 결제요청 api 호출
func requestPayPayletter() {

}
