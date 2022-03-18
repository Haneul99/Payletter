package main

import (
	"fmt"
	"net/http"

	handlers "Haneul99/Payletter/handlers/api"
	"Haneul99/Payletter/util"

	"github.com/labstack/echo/v4"
)

type ottservice struct {
	OTTserviceId int64
	platform     string
	membership   string
	price        int64
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func main() {
	if util.ServerConfig.LoadConfig() != nil {
		panic("설정파일 읽기 실패")
	}

	//fmt.Println(util.ServerConfig.GetData())
	if util.DBConnect() != nil {
		panic("DB connect 실패")
	}

	results, err := util.GetOTTservices()
	if err != nil {
		panic("SELECT DB 실패")
	}
	fmt.Println(results)

	e := echo.New()

	apiHandlers(e)

	e.Logger.Fatal(e.Start(":8080"))
}

func apiHandlers(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello Worldddd\n") })
	e.GET("/api/loadProdouctsList", handlers.LoadPlatformsList)
	e.POST("/api/signUp", handlers.SignUp)
	e.GET("/api/loadPlatformDetail", handlers.LoadPlatformDetail)
	e.POST("/api/requestPay", handlers.ReqestPay)
}
