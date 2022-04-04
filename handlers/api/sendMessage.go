package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thedevsaddam/retry"
)

type ReqSendMessage struct {
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	Message     string `json:"message"`
	AccessToken string `json:"accessToken"`
}

type ResSendMessage struct {
	ErrCode int `json:"errCode"`
}

func SendMessage(c echo.Context) error {
	reqSendMessage := ReqSendMessage{}
	resSendMessage := ResSendMessage{}

	if err := c.Bind(&reqSendMessage); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_SEND_MESSAGE_REQUEST_BINDING)
	}

	if errCode, err := util.IsValidAccessToken(reqSendMessage.AccessToken, reqSendMessage.Sender); err != nil {
		return handleError.ReturnResFail(c, http.StatusUnauthorized, err, errCode)
	}

	go waitTenSecond(reqSendMessage) // sleep을 지우지 않고 10초 작동 하는데, 유저한테 결과를 바로 알려주도록 하기

	return c.JSON(http.StatusOK, resSendMessage)
}

func insertMessageDB(req ReqSendMessage) (int, error) {
	datetime := time.Now().Format("2006-01-02 15:04:05")
	query := fmt.Sprintf("INSERT INTO MESSAGES(sender, receiver, message, date) VALUES (\"%s\", \"%s\", \"%s\", \"%s\")", req.Sender, req.Receiver, req.Message, datetime) // sender, receiver, message, datetime
	if _, err := util.GetDB().Exec(query); err != nil {
		return handleError.ERR_SEND_MESSAGE_GET_DB, err
	}

	return handleError.SUCCESS, nil
}

func waitTenSecond(req ReqSendMessage) (int, error) {
	time.Sleep(10 * time.Second)
	errCode, err := retry.Do(10, 1*time.Second, insertMessageDB, req)
	if err != nil {
		return errCode[0].(int), err
	}
	return handleError.SUCCESS, nil
}

// 1. Thread 조사
// 2. 이거보다 더 가벼운 방법 찾기(Thread는 그대로 사용)

// https://fransoaardi.github.io/posts/goroutine_lifecycle/
// main이 아닌 func 내에서 호출된 goroutine 은 해당 func 가 종료되었더라도, 종료되지않고 독립된 lifecycle 을 가진다.
// main 아닌 함수에서 goroutine을 해도 main은 계속해서 유지되고 있기 때문에 return 후에도 종료되지 않고 수행이 이어짐.
