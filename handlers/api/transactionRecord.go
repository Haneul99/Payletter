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

type ReqTransactionRecord struct {
	SubscribedServiceId int    `json:"subscribedServiceId"`
	Username            string `json:"username"`
	AccessToken         string `json:"accessToken"`
}

type ResTransactionRecord struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	ReceiptURL string `json:"receipt_url"`
}

func TransactionRecord(c echo.Context) error {
	reqTransactionRecord := ReqTransactionRecord{}

	// Bind
	if err := c.Bind(&reqTransactionRecord); err != nil {
		return c.JSON(http.StatusInternalServerError, handleError.ERR_TRANSACTION_RECORD_REQUEST_BINDING)
	}

	// CheckParam
	if errCode, err := util.IsValidAccessToken(reqTransactionRecord.AccessToken, reqTransactionRecord.Username); err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	tid, amount, transactionDate, errCode, err := getSubscribedInfo(reqTransactionRecord)
	if err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// Process
	respBody, errCode, err := util.RequestTransactionRecordAPI(tid, amount, transactionDate)
	if err != nil {
		handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	resTransactionRecord := ResTransactionRecord{}
	if err := json.Unmarshal(respBody, &resTransactionRecord); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_TRANSACTION_RECORD_JSON_UNMARSHAL)
	}

	// Return
	return c.JSON(http.StatusOK, resTransactionRecord)
}

// DB에서 tid로 거래 정보 불러오기 or OTTServiceId로 거래정보 불러오기
func getSubscribedInfo(req ReqTransactionRecord) (string, int, string, int, error) {
	tid := ""
	amount := 0
	transactionDate := ""
	query := fmt.Sprintf("SELECT tid, price, subscribedDate FROM %s WHERE username = \"%s\" && subscribedServiceId = %d", "subscribedServices", req.Username, req.SubscribedServiceId)
	if err := util.GetDB().QueryRow(query).Scan(&tid, &amount, &transactionDate); err != nil {
		if err == sql.ErrNoRows {
			return "", 0, "", handleError.ERR_TRANSACTION_RECORD_SQL_NO_RESULT, err
		}
		return "", 0, "", handleError.ERR_TRANSACTION_RECORD_GET_DB, err
	}

	return tid, amount, transactionDate, 0, nil
} // tid, amount, transactionDate return
