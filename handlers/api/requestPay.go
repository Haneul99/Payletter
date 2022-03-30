package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"database/sql"
	"encoding/json"
	"errors"
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
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Token     int    `json:"token"`
	OnlineURL string `json:"online_url"`
	MobileURL string `json:"mobile_url"`
}

func RequestPay(c echo.Context) error {
	reqRequestPay := ReqRequestPay{}

	// Bind
	if err := c.Bind(&reqRequestPay); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_REQUEST_PAY_REQUEST_BINDING)
	}

	// CheckParam
	if errCode, err := util.IsValidAccessToken(reqRequestPay.AccessToken, reqRequestPay.Username); err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	if status, errCode, err := requestPayCheckParam(reqRequestPay.Username, reqRequestPay.OTTserviceId); err != nil {
		return handleError.ReturnResFail(c, status, err, errCode)
	}

	// Process
	respBody, errCode, err := util.RequestPayAPI(reqRequestPay.Username, reqRequestPay.Platform, reqRequestPay.Membership, reqRequestPay.OTTserviceId, reqRequestPay.Price)
	if err != nil {
		handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	fmt.Println(string(respBody))

	resRequestPay := ResRequestPay{}
	if err := json.Unmarshal(respBody, &resRequestPay); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_REQUEST_PAY_JSON_UNMARSHAL)
	}

	// Return
	resRequestPay.Code = 0
	return c.JSON(http.StatusOK, resRequestPay)
}

// 중복결제 하면 안되니까 구독 정보 중에 해당 유저 + 해당 serviceId가 일치하는게 있는지 확인
func requestPayCheckParam(username string, OTTserviceId int) (int, int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = \"%s\" && OTTserviceId = %d && canceled = 0", "subscribedServices", username, OTTserviceId)
	exist := 0
	if err := util.GetDB().QueryRow(query).Scan(&exist); err != nil {
		if err == sql.ErrNoRows {
			return http.StatusInternalServerError, handleError.ERR_REQUEST_PAY_SQL_NO_RESULT, err
		}
		return http.StatusInternalServerError, handleError.ERR_REQUEST_PAY_GET_DB, err
	}
	if exist != 0 {
		return http.StatusForbidden, handleError.ERR_REQUEST_PAY_ALREADY_PAID, errors.New("ERR_REQUEST_PAY_ALREADY_PAID")
	}
	return http.StatusOK, handleError.SUCCESS, nil
}
