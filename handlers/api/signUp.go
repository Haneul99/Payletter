package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ResSignUp struct {
	ErrCode  int    `json:"errCode"`
	Username string `json:"username"`
}

func SignUp(c echo.Context) error {
	user := ReqSignUp{}
	resSignUp := ResSignUp{}

	// Bind
	if err := c.Bind(&user); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_SIGN_UP_REQUEST_BINDING)
	}

	// CheckParam
	if isAvailable, status, errCode, err := signUpCheckParam(user); !isAvailable || err != nil {
		return handleError.ReturnResFail(c, status, err, errCode)
	}

	// Process
	if errCode, err := insertUserDB(user); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// Return
	resSignUp.ErrCode = 0
	resSignUp.Username = user.Username
	return c.JSON(http.StatusOK, resSignUp)
}

// 유저가 입력한 정보가 가입 가능한 정보인지 체크
func signUpCheckParam(user ReqSignUp) (bool, int, int, error) {
	if isAvailable, status, errCode, err := checkParamId(user.Username); !isAvailable || err != nil {
		return false, status, errCode, err
	}
	// password가 빈 값인지 확인
	if isNotNull, status, errCode, err := checkParamPassword(user.Password); !isNotNull || err != nil {
		return false, status, errCode, err
	}
	return true, http.StatusOK, 0, nil
}

// 아이디 중복 체크
func checkParamId(username string) (bool, int, int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = \"%s\"", "user", username)
	exist := 0
	if err := util.GetDB().QueryRow(query).Scan(&exist); err != nil {
		return false, http.StatusInternalServerError, handleError.ERR_SIGN_UP_GET_DB, err
	}
	if exist != 0 {
		return false, http.StatusBadRequest, handleError.ERR_SIGN_UP_DUPLICATED_ID, errors.New("ERR_SIGN_UP_DUPLICATED_ID")
	}
	return true, http.StatusOK, 0, nil
}

func checkParamPassword(password string) (bool, int, int, error) {
	if len(password) == 0 {
		return false, http.StatusBadRequest, handleError.ERR_SIGN_UP_NULL_PASSWORD, errors.New("ERR_SIGNUP_NULL_PASSWORD")
	}
	return true, http.StatusOK, 0, nil
}

// DB에 유저 정보 삽입
func insertUserDB(user ReqSignUp) (int, error) {
	query := fmt.Sprintf("INSERT INTO USER(username, password, email) VALUE(\"%s\", \"%s\", \"%s\")", user.Username, user.Password, user.Email)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return handleError.ERR_SIGN_UP_SQL_NO_RESULT, err
		}
		return handleError.ERR_SIGN_UP_GET_DB, err
	}
	return 0, nil
}
