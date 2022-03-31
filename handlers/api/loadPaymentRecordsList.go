package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqLoadPaymentRecordsList struct {
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}

type PaymentRecord struct {
	SubscribedServiceId int    `json:"subscribedServiceId"`
	OTTServiceId        int    `json:"OttServiceId"`
	Platform            string `json:"platform"`
	Membership          string `json:"membership"`
	SubscribedDate      string `json:"subscribedDate"`
	ExpireDate          string `json:"expireDate"`
	PaymentType         int    `json:"paymentType"`
	TID                 string `json:"tid"`
	Price               int    `json:"price"`
	PgCode              string `json:"pgCode"`
	Canceled            int    `json:"canceled"`
}

type ResLoadPaymentRecordsList struct {
	ErrCode  int             `json:"errCode"`
	Contents []PaymentRecord `json:"contents"`
}

func LoadPaymentRecordslist(c echo.Context) error {
	reqLoadPaymentRecordsList := ReqLoadPaymentRecordsList{}
	resLoadPaymentRecordsList := ResLoadPaymentRecordsList{}

	// Bind
	if err := c.Bind(&reqLoadPaymentRecordsList); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_LOAD_PAYMENT_RECORDS_LIST_REQUEST_BINDING)
	}

	// CheckParam
	if errCode, err := util.IsValidAccessToken(reqLoadPaymentRecordsList.AccessToken, reqLoadPaymentRecordsList.Username); err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	// Process
	record, errCode, err := getPaymentRecordsList(reqLoadPaymentRecordsList)
	if err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// Return
	resLoadPaymentRecordsList.Contents = record
	return c.JSON(handleError.SUCCESS, resLoadPaymentRecordsList)
}

func getPaymentRecordsList(req ReqLoadPaymentRecordsList) ([]PaymentRecord, int, error) {
	record := []PaymentRecord{}

	query := fmt.Sprintf("SELECT subscribedServiceId, OTTServiceId, platform, membership, subscribedDate, ExpireDate, paymentType, tid, subscribedServices.price, pgcode, canceled "+
		"FROM subscribedServices LEFT JOIN ottservices ON subscribedServices.OTTServiceId = ottservices.OTTServicesId "+
		"WHERE username = \"%s\"", req.Username)

	rows, err := util.GetDB().Query(query)

	if err != nil {
		return nil, handleError.ERR_LOAD_PAYMENT_RECORDS_LIST_GET_DB, err
	}
	defer rows.Close()

	for rows.Next() {
		var p PaymentRecord
		if err = rows.Scan(&p.SubscribedServiceId, &p.OTTServiceId, &p.Platform, &p.Membership, &p.SubscribedDate, &p.ExpireDate, &p.PaymentType, &p.TID, &p.Price, &p.PgCode, &p.Canceled); err != nil {
			return nil, handleError.ERR_LOAD_PAYMENT_RECORDS_LIST_SELECT_DB, err
		}
		record = append(record, p)
	}

	return record, handleError.SUCCESS, nil
}
