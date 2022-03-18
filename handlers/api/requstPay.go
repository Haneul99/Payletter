package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqRequestPay struct {
	OTTserviceId int    `json: "OTTserviceId"`
	Platform     string `json: "platform`
	Membership   string `json: "membership"`
	Price        int    `json: "price"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
}

type ResRequestPay struct {
	Success bool `json: "success`
}

func ReqestPay(c echo.Context) error {
	reqRequestPay := ReqRequestPay{}
	if err := c.Bind(&reqRequestPay); err != nil {
		return err
	}
	if err := checkPayRequestValidity(); err != nil {
		return c.JSON(http.StatusBadRequest, "request failure")
	}
	callPayletter()

	resRequestPay := ResRequestPay{}
	resRequestPay.Success = true

	return c.JSON(http.StatusOK, resRequestPay)
}

// Payletter 결제요청 api 호출
func callPayletter() {

}

func checkPayRequestValidity() error {
	return errors.New("failure") // failure test
}
