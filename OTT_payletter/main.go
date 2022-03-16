package main

import (
	"fmt"
	"log"
	"net/http"

	"Haneul99/OTT_payletter/util"
	"database/sql"
	//	"github.com/labstack/echo"
)

var db *sql.DB

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func dbConnect() {
	var err error

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		util.ServerConfig.GetStringData("DB_OTTS_LoginID"),
		util.ServerConfig.GetStringData("DB_OTTS_Password"),
		util.ServerConfig.GetStringData("DB_OTTS_IP"),
		util.ServerConfig.GetStringData("DB_OTTS_Port"),
		util.ServerConfig.GetStringData("DB_OTTS_Name"))

	db, err = sql.Open("mysql", dbURL)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("db connection success")
	}
}

func main() {
	if !util.ServerConfig.LoadConfig() {
		panic("설정파일 읽기 실패")
	}

	//fmt.Println(util.ServerConfig.GetData())
	dbConnect()
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
