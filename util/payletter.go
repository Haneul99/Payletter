package util

import (
	handleError "Haneul99/Payletter/handlers/error"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReqPayletterRequestData struct {
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

type ReqPayletterCancelData struct {
	PgCode   string `json:"pgcode"`
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
	TID      string `json:"tid"`
	Amount   int    `json:"amount"`
	IpAddr   string `json:"ip_addr"`
}

func RequestPayAPI(username string, platform string, membership string, OTTserviceId int, amount int) ([]byte, int, error) {
	reqPayletterRequestData := ReqPayletterRequestData{}
	reqPayletterRequestData.PgCode = "creditcard"
	reqPayletterRequestData.ClientID = ServerConfig.GetStringData("Payletter_CLIENT_ID")
	reqPayletterRequestData.UserID = username
	reqPayletterRequestData.Amount = amount
	reqPayletterRequestData.ProductName = fmt.Sprintf("%d_%s_%s", OTTserviceId, platform, membership)
	reqPayletterRequestData.ReturnURL = "http://127.0.0.1:8080/api/payletterReturn"
	reqPayletterRequestData.CallbackURL = "http://127.0.0.1:8080/api/payletterCallback"
	reqPayletterRequestData.CancelURL = "https://testpg.payletter.com/cancel"

	jsonData, err := json.Marshal(reqPayletterRequestData)
	if err != nil {
		return nil, handleError.ERR_PAYLETTER_JSON_MARSHAL, err
	}
	return requestPayletterAPI(http.MethodPost, "v1.0/payments/request", jsonData, "PAYMENT")
}

func RequestCancelAPI(username string, pgcode string, tid string, amount int) ([]byte, int, error) {
	reqPayletterCancelData := ReqPayletterCancelData{}
	reqPayletterCancelData.PgCode = pgcode
	reqPayletterCancelData.ClientID = ServerConfig.GetStringData("Payletter_CLIENT_ID")
	reqPayletterCancelData.UserID = username
	reqPayletterCancelData.TID = tid
	reqPayletterCancelData.Amount = amount
	reqPayletterCancelData.IpAddr = "127.0.0.1"

	jsonData, err := json.Marshal(reqPayletterCancelData)
	if err != nil {
		return nil, handleError.ERR_PAYLETTER_JSON_MARSHAL, err
	}
	return requestPayletterAPI(http.MethodPost, "v1.0/payments/cancel", jsonData, "PAYMENT")
}

func RequestTransactionRecordAPI(tid string, amount int, transaction_date string) ([]byte, int, error) {
	transaction_date = transaction_date[:4] + transaction_date[5:7] + transaction_date[8:10]
	uri := fmt.Sprintf("v1.0/receipt/info/%s/?client_Id=%s&amount=%d&transaction_date=%s", tid, ServerConfig.GetStringData("Payletter_CLIENT_ID"), amount, transaction_date)
	return requestPayletterAPI(http.MethodGet, uri, nil, "SEARCH")
}

// Payletter api 호출
func requestPayletterAPI(method string, uri string, jsonData []byte, authType string) ([]byte, int, error) {
	client := httpClient()
	req, err := http.NewRequest(method, ServerConfig.GetStringData("Payletter_ENDPOINT")+uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, handleError.ERR_PAYLETTER_NEW_REQUEST, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "PLKEY "+ServerConfig.GetStringData(fmt.Sprintf("Payletter_%s_API_KEY", authType)))

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

// Return URL, CallBack URL로 전달된 payhash값과 검증하는 단계
// user_id + amount + tid + 결제용 API KEY 로 sha256 hash 값 생성
func VerifyPayment(payhash, username, tid string, amount int) (bool, int, error) {
	data := username + strconv.Itoa(amount) + tid + ServerConfig.GetStringData("Payletter_PAYMENT_API_KEY")
	fmt.Println(data)

	hash := sha256.New()
	hash.Write([]byte(data))
	hashData := hex.EncodeToString(hash.Sum(nil))

	if !strings.EqualFold(payhash, hashData) {
		return false, handleError.ERR_PAYLETTER_PAYHASH_INVALID, errors.New("ERR_PAYLETTER_PAYHASH_INVALID")
	}
	return true, 0, nil
}
