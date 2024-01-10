package controller

import (
	"api/resources/static"
	"api/src/model"
	"api/src/service"
	"fmt"
	"log"
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
	// 検索(グループ)
	SearchGroups(e echo.Context) error
	// グループ登録
	InsertGroup(e echo.Context) error
	// スケジュール登録種別一覧
	ListScheduleType(e echo.Context) error
	// スケジュール登録
	InsertSchedules(e echo.Context) error
	// スケジュール一覧
	Schedules(e echo.Context) error
	// スケジュール削除
	DeleteSchedule(e echo.Context) error
	// 予約表提示
	DispReserveTable(e echo.Context) error
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
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 登録
func (c *UserController) Create(e echo.Context) error {
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	res, err := c.s.Create(&req)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// ロール一覧
func (c *UserController) RoleList(e echo.Context) error {
	res, err := c.s.RoleList()
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 検索(グループ)
func (c *UserController) SearchGroups(e echo.Context) error {
	res, err := c.s.SearchGroups()
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// グループ登録
func (c *UserController) InsertGroup(e echo.Context) error {
	req := model.UserGroup{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.CreateGroup(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// スケジュール登録種別一覧
func (c *UserController) ListScheduleType(e echo.Context) error {
	res, err := c.s.ListScheduleType()
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// スケジュール登録
func (c *UserController) InsertSchedules(e echo.Context) error {
	req := model.UserScheduleRequest{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.CreateSchedule(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// スケジュール一覧
func (c *UserController) Schedules(e echo.Context) error {
	res, err := c.s.Schedules()
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// スケジュール削除
func (c *UserController) DeleteSchedule(e echo.Context) error {
	req := model.UserSchedule{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.DeleteSchedule(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// 予約表提示
func (c *UserController) DispReserveTable(e echo.Context) error {
	res, err := c.s.DispReserveTable()
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}
