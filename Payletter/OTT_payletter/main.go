package main

import (
	"fmt"
	"log"
	"net/http"

	"Haneul99/Payletter/OTT_payletter/util"
	"database/sql"
	//	"github.com/labstack/echo"
)

type ottservice struct {
	OTTserviceId int64
	platform     string
	membership   string
	price        int64
}

var db *sql.DB

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func dbConnect() {
	var err error

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		util.ServerConfig.GetStringData("DB_LoginID"),
		util.ServerConfig.GetStringData("DB_Password"),
		util.ServerConfig.GetStringData("DB_IP"),
		util.ServerConfig.GetStringData("DB_Port"),
		util.ServerConfig.GetStringData("DB_Name"))

	db, err = sql.Open("mysql", dbURL)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("db connection success")
	}
}

// OTTservices Table 정보 SELECT
func getOTTservices() []ottservice {
	query := fmt.Sprintf("SELECT * FROM %s", "ottservices")
	fmt.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err, "db query failed")
	}
	defer rows.Close()

	results := []ottservice{}

	for rows.Next() {
		var ott ottservice
		err = rows.Scan(&ott.OTTserviceId, &ott.platform, &ott.membership, &ott.price)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, ott)
	}
	return results
}

func main() {
	if !util.ServerConfig.LoadConfig() {
		panic("설정파일 읽기 실패")
	}

	//fmt.Println(util.ServerConfig.GetData())
	dbConnect()

	results := getOTTservices()
	fmt.Println(results)

	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
