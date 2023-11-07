package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ILoginController interface {
	// ログイン
	Login(e echo.Context) error
}

type LoginController struct{}

func NewLoginController() ILoginController {
	return &LoginController{}
}

// ログイン
func (c *LoginController) Login(e echo.Context) error {
	return e.JSON(http.StatusOK, "OK")
}
