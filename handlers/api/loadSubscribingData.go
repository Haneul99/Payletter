package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"fmt"
	"net/http"

	"Haneul99/Payletter/util"

	"github.com/labstack/echo/v4"
)

type ReqLoadSubscribingData struct {
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}

type OttService struct {
	SubscribedServiceId int    `json:"subscribedServiceId"`
	OTTServiceId        int    `json:"OTTServiceId"`
	SubscribedDate      string `json:"subscribedDate"`
	ExpireDate          string `json:"expireDate"`
	PaymentType         int    `json:"paymentType"`
}

type ResLoadSubscribingData struct {
	ErrCode  int          `json:"errCode"`
	Contents []OttService `json:"contents"`
}

func LoadSubscribingData(c echo.Context) error {
	reqLoadSubscribingData := ReqLoadSubscribingData{}
	resLoadSubscribingData := ResLoadSubscribingData{}

	// Bind
	if err := c.Bind(&reqLoadSubscribingData); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_LOAD_SUBSCRIBING_DATA_REQUEST_BINDING)
	}

	// CheckParam
	if isValid, errCode, err := util.IsValidAccessToken(reqLoadSubscribingData.AccessToken, reqLoadSubscribingData.Username); !isValid || err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	// Process
	if subscribed, errCode, err := getSubscribingData(ReqLoadPeronsalData(reqLoadSubscribingData)); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	} else {
		resLoadSubscribingData.Contents = subscribed
	}

	// Return
	resLoadSubscribingData.ErrCode = 0
	return c.JSON(http.StatusOK, resLoadSubscribingData)
}

func getSubscribingData(user ReqLoadPeronsalData) ([]OttService, int, error) {
	subscribed := []OttService{}

	query := fmt.Sprintf("SELECT subscribedServiceId, OTTServiceId, subscribedDate, ExpireDate, paymentType FROM subscribedServices WHERE username = \"%s\"", user.Username)
	rows, err := util.GetDB().Query(query)

	if err != nil {
		return nil, handleError.ERR_LOAD_SUBSCRIBING_DATA_GET_DB, err
	}
	defer rows.Close()

	for rows.Next() {
		var service OttService
		if err = rows.Scan(&service.SubscribedServiceId, &service.OTTServiceId, &service.SubscribedDate, &service.ExpireDate, &service.PaymentType); err != nil {
			return nil, handleError.ERR_LOAD_SUBSCRIBING_DATA_SELECT_DB, err
		}
		subscribed = append(subscribed, service)
	}

	return subscribed, 0, nil
}
