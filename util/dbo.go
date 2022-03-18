package util

import (
	"database/sql"
	"fmt"
)

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

func GetDB() *sql.DB {
	if db != nil {
		return db
	}
	return nil
}
