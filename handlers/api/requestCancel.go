package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqRequestCancel struct {
	SubscribedOTTId int    `json:"subscribedOTTId"`
	OTTserviceId    int    `json:"OTTserviceId"`
	Platform        string `json:"platform"`
	Membership      string `json:"membership"`
	Price           int    `json:"price"`
	Username        string `json:"username"`
	AccessToken     string `json:"accessToken"`
}

type ResRequestCancel struct {
	ErrCode    int `json:"errCode"`
	CanceledId int `json:"canceledId"`
}

func RequestCancel(c echo.Context) error {
	reqRequestCancel := ReqRequestCancel{}
	resRequestCancel := ResRequestCancel{}

	if err := c.Bind(&reqRequestCancel); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_REQUEST_CANCEL_REQUEST_BINDING)
	}

	if isValid, errCode, err := util.IsValidAccessToken(reqRequestCancel.AccessToken, reqRequestCancel.Username); !isValid || err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	//Payletter cancel API 호출
	requestCancelPayletter()

	resRequestCancel.ErrCode = 0
	resRequestCancel.CanceledId = reqRequestCancel.SubscribedOTTId
	return c.JSON(http.StatusOK, resRequestCancel)
}

func requestCancelPayletter() {

}
