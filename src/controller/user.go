package controller

import (
	"api/src/model"
	"api/src/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	// 一覧
	List(e echo.Context) error
	// 登録
	Create(e echo.Context) error
	// ロール一覧
	RoleList(e echo.Context) error
}

type UserController struct {
	s service.IUserService
}

func NewUserController(s service.IUserService) IUserController {
	return &UserController{s}
}

// 一覧
func (c *UserController) List(e echo.Context) error {
	res, err := c.s.List()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, res)
}

// 登録
func (c *UserController) Create(e echo.Context) error {
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := c.s.Create(&req)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, res)
}

// ロール一覧
func (c *UserController) RoleList(e echo.Context) error {
	res, err := c.s.RoleList()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, res)
}
