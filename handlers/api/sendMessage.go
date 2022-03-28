package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

	if errCode, err := insertMessageDB(reqSendMessage); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

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
