package util

import (
	_ "github.com/go-sql-driver/mysql"
)

/*
func GetDBConnection(key string) (*sql.DB, error) {
	db, exists := sql.Open("mysql", "")
	if !exists {
		return nil, fmt.Errorf("not found db : %s", key)
	}
	return db, nil
}
*/
