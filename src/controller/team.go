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

type ITeamController interface {
	// 検索
	Search(e echo.Context) error
	// 登録
	Create(e echo.Context) error
	// 更新
	Update(e echo.Context) error
	// 基本情報更新
	UpdateBasic(e echo.Context) error
	// 取得
	Get(e echo.Context) error
	// 自取得
	GetOwn(e echo.Context) error
	// 削除
	Delete(e echo.Context) error
	// 検索_同一企業
	SearchByCompany(e echo.Context) error
	// 毎ステータスイベント取得
	StatusEvents(e echo.Context) error
}

type TeamController struct {
	s     service.ITeamService
	a     service.IApplicantService
	login service.ILoginService
	role  service.IRoleService
}

func NewTeamController(
	s service.ITeamService,
	a service.IApplicantService,
	login service.ILoginService,
	role service.IRoleService,
) ITeamController {
	return &TeamController{s, a, login, role}
}

func (c *TeamController) GetLoginService() service.ILoginService {
	return c.login
}

// 検索
func (c *TeamController) Search(e echo.Context) error {
	req := request.SearchTeam{}
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
		ID: static.ROLE_MANAGEMENT_TEAM_READ,
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

// 登録
func (c *TeamController) Create(e echo.Context) error {
	req := request.CreateTeam{}
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
		ID: static.ROLE_MANAGEMENT_TEAM_CREATE,
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
func (c *TeamController) Update(e echo.Context) error {
	req := request.UpdateTeam{}
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
		ID: static.ROLE_MANAGEMENT_TEAM_EDIT,
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

// 基本情報更新
func (c *TeamController) UpdateBasic(e echo.Context) error {
	req := request.UpdateBasicTeam{}
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
		ID: static.ROLE_MANAGEMENT_SETTING_TEAM,
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

	if err := c.s.UpdateBasic(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// 削除
func (c *TeamController) Delete(e echo.Context) error {
	req := request.DeleteTeam{}
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
		ID: static.ROLE_MANAGEMENT_TEAM_DELETE,
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

// 取得
func (c *TeamController) Get(e echo.Context) error {
	req := request.GetTeam{}
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
		ID: static.ROLE_MANAGEMENT_TEAM_DETAIL_READ,
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

	res, err := c.s.Get(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 自取得
func (c *TeamController) GetOwn(e echo.Context) error {
	req := request.GetOwnTeam{}
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
		ID: static.ROLE_MANAGEMENT_SETTING_TEAM,
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

	res, err := c.s.GetOwn(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 検索_同一企業
func (c *TeamController) SearchByCompany(e echo.Context) error {
	req := request.SearchTeamByCompany{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.HashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.HashKey,
		},
		ID: static.ROLE_MANAGEMENT_TEAM_READ,
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

	res, err := c.s.SearchByCompany(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 毎ステータスイベント取得
func (c *TeamController) StatusEvents(e echo.Context) error {
	req := request.StatusEventsByTeam{}
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
		ID: static.ROLE_MANAGEMENT_SETTING_TEAM,
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

	res, err := c.s.StatusEvents(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}
