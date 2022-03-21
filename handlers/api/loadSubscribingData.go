package handlers

import (
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
	Success  bool         `json:"success"`
	Contents []OttService `json:"contents"`
}

func LoadSubscribingData(c echo.Context) error {
	reqLoadSubscribingData := ReqLoadSubscribingData{}
	resLoadSubscribingData := ResLoadSubscribingData{}

	if err := c.Bind(&reqLoadSubscribingData); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_REQUEST_BINDING})
	}

	if isValid, err := util.IsValidAccessToken(reqLoadSubscribingData.AccessToken, reqLoadSubscribingData.Username); !isValid || err != nil {
		return c.JSON(http.StatusUnauthorized, ResFail{ErrCode: false, Message: ERR_ACCESSTOKEN})
	}

	if subscribed, err := getSubscribingData(ReqLoadPeronsalData(reqLoadSubscribingData)); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_SELECT_DB})
	} else {
		resLoadSubscribingData.Contents = subscribed
	}
	resLoadSubscribingData.Success = true

	return c.JSON(http.StatusOK, resLoadSubscribingData)
}

func getSubscribingData(user ReqLoadPeronsalData) ([]OttService, error) {
	subscribed := []OttService{}

	query := fmt.Sprintf("SELECT subscribedServiceId, OTTServiceId, subscribedDate, ExpireDate, paymentType FROM subscribedServices WHERE username = \"%s\"", user.Username)
	rows, err := util.GetDB().Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var service OttService
		if err = rows.Scan(&service.SubscribedServiceId, &service.OTTServiceId, &service.SubscribedDate, &service.ExpireDate, &service.PaymentType); err != nil {
			return nil, err
		}
		subscribed = append(subscribed, service)
	}

	return subscribed, nil
}
