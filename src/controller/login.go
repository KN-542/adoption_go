package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ILoginController interface {
	// ログイン
	Login(e echo.Context) error
}

type loginController struct{}

func NewLoginController() ILoginController {
	return &loginController{}
}

// ログイン
func (c *loginController) Login(e echo.Context) error {
	return e.JSON(http.StatusOK, "OK")
}
