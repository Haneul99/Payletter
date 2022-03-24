package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type CashReceipt struct {
	Code      string `json:"code" form:"code"`
	Message   string `json:"message" form:"message"`
	CID       string `json:"cid" form:"cid"`
	DealNo    string `json:"deal_no" form:"deal_no"`
	IssueType string `json:"issue_type" form:"issue_type"`
	PayerSID  string `json:"payer_sid" form:"payer_sid"`
	Type      string `json:"type" form:"type"`
}

type ReqPayletterCallback struct {
	Code            string      `json:"code" form:"code"`
	Message         string      `json:"message" form:"message"`
	UserID          string      `json:"user_id" form:"user_id"`
	UserName        string      `json:"user_name" form:"user_name"`
	OrderNo         string      `json:"order_no" form:"order_no"`
	ServiceName     string      `json:"service_name" form:"service_name"`
	ProductName     string      `json:"product_name" form:"product_name"`
	CustomParameter string      `json:"custom_parameter" form:"custom_parameter"`
	TID             string      `json:"tid" form:"tid"`
	CID             string      `json:"cid" form:"cid"`
	Amount          int         `json:"amount" form:"amount"`
	TaxFreeAmount   int         `json:"taxfree_amount" form:"taxfree_amount"`
	TaxAmount       int         `json:"tax_amount" form:"tax_amount"`
	PayInfo         string      `json:"pay_info" form:"pay_info"`
	PgCode          string      `json:"pgcode" form:"pacode"`
	BillKey         string      `json:"billkey" form:"billkey"`
	DomesticFlag    string      `json:"domestic_flag" form:"domestic_flag"`
	TransactionDate string      `json:"transaction_date" form:"transaction_date"`
	InstallMonth    int         `json:"install_month" form:"install_month"`
	CardInfo        string      `json:"card_info" form:"card_info"`
	PayHash         string      `json:"payhash" form:"payhash"`
	CashReCeipt     CashReceipt `json:"cash_receipt" form:"cash_receipt"`
}

type ResPayletterCallback struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func PayletterCallback(c echo.Context) error {
	reqPayletterCallback := ReqPayletterCallback{}
	resPayletterCallback := ResPayletterCallback{}

	// Bind
	if err := c.Bind(&reqPayletterCallback); err != nil {
		resPayletterCallback.Code = http.StatusInternalServerError
		resPayletterCallback.Message = err.Error()
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_PAYLETTER_CALLBACK_REQUEST_BINDING)
	}

	// CheckParam
	if isVerified, errCode, err := util.VerifyPayment(reqPayletterCallback.PayHash, reqPayletterCallback.UserID, reqPayletterCallback.TID, reqPayletterCallback.Amount); !isVerified || err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// Process
	if errCode, err := insertPayInfo(reqPayletterCallback); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// Return
	return c.JSON(http.StatusOK, resPayletterCallback)
}

// 결제 성공 후 DB에 결제 정보 삽입
func insertPayInfo(req ReqPayletterCallback) (int, error) {
	OTTServiceId, err := strconv.Atoi(strings.Split(req.ProductName, "_")[0])
	if err != nil {
		return handleError.ERR_PAYLETTER_CALLBACK_CONVERT_STR_TO_INT, err
	}

	subscribedDate, expireDate, errCode, err := getPayInfoDate(req.TransactionDate)
	if err != nil {
		return errCode, err
	}

	paymentType := getPayInfoBillkey(req.BillKey)

	query := fmt.Sprintf("INSERT INTO subscribedServices(username, OTTServiceId, subscribedDate, expireDate, paymentType, tid, price, pgcode) VALUES(\"%s\", %d, \"%s\", \"%s\", %d, \"%s\", %d, \"%s\")", req.UserID, OTTServiceId, subscribedDate, expireDate, paymentType, req.TID, req.Amount, req.PgCode)
	_, err = util.GetDB().Exec(query)
	if err != nil {
		return handleError.ERR_PAYLETTER_CALLBACK_GET_DB, err
	}
	return 0, nil
}

func getPayInfoDate(transactionDate string) (string, string, int, error) {
	subscribedDate, err := time.Parse("2006-01-02", transactionDate[:10])
	if err != nil {
		return "", "", handleError.ERR_PAYLETTER_CALLBACK_CONVERT_STR_TO_DATE, err
	}
	expireDate := subscribedDate.AddDate(0, 1, 0).String()[:10] // 한 달 후 구독 만료

	return subscribedDate.String()[:10], expireDate, 0, nil
}

func getPayInfoBillkey(billKey string) int {
	if len(billKey) == 0 {
		return 0
	}
	return 1
} // 0은 일반결제, 1은 자동결제
