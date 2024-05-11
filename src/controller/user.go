package controller

import (
	"api/resources/static"
	"api/src/model/ddl"
	"api/src/model/response"
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
	// 検索(チーム)
	SearchTeams(e echo.Context) error
	// チーム登録
	InsertTeam(e echo.Context) error
	// スケジュール登録種別一覧
	ListScheduleType(e echo.Context) error
	// スケジュール登録
	InsertSchedules(e echo.Context) error
	// スケジュール更新
	UpdateSchedule(e echo.Context) error
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
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 登録
func (c *UserController) Create(e echo.Context) error {
	req := ddl.User{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	res, err := c.s.Create(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 検索(チーム)
func (c *UserController) SearchTeams(e echo.Context) error {
	res, err := c.s.SearchTeams()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// チーム登録
func (c *UserController) InsertTeam(e echo.Context) error {
	req := ddl.TeamRequest{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.CreateTeam(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// スケジュール登録種別一覧
func (c *UserController) ListScheduleType(e echo.Context) error {
	res, err := c.s.ListScheduleType()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// スケジュール登録
func (c *UserController) InsertSchedules(e echo.Context) error {
	req := ddl.UserScheduleRequest{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	_, err := c.s.CreateSchedule(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// スケジュール更新
func (c *UserController) UpdateSchedule(e echo.Context) error {
	req := ddl.UserScheduleRequest{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.UpdateSchedule(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// スケジュール一覧
func (c *UserController) Schedules(e echo.Context) error {
	res, err := c.s.Schedules()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// スケジュール削除
func (c *UserController) DeleteSchedule(e echo.Context) error {
	req := ddl.UserSchedule{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.DeleteSchedule(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// 予約表提示
func (c *UserController) DispReserveTable(e echo.Context) error {
	res, err := c.s.DispReserveTable()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}
