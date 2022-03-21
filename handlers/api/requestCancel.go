package handlers

import (
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
	Success    bool `json:"success"`
	CanceledId int  `json:"canceledId"`
}

func RequestCancel(c echo.Context) error {
	reqRequestCancel := ReqRequestCancel{}
	resRequestCancel := ResRequestCancel{}

	if err := c.Bind(&reqRequestCancel); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_REQUEST_BINDING})
	}

	if isValid, err := util.IsValidAccessToken(reqRequestCancel.AccessToken, reqRequestCancel.Username); !isValid || err != nil {
		return c.JSON(http.StatusUnauthorized, ResFail{ErrCode: false, Message: ERR_ACCESSTOKEN})
	}

	//Payletter cancel API 호출
	requestCancelPayletter()

	resRequestCancel.Success = true
	resRequestCancel.CanceledId = reqRequestCancel.SubscribedOTTId
	return c.JSON(http.StatusOK, resRequestCancel)
}

func requestCancelPayletter() {

}
