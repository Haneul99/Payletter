package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"fmt"
	"net/http"

	"Haneul99/Payletter/util"

	"github.com/labstack/echo/v4"
)

type Product struct {
	Platform string `json:"platform"`
}

type ResLoadPlatFormsList struct {
	ErrCode  int      `json:"errCode"`
	Contents []string `json:"contents"`
}

func LoadPlatformsList(c echo.Context) error {
	resLoadPlatFormsList := ResLoadPlatFormsList{}

	// Process
	results, errCode, err := getPlatformList()
	if err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}
	for _, value := range results {
		resLoadPlatFormsList.Contents = append(resLoadPlatFormsList.Contents, value.Platform)
	}

	// Return
	resLoadPlatFormsList.ErrCode = 0
	return c.JSON(http.StatusOK, resLoadPlatFormsList)
}

// Platform 정보 SELECT
func getPlatformList() ([]Product, int, error) {
	query := fmt.Sprintf("SELECT DISTINCT platform FROM %s", "ottservices")
	rows, err := util.GetDB().Query(query)
	if err != nil {
		return nil, handleError.ERR_LOAD_PLATFORMS_LIST_GET_DB, err
	}
	defer rows.Close()

	results := []Product{}

	for rows.Next() {
		var platform Product
		if err = rows.Scan(&platform.Platform); err != nil {
			return nil, handleError.ERR_LOAD_PLATFORMS_LIST_SELECT_DB, err
		}
		results = append(results, platform)
	}
	return results, 0, nil
}
