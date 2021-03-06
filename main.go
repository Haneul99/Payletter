package main

import (
	handlers "Haneul99/Payletter/handlers/api"
	"Haneul99/Payletter/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	if util.ServerConfig.LoadConfig() != nil {
		panic("설정파일 읽기 실패")
	}

	if util.DBConnect() != nil {
		panic("DB connect 실패")
	}

	e := echo.New()
	apiHandlers(e)
	e.Logger.Fatal(e.Start(":8080"))
}

func test(c echo.Context) error {
	return c.JSON(http.StatusOK, "connection success")
}

func apiHandlers(e *echo.Echo) {
	e.GET("/", test)
	e.GET("/api/loadPlatformsList", handlers.LoadPlatformsList)
	e.GET("/api/loadPlatformDetail", handlers.LoadPlatformDetail)
	e.POST("/api/signUp", handlers.SignUp)
	e.POST("/api/login", handlers.Login)
	e.POST("/api/logout", handlers.Logout)
	e.POST("/api/loadPersonalData", handlers.LoadPersonalData)
	e.POST("/api/loadSubscribingData", handlers.LoadSubscribingData)
	e.POST("/api/requestPay", handlers.RequestPay)
	e.POST("/api/requestCancel", handlers.RequestCancel)
	e.POST("/api/requestTransactionRecord", handlers.RequestTransactionRecord)
	e.POST("/api/payletterReturn", handlers.PayletterReturn)
	e.POST("/api/payletterCallback", handlers.PayletterCallback)
	e.POST("/api/sendMessage", handlers.SendMessage)
	e.POST("/api/loadPaymentRecordsList", handlers.LoadPaymentRecordslist)
}
