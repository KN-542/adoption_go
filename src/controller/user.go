package controller

import (
	"api/src/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	// 一覧
	List(e echo.Context) error
}

type userController struct {
	s service.IUserService
}

func NewUserController(s service.IUserService) IUserController {
	return &userController{s}
}

// 一覧
func (c *userController) List(e echo.Context) error {
	res, err := c.s.List()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, res)
}
