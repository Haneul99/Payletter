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
	Platform            string `json:"platform"`
	Membership          string `json:"membership"`
	Price               int    `json:"price"`
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
	if errCode, err := util.IsValidAccessToken(reqLoadSubscribingData.AccessToken, reqLoadSubscribingData.Username); err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	// Process
	if subscribed, errCode, err := getSubscribingData(ReqLoadPeronsalData(reqLoadSubscribingData)); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	} else {
		resLoadSubscribingData.Contents = subscribed
	}

	// Return
	resLoadSubscribingData.ErrCode = handleError.SUCCESS
	return c.JSON(http.StatusOK, resLoadSubscribingData)
}

func getSubscribingData(user ReqLoadPeronsalData) ([]OttService, int, error) {
	subscribed := []OttService{}

	query := fmt.Sprintf("SELECT subscribedServiceId, OTTServiceId, subscribedDate, ExpireDate, paymentType, platform, membership, subscribedServices.price "+
		"FROM subscribedServices LEFT JOIN ottservices ON subscribedServices.OTTServiceId = ottservices.OTTServicesId "+
		"WHERE username = \"%s\" && canceled = 0", user.Username)
	rows, err := util.GetDB().Query(query)

	if err != nil {
		return nil, handleError.ERR_LOAD_SUBSCRIBING_DATA_GET_DB, err
	}
	defer rows.Close()

	for rows.Next() {
		var ott OttService
		if err = rows.Scan(&ott.SubscribedServiceId, &ott.OTTServiceId, &ott.SubscribedDate, &ott.ExpireDate, &ott.PaymentType, &ott.Platform, &ott.Membership, &ott.Price); err != nil {
			return nil, handleError.ERR_LOAD_SUBSCRIBING_DATA_SELECT_DB, err
		}
		subscribed = append(subscribed, ott)
	}

	return subscribed, handleError.SUCCESS, nil
}
