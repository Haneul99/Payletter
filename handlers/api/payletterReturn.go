package handlers

import (
	handleError "Haneul99/Payletter/handlers/error"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReqPayletterReturn struct {
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
	PgCode          string `json:"pgcode" form:"pgcode"`
	BillKey         string `json:"billkey" form:"billkey"`
	DomesticFlag    string `json:"domestic_flag" form:"domestic_flag"`
	TransactionDate string `json:"transaction_date" form:"transaction_date"`
	InstallMonth    int    `json:"install_month" form:"install_month"`
	CardInfo        string `json:"card_info" form:"card_info"`
	PayHash         string `json:"payhash" form:"payhash"`
}

type ResPayletterReturn struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func PayletterReturn(c echo.Context) error {
	reqPayletterReturn := ReqPayletterReturn{}
	//resPayletterReturn := ResPayletterReturn{}

	// Bind
	if err := c.Bind(&reqPayletterReturn); err != nil {
		return handleError.ReturnResFail(c, http.StatusInternalServerError, err, handleError.ERR_PAYLETTER_RETURN_REQUEST_BINDING)
	}

	// CheckParam
	//if isVerified, errCode, err := util.VerifyPayment(reqPayletterReturn.PayHash, reqPayletterReturn.UserID, reqPayletterReturn.TID, reqPayletterReturn.Amount); !isVerified || err != nil {
	//	return handleError.ReturnResFail(c, http.StatusInternalServerError, err, errCode)
	//}

	// Return
	return c.JSON(http.StatusOK, reqPayletterReturn)
}
