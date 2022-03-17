package handlers

import (
	"fmt"
	"net/http"

	"Haneul99/Payletter/util"

	"github.com/labstack/echo/v4"
)

type Product struct {
	Platform string `json: "platform"`
}

// Product
func LoadPlatformsList(c echo.Context) error {
	results := SelectPlatformList()
	fmt.Println(results)
	str := ""
	for _, value := range results {
		str += value.Platform + " "
	}
	return c.String(http.StatusOK, str)
}

// Platform 정보 SELECT
func SelectPlatformList() []Product {
	query := fmt.Sprintf("SELECT DISTINCT platform FROM %s", "ottservices")
	fmt.Println(query)
	rows, err := util.GetDB().Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	results := []Product{}

	for rows.Next() {
		var platform Product
		err = rows.Scan(&platform.Platform)
		if err != nil {
			return nil
		}
		results = append(results, platform)
	}
	return results
}
