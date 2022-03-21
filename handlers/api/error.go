package handlers

type ResFail struct {
	ErrCode int    `json:"errCode"`
	Message string `json:"message"`
}

const (
	SUCCESS = 0

	//SignUp api 처리 시 발생 에러
	ERR_SIGNUP_REQUEST_BINDING = 10000 // SignUp 데이터 바인딩 에러
	ERR_SIGNUP_DUPLICATED_ID   = 10001 // 중복된 Username
	ERR_SIGNUP_NULL_PASSWORD   = 10002 // 빈 Password
	ERR_SIGNUP_GET_DB          = 10003 // DB Conn 실패

	// Login api 처리 시 발생 에러

	// Logout api 처리 시 발생 에러

	// LoadPlatformsList api 처리 시 발생 에러

	// LoadPlatformDetail api 처리 시 발생 에러

	// LoadPersonalData api 처리 시 발생 에러

	// LoadSubscribingData api 처리 시 발생 에러

	// RequestPay api 처리 시 발생 에러

	// RequestCancel api 처리 시 발생 에러

	ERR_ACCESSTOKEN     = "AccessToken Error"
	ERR_REQUEST_BINDING = "Request Binding Error"

	ERR_SELECT_DB = 20000 //DB 데이터 조회 에러
	ERR_INSERT_DB = 20001 //DB 데이터 삽입 에러
	ERR_DELETE_DB = 20002 //DB 데이터 삭제 에러
	ERR_UPDATE_DB = 20003 //DB 데이터 갱신 에러

	ERR_INCORRECT_PASSWORD = "Incorrect Password Error"
	ERR_DUPLICATE_ID       = "Duplicated USER ID Error"
)
