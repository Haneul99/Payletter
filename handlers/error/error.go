package handlers

import (
	"github.com/labstack/echo/v4"
)

type ResFail struct {
	ErrCode int    `json:"errCode"`
	Message string `json:"message"`
}

func ReturnResFail(c echo.Context, statusCode int, err error, errCode int) error {
	resFail := ResFail{}
	resFail.ErrCode = errCode
	if err != nil {
		resFail.Message = err.Error()
	}

	return c.JSON(statusCode, resFail)
}

const (
	SUCCESS = 0

	//SignUp api 처리 시 발생 에러
	ERR_SIGN_UP_REQUEST_BINDING = 10000 // SignUp 데이터 바인딩 에러
	ERR_SIGN_UP_DUPLICATED_ID   = 10001 // 중복된 Username
	ERR_SIGN_UP_NULL_PASSWORD   = 10002 // 빈 Password
	ERR_SIGN_UP_GET_DB          = 10003 // DB Conn 실패

	// Login api 처리 시 발생 에러
	ERR_LOGIN_REQUEST_BINDING    = 10100 // Login 데이터 바인딩 에러
	ERR_LOGIN_INCORRECT_PASSWORD = 10101 // 틀린 비밀번호
	ERR_LOGIN_GET_DB             = 10102 // DB Conn 실패
	ERR_LOGIN_CREATE_ACCESSTOKEN = 10103 // AccessToken 생성 실패
	ERR_LOGIN_UPDATE_DB          = 10104 // DB Update 실패

	// Logout api 처리 시 발생 에러
	ERR_LOGOUT_REQEUST_BINDING     = 10200 // Logout 데이터 바인딩 에러
	ERR_LOGOUT_INVALID_ACCESSTOKEN = 10201 // DB 정보와 일치하지 않는 AccessToken
	ERR_LOGOUT_GET_DB              = 10202 // DB Conn 실패
	ERR_LOGOUT_DELETE_DB           = 10203 // DB Delete 실패

	// LoadPlatformsList api 처리 시 발생 에러
	ERR_LOAD_PLATFORMS_LIST_GET_DB    = 10300 // DB Conn 실패
	ERR_LOAD_PLATFORMS_LIST_SELECT_DB = 10301 // DB Select 실패

	// LoadPlatformDetail api 처리 시 발생 에러
	ERR_LOAD_PLATFORM_DETAIL_REQUEST_BINDING = 10400 // LoadPlatformDetail 데이터 바인딩 에러
	ERR_LOAD_PLATFORM_DETAIL_GET_DB          = 10401 // DB Conn 실패
	ERR_LOAD_PLATFORM_DETAIL_SELECT_DB       = 10402 // DB Select 실패

	// LoadPersonalData api 처리 시 발생 에러
	ERR_LOAD_PERSONAL_DATA_REQUEST_BINDING     = 10500 // LoadPersonalData 데이터 바인딩 에러
	ERR_LOAD_PERSONAL_DATA_INVALID_ACCESSTOKEN = 10501 // DB 정보와 일치하지 않는 AccessToken
	ERR_LOAD_PERSONAL_DATA_GET_DB              = 10502 // DB Conn 실패
	ERR_LOAD_PERSONAL_DATA_SELECT_DB           = 10503 // DB Select 실패

	// LoadSubscribingData api 처리 시 발생 에러
	ERR_LOAD_SUBSCRIBING_DATA_REQUEST_BINDING     = 10600 //LoadSubscribingData 데이터 바인딩 에러
	ERR_LOAD_SUBSCRIBING_DATA_INVALID_ACCESSTOKEN = 10601 // DB 정보와 일치하지 않는 AccessToken
	ERR_LOAD_SUBSCRIBING_DATA_GET_DB              = 10602 // DB Conn 실패
	ERR_LOAD_SUBSCRIBING_DATA_SELECT_DB           = 10603 // DB Select 실패

	// RequestPay api 처리 시 발생 에러
	ERR_REQUEST_PAY_REQUEST_BINDING     = 10700 // RequestPay 데이터 바인딩 에러
	ERR_REQUEST_PAY_INVALID_ACCESSTOKEN = 10701 // DB 정보와 일치하지 않는 AccessToken
	ERR_REQUEST_PAY_JSON_UNMARSHAL      = 10702 // JSON Unmarshal Error
	ERR_REQUEST_PAY_ALREADY_PAID        = 10703 // 이미 결제된 서비스
	ERR_REQUEST_PAY_GET_DB              = 10704 // DB Conn 실패

	// RequestCancel api 처리 시 발생 에러
	ERR_REQUEST_CANCEL_REQUEST_BINDING     = 10800 // RequestCancel 데이터 바인딩 에러
	ERR_REQUEST_CANCEL_INVALID_ACCESSTOKEN = 10801 // DB 정보와 일치하지 않는 AccessToken
	ERR_REQUEST_CANCEL_GET_DB              = 10802
	ERR_REQUEST_CANCEL_JSON_UNMARSHAL      = 10803
	ERR_REQUEST_CANCEL_DB_DELETE           = 10804

	// payletterReturn api 처리 시 발생 에러
	ERR_PAYLETTER_RETURN_REQUEST_BINDING = 10900

	//payletterCallback api 처리 시 발생 에러
	ERR_PAYLETTER_CALLBACK_REQUEST_BINDING     = 11000
	ERR_PAYLETTER_CALLBACK_GET_DB              = 11001
	ERR_PAYLETTER_CALLBACK_CONVERT_STR_TO_INT  = 11002
	ERR_PAYLETTER_CALLBACK_CONVERT_STR_TO_DATE = 11003

	//transactionRecord api 처리 시 발생 에러
	ERR_TRANSACTION_RECORD_REQUEST_BINDING = 11101
	ERR_TRANSACTION_RECORD_GET_DB          = 11102
	ERR_TRANSACTION_RECORD_JSON_UNMARSHAL  = 11103

	// jwt.go Error
	ERR_JWT_CREATE_ACCESSTOKEN    = 20000
	ERR_JWT_NULL_ACCESSTOKEN      = 20001
	ERR_JWT_INVALID_ACCESSTOKEN   = 20002
	ERR_JWT_GET_DB                = 20003
	ERR_JWT_INCORRECT_ACCESSTOKEN = 20004

	// payletter.go Error
	ERR_PAYLETTER_JSON_MARSHAL    = 20100
	ERR_PAYLETTER_NEW_REQUEST     = 20101
	ERR_PAYLETTER_CLIENT_DO       = 20102
	ERR_PAYLETTER_IOUTIL_READALL  = 20103
	ERR_PAYLETTER_PAYHASH_INVALID = 20104
)
