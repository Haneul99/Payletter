package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqPayletterResult struct {
	Code            string `json:"code" form:"code"`
	Message         string `json:"message" form:"message"`
	UserID          string `json:"user_id" form:"user_id"`
	UserName        string `json:"user_name" form:"user_name"`
	OrderNo         string `json:"order_no" form:"order_no"`
	ServiceName     string `json:"service_name" form:"service_name"`
	ProductName     string `json:"product_name" form:"product_name"`
	CustomParameter string `json:"custom_parameter" form:"custom_parameter"`
	TID             string `json:"tid" form:"tid"`
	CID             string `json:"cid" form:"cid"`
	Amount          int    `json:"amount" form:"amount"`
	TaxFreeAmount   int    `json:"taxfree_amount" form:"taxfree_amount"`
	TaxAmount       int    `json:"tax_amount" form:"tax_amount"`
	PayInfo         string `json:"pay_info" form:"pay_info"`
	PgCode          string `json:"pgcode" form:"pacode"`
	BillKey         string `json:"billkey" form:"billkey"`
	DomesticFlag    string `json:"domestic_flag" form:"domestic_flag"`
	TransactionDate string `json:"transaction_date" form:"transaction_date"`
	InstallMonth    int    `json:"install_month" form:"install_month"`
	CardInfo        string `json:"card_info" form:"card_info"`
	PayHash         string `json:"payhash" form:"payhash"`
}

type ResReturnPayResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func PayletterResult(c echo.Context) error {
	reqPayletterResult := ReqPayletterResult{}

	if err := c.Bind(&reqPayletterResult); err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, reqPayletterResult)
}
