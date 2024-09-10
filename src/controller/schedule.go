package controller

import (
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/service"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IScheduleController interface {
	// 登録種別一覧
	SearchScheduleType(e echo.Context) error
	// 登録
	Insert(e echo.Context) error
	// 更新
	Update(e echo.Context) error
	// 検索
	Search(e echo.Context) error
	// 削除
	Delete(e echo.Context) error
}

type ScheduleController struct {
	s     service.IScheduleService
	a     service.IApplicantService
	login service.ILoginService
	role  service.IRoleService
}

func NewScheduleController(
	s service.IScheduleService,
	a service.IApplicantService,
	login service.ILoginService,
	role service.IRoleService,
) IScheduleController {
	return &ScheduleController{s, a, login, role}
}

func (c *ScheduleController) GetLoginService() service.ILoginService {
	return c.login
}

// 登録種別一覧
func (c *ScheduleController) SearchScheduleType(e echo.Context) error {
	res, err := c.s.SearchScheduleType()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 登録
func (c *ScheduleController) Insert(e echo.Context) error {
	req := request.CreateSchedule{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_SCHEDULE_CREATE,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusForbidden,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	if err := c.s.Create(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// 更新
func (c *ScheduleController) Update(e echo.Context) error {
	req := request.UpdateSchedule{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_SCHEDULE_EDIT,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusForbidden,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	if err := c.s.Update(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// 検索
func (c *ScheduleController) Search(e echo.Context) error {
	req := request.SearchSchedule{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_SCHEDULE_READ,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusNoContent,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	res, err := c.s.Search(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 削除
func (c *ScheduleController) Delete(e echo.Context) error {
	req := request.DeleteSchedule{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_SCHEDULE_DELETE,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusForbidden,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	if err := c.s.Delete(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}
