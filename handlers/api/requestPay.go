package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"encoding/json"
	"fmt"
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
	ErrCode   int    `json:"errCode"`
	OnlineURL string `json:"online_url"`
	MobileURL string `json:"mobile_url"`
}

func RequestPay(c echo.Context) error {
	reqRequestPay := ReqRequestPay{}
	if err := c.Bind(&reqRequestPay); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_REQUEST_PAY_REQUEST_BINDING)
	}

	if isValid, errCode, err := util.IsValidAccessToken(reqRequestPay.AccessToken, reqRequestPay.Username); !isValid || err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	respBody, errCode, err := util.RequestPayletterAPI(http.MethodPost, "v1.0/payments/request", reqRequestPay.Username, reqRequestPay.Price, reqRequestPay.Platform, reqRequestPay.Membership)
	if err != nil {
		handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}
	fmt.Println(string(respBody))

	resRequestPay := ResRequestPay{}
	err = json.Unmarshal(respBody, &resRequestPay)
	if err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_REQUEST_PAY_JSON_UNMARSHAL)
	}
	resRequestPay.ErrCode = 0

	return c.JSON(http.StatusOK, resRequestPay)
}
