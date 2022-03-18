package handlers

import (
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var ERR_ID_OVERLAP = 40000

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ResSignUp struct {
	Success  bool   `json:"success"`
	Message  string `json: "message"`
	Username string `json: "username"`
}

func SignUp(c echo.Context) error {
	user := User{}
	resSignUp := ResSignUp{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	isAvailable, err := checkValidity(user)
	if !isAvailable || err != nil {
		return err // 실패 response return 할 것
	}

	if err := insertUserDB(user); err != nil {
		return err
	}

	resSignUp.Success = true
	resSignUp.Message = "회원가입 성공"
	resSignUp.Username = user.Username
	return c.JSON(http.StatusOK, resSignUp)
}

// 유저가 입력한 정보가 가입 가능한 정보인지 체크
func checkValidity(user User) (bool, error) {
	isAvailable, err := checkIdUnique(user.Username)
	if err != nil {
		return false, err
	}
	if !isAvailable {
		return false, nil
	}
	return true, nil
}

// 아이디 중복 체크
func checkIdUnique(username string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = \"%s\"", "user", username)
	exist := 0
	err := util.GetDB().QueryRow(query).Scan(&exist)
	if err != nil {
		return false, err
	} // 읽기 실패
	if exist != 0 {
		return false, nil
	} // 이미 존재하는 username
	return true, nil
}

// INSERT USER DB
func insertUserDB(user User) error {
	var key = []byte("thisis32bitlongpassphraseimusing")
	a, err := util.NewAesCipher(key)

	if err != nil {
		fmt.Println(err)
		return err
	}

	encryptedPassword := a.EncryptString(user.Password) // 비밀번호 암호와
	//d := a.DecryptString(e) // 복호화된 비밀번호
	// rows, err := util.GetDB().Query(query)
	query := fmt.Sprintf("INSERT INTO USER(username, password, email) VALUE(\"%s\", \"%s\", \"%s\")", user.Username, encryptedPassword, user.Email)
	_, err = util.GetDB().Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
