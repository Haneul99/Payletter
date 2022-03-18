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

type ResLoadPlatFormsList struct {
	Success  bool     `json: "success`
	Message  string   `json: "message"`
	Contents []string `json: "contents"`
}

// Product
func LoadPlatformsList(c echo.Context) error {
	resLoadPlatFormsList := ResLoadPlatFormsList{}
	results := SelectPlatformList()
	if results == nil {
		return echo.NewHTTPError(http.StatusBadRequest) // 임시
	}
	for _, value := range results {
		resLoadPlatFormsList.Contents = append(resLoadPlatFormsList.Contents, value.Platform)
	}
	resLoadPlatFormsList.Success = true
	resLoadPlatFormsList.Message = "Load Platform Lists"
	return c.JSON(http.StatusOK, resLoadPlatFormsList)
}

// Platform 정보 SELECT
func SelectPlatformList() []Product {
	query := fmt.Sprintf("SELECT DISTINCT platform FROM %s", "ottservices")
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
