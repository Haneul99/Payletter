package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func SignUp(c echo.Context) error {
	userInfo := User{}
	if err := c.Bind(&userInfo); err != nil {
		fmt.Println(err)
	} // unsupported media type

	return c.JSON(http.StatusOK, userInfo)
}
