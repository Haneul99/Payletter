package handlers

import (
	"Haneul99/Payletter/util"
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
	Success  bool   `json:"success"`
	Username string `json:"username"`
}

func SignUp(c echo.Context) error {
	user := ReqSignUp{}
	resSignUp := ResSignUp{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: ERR_SIGNUP_REQUEST_BINDING, Message: "ERR_SIGNUP_REQUEST_BINDING"})
	}

	if isAvailable, err := checkParam(user); !isAvailable || err != nil {
		return c.JSON(http.StatusBadRequest, ResFail{ErrCode: err, Message: ERR_DUPLICATE_ID})
	}

	if err := insertUserDB(user); err != nil {
		return c.JSON(http.StatusInternalServerError, ResFail{ErrCode: false, Message: ERR_INSERT_DB})
	}

	resSignUp.Success = true
	resSignUp.Username = user.Username
	return c.JSON(http.StatusOK, resSignUp)
}

// 유저가 입력한 정보가 가입 가능한 정보인지 체크
func checkParam(user ReqSignUp) (bool, error) {
	if isAvailable, err := checkParamId(user.Username); !isAvailable || err != nil {
		return false, fmt.Errorf(ERR_DUPLICATE_ID)
	}
	// password가 빈 값인지 확인
	if isNull, err := checkParamPassword(user.Password); isNull || err != nil {
		return false, fmt.Errorf("ERR_SIGNUP_NULL_PASSWORD") // 받아온 err 값을 return 하도록 수정 필요.
	}
	return true, nil
}

// 아이디 중복 체크
func checkParamId(username string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = \"%s\"", "user", username)
	exist := 0
	if err := util.GetDB().QueryRow(query).Scan(&exist); err != nil {
		return false, err
	} // 읽기 실패
	if exist != 0 {
		return false, nil
	} // 이미 존재하는 username
	return true, nil
}

func checkParamPassword(password string) (bool, error) {
	if len(password) == 0 {
		return false, fmt.Errorf("ERR_SIGNUP_NULL_PASSWORD")
	}
	return true, nil
}

// DB에 유저 정보 삽입
func insertUserDB(user ReqSignUp) error {
	query := fmt.Sprintf("INSERT INTO USER(username, password, email) VALUE(\"%s\", \"%s\", \"%s\")", user.Username, user.Password, user.Email)
	_, err := util.GetDB().Exec(query)
	if err != nil {
		return err
	}
	return nil
}
