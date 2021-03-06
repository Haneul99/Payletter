package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqRequestCancel struct {
	SubscribedServiceId int    `json:"subscribedServiceId"`
	Username            string `json:"username"`
	AccessToken         string `json:"accessToken"`
}

type ResPayletterRequestCancel struct {
	ErrCode    int    `json:"code"`
	Message    string `json:"message"`
	TID        string `json:"tid"`
	CID        string `json:"cid"`
	Amount     int    `json:"amount"`
	CancelDate string `json:"cancel_date"`
}

func RequestCancel(c echo.Context) error {
	reqRequestCancel := ReqRequestCancel{}
	resPayletterRequestCancel := ResPayletterRequestCancel{}

	if err := c.Bind(&reqRequestCancel); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_REQUEST_CANCEL_REQUEST_BINDING)
	}

	if errCode, err := util.IsValidAccessToken(reqRequestCancel.AccessToken, reqRequestCancel.Username); err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	tid, price, pgcode, errCode, err := getPayInfo(reqRequestCancel)
	if err != nil {
		handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	respBody, errCode, err := util.RequestCancelAPI(reqRequestCancel.Username, pgcode, tid, price)
	if err != nil {
		handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	if err := json.Unmarshal(respBody, &resPayletterRequestCancel); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_REQUEST_CANCEL_JSON_UNMARSHAL)
	}

	// 결제 취소가 성공했다면
	if resPayletterRequestCancel.ErrCode == 0 {
		if errCode, err := deletePayInfo(reqRequestCancel); err != nil {
			return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
		}
	}

	// Return
	return c.JSON(http.StatusOK, resPayletterRequestCancel)
}

func getPayInfo(req ReqRequestCancel) (string, int, string, int, error) {
	query := fmt.Sprintf("SELECT tid, price, pgcode FROM subscribedServices WHERE SubscribedServiceId = %d", req.SubscribedServiceId)
	tid := ""
	price := 0
	pgcode := ""
	if err := util.GetDB().QueryRow(query).Scan(&tid, &price, &pgcode); err != nil {
		if err == sql.ErrNoRows {
			return "", 0, "", handleError.ERR_REQUEST_CANCEL_SQL_NO_RESULT, err
		}
		return "", 0, "", handleError.ERR_REQUEST_CANCEL_GET_DB, err
	}
	return tid, price, pgcode, handleError.SUCCESS, nil
}

func deletePayInfo(req ReqRequestCancel) (int, error) {
	query := fmt.Sprintf("UPDATE subscribedServices SET canceled = 1 WHERE subscribedServiceId = %d", req.SubscribedServiceId)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		return handleError.ERR_REQUEST_CANCEL_DB_DELETE, err
	}
	return handleError.SUCCESS, nil
}
