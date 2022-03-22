package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
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
	if err := c.Bind(&user); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_SIGN_UP_REQUEST_BINDING)
	}

	if isAvailable, errCode, err := checkParam(user); !isAvailable || err != nil {
		return handleError.ReturnResFail(c, http.StatusBadRequest, err, errCode)
	}

	if errCode, err := insertUserDB(user); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	resSignUp.ErrCode = 0
	resSignUp.Username = user.Username
	return c.JSON(http.StatusOK, resSignUp)
}

// 유저가 입력한 정보가 가입 가능한 정보인지 체크
func checkParam(user ReqSignUp) (bool, int, error) {
	if isAvailable, errCode, err := checkParamId(user.Username); !isAvailable || err != nil {
		return false, errCode, err
	}
	// password가 빈 값인지 확인
	if isNull, errCode, err := checkParamPassword(user.Password); isNull || err != nil {
		return false, errCode, err // 받아온 err 값을 return 하도록 수정 필요.
	}
	return true, 0, nil
}

// 아이디 중복 체크
func checkParamId(username string) (bool, int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = \"%s\"", "user", username)
	exist := 0
	if err := util.GetDB().QueryRow(query).Scan(&exist); err != nil {
		return false, handleError.ERR_SIGN_UP_GET_DB, err
	} // 읽기 실패
	if exist != 0 {
		return false, handleError.ERR_SIGN_UP_DUPLICATED_ID, nil
	} // 이미 존재하는 username
	return true, 0, nil
}

func checkParamPassword(password string) (bool, int, error) {
	if len(password) == 0 {
		return false, handleError.ERR_SIGN_UP_NULL_PASSWORD, errors.New("ERR_SIGNUP_NULL_PASSWORD")
	}
	return true, 0, nil
}

// DB에 유저 정보 삽입
func insertUserDB(user ReqSignUp) (int, error) {
	query := fmt.Sprintf("INSERT INTO USER(username, password, email) VALUE(\"%s\", \"%s\", \"%s\")", user.Username, user.Password, user.Email)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		return handleError.ERR_SIGN_UP_GET_DB, err
	}
	return 0, nil
}
