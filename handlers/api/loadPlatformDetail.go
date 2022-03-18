package handlers

import (
	"Haneul99/Payletter/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PlatformDetail struct {
	OTTserviceId int    `json: "OTTserviceId"`
	Platform     string `json: "platform`
	Membership   string `json: "membership"`
	Price        int    `json: "price"`
}

type ResLoadPlatformDetail struct {
	Success  bool             `json: "success"`
	Contents []PlatformDetail `json: "contents"`
}

func LoadPlatformDetail(c echo.Context) error {
	platformName := c.QueryParam("platform")

	resLoadPlatformDetail := ResLoadPlatformDetail{}
	results := SelectPlatformDetail(platformName)

	resLoadPlatformDetail.Success = true
	resLoadPlatformDetail.Contents = results

	return c.JSON(http.StatusOK, resLoadPlatformDetail)
}

// Platform 정보 SELECT
func SelectPlatformDetail(platformName string) []PlatformDetail {
	query := fmt.Sprintf("SELECT * FROM %s WHERE platform = \"%s\"", "ottservices", platformName)
	rows, err := util.GetDB().Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	results := []PlatformDetail{}

	for rows.Next() {
		var platformDetail PlatformDetail
		err = rows.Scan(&platformDetail.OTTserviceId, &platformDetail.Platform, &platformDetail.Membership, &platformDetail.Price)
		if err != nil {
			return nil
		}
		results = append(results, platformDetail)
	}
	return results
}
