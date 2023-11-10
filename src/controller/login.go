package controller

import (
	"api/src/model"
	"api/src/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ILoginController interface {
	// ログイン
	Login(e echo.Context) error
}

type LoginController struct {
	s service.ILoginService
}

func NewLoginController(s service.ILoginService) ILoginController {
	return &LoginController{s}
}

// ログイン
func (c *LoginController) Login(e echo.Context) error {
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	// ログイン TODO

	// JWT＆Cookie
	cookie, err := c.s.JWT(&req.Email)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}
