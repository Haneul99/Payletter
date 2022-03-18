package util

import (
	"database/sql"
	"fmt"
)

type ottservice struct {
	OTTserviceId int64
	platform     string
	membership   string
	price        int64
}

var db *sql.DB

// 데이터베이스 open
func DBConnect() error {
	var err error

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		ServerConfig.GetStringData("DB_LoginID"),
		ServerConfig.GetStringData("DB_Password"),
		ServerConfig.GetStringData("DB_IP"),
		ServerConfig.GetStringData("DB_Port"),
		ServerConfig.GetStringData("DB_Name"))

	db, err = sql.Open("mysql", dbURL)
	if err != nil {
		return err
	}
	fmt.Println("db connection success")
	return nil
}

// OTTservices Table 정보 SELECT
func GetOTTservices() ([]ottservice, error) {
	query := fmt.Sprintf("SELECT * FROM %s", "ottservices")
	fmt.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []ottservice{}

	for rows.Next() {
		var ott ottservice
		err = rows.Scan(&ott.OTTserviceId, &ott.platform, &ott.membership, &ott.price)
		if err != nil {
			return nil, err
		}
		results = append(results, ott)
	}
	return results, nil
}

func GetDB() *sql.DB {
	if db != nil {
		return db
	}
	return nil
}

// 테이블 전체 값 SELECT
func SelectAll(table string) {

}

// 테이블 중 하나의 COLUMN SELECT
func SelectOneColumn(table, target string) {

}
