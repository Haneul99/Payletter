package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PlatformDetail struct {
	OTTserviceId int    `json:"OTTserviceId"`
	Platform     string `json:"platform"`
	Membership   string `json:"membership"`
	Price        int    `json:"price"`
}

type ResLoadPlatformDetail struct {
	ErrCode  int              `json:"errCode"`
	Contents []PlatformDetail `json:"contents"`
}

func LoadPlatformDetail(c echo.Context) error {
	resLoadPlatformDetail := ResLoadPlatformDetail{}

	// Bind
	platformName := c.QueryParam("platform")

	// Process
	results, errCode, err := getPlatformDetail(platformName)
	if err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	}

	// Return
	resLoadPlatformDetail.ErrCode = handleError.SUCCESS
	resLoadPlatformDetail.Contents = results
	return c.JSON(http.StatusOK, resLoadPlatformDetail)
}

// Platform 정보 SELECT
func getPlatformDetail(platformName string) ([]PlatformDetail, int, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE platform = \"%s\"", "ottservices", platformName)
	rows, err := util.GetDB().Query(query)
	if err != nil {
		return nil, handleError.ERR_JWT_GET_DB, err
	}
	defer rows.Close()

	results := []PlatformDetail{}

	for rows.Next() {
		var platformDetail PlatformDetail
		if err = rows.Scan(&platformDetail.OTTserviceId, &platformDetail.Platform, &platformDetail.Membership, &platformDetail.Price); err != nil {
			return nil, handleError.ERR_LOAD_PLATFORM_DETAIL_SELECT_DB, err
		}
		results = append(results, platformDetail)
	}
	return results, handleError.SUCCESS, nil
}
