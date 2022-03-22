package util

import (
	handleError "Haneul99/Payletter/handlers/error"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ReqPayletterData struct {
	PgCode          string `json:"pgcode"`
	ClientID        string `json:"client_id"`
	ServiceName     string `json:"service_name"`
	UserID          string `json:"user_id"`
	Amount          int    `json:"amount"`
	TaxAmount       int    `json:"tax_amount"`
	ProductName     string `json:"product_name"`
	ReturnURL       string `json:"return_url"`
	CallbackURL     string `json:"callback_url"`
	CancelURL       string `json:"cancel_url"`
	Email           string `json:"email_addr"`
	EmailFlag       string `json:"email_flag"`
	AutoPayFlag     string `json:"autopay_flag"`
	CustomParameter string `json:"custom_parameter"`
	OrderNo         string `json:"order_no"`
}

type ResPayletterData struct {
	Token     int    `json:"token"`
	OnlineURL string `json:"online_url"`
	MobileURL string `json:"mobile_url"`
}

type ReqRequestPayletterAPI struct {
	OTTserviceId int    `json:"OTTserviceId"`
	Platform     string `json:"platform"`
	Membership   string `json:"membership"`
	Price        int    `json:"price"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
}

// Payletter 결제요청 api 호출
func RequestPayletterAPI(method string, uri string, username string, price int, platform string, membership string) ([]byte, int, error) {
	client := httpClient()

	reqPayletterData := ReqPayletterData{}
	reqPayletterData.PgCode = "creditcard"
	reqPayletterData.ClientID = "pay_test"
	reqPayletterData.UserID = username
	reqPayletterData.Amount = price
	reqPayletterData.ProductName = fmt.Sprintf("%s_%s", platform, membership)
	reqPayletterData.ReturnURL = "http://127.0.0.1:8080/api/payletterResult"
	reqPayletterData.CallbackURL = "https://testpg.payletter.com/callback"
	reqPayletterData.CancelURL = "https://testpg.payletter.com/cancel"

	jsonData, err := json.Marshal(reqPayletterData)
	if err != nil {
		return nil, handleError.ERR_PAYLETTER_JSON_MARSHAL, err
	}
	req, err := http.NewRequest(method, ServerConfig.GetStringData("Payletter_ENDPOINT")+uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, handleError.ERR_PAYLETTER_NEW_REQUEST, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", ServerConfig.GetStringData("Payletter_PAYMENT_API_KEY"))

	response, err := client.Do(req)
	if err != nil {
		return nil, handleError.ERR_PAYLETTER_CLIENT_DO, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, handleError.ERR_PAYLETTER_IOUTIL_READALL, err
	}

	return body, 0, nil
}

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}
