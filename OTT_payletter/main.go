package main

import (
	"fmt"
	"log"
	"net/http"

	"Haneul99/OTT_payletter/util"
	//	"github.com/labstack/echo"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func main() {
	if !util.ServerConfig.LoadConfig() {
		panic("설정파일 읽기 실패")
	}

	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func registLatesHandler(e *echo.Echo, ServiceRoot string) {
//
// }
