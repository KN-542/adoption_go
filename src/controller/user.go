package controller

import (
	"api/src/model/ddl"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/service"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	// 登録
	Create(e echo.Context) error
	// 検索
	Search(e echo.Context) error
	// 検索_同一企業
	SearchByCompany(e echo.Context) error
	// 応募者ステータス変更
	UpdateStatus(e echo.Context) error
	// ステータスイベントマスタ一覧
	ListStatusEvent(e echo.Context) error
	// アサイン関連マスタ取得
	AssignMaster(e echo.Context) error
	// 面接官割り振り方法更新
	UpdateAssignMethod(e echo.Context) error
	// 書類提出ルールマスタ取得
	DocumentRuleMaster(e echo.Context) error
	// 職種マスタ取得
	OccupationMaster(e echo.Context) error
}

type UserController struct {
	s     service.IUserService
	a     service.IApplicantService
	login service.ILoginService
	role  service.IRoleService
}

func NewUserController(
	s service.IUserService,
	a service.IApplicantService,
	login service.ILoginService,
	role service.IRoleService,
) IUserController {
	return &UserController{s, a, login, role}
}

func (c *UserController) GetLoginService() service.ILoginService {
	return c.login
}

// 登録
func (c *UserController) Create(e echo.Context) error {
	req := request.CreateUser{}
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

	// ログイン種別取得
	loginType, loginTypeErr := c.login.GetLoginType(&request.GetLoginType{
		User: ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: req.HashKey,
			},
		},
	})
	if loginTypeErr != nil {
		return e.JSON(loginTypeErr.Status, response.ErrorConvert(*loginTypeErr))
	}

	id := static.ROLE_ADMIN_USER_CREATE
	if loginType.LoginType == static.LOGIN_TYPE_MANAGEMENT {
		id = static.ROLE_MANAGEMENT_USER_CREATE
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.HashKey,
		},
		ID: id,
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

	res, sErr := c.s.Create(&req)
	if sErr != nil {
		return e.JSON(sErr.Status, response.ErrorConvert(*sErr))
	}

	return e.JSON(http.StatusOK, res)
}

// 検索
func (c *UserController) Search(e echo.Context) error {
	req := request.SearchUser{}
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

	// ログイン種別取得
	loginType, loginTypeErr := c.login.GetLoginType(&request.GetLoginType{
		User: ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: req.HashKey,
			},
		},
	})
	if loginTypeErr != nil {
		return e.JSON(loginTypeErr.Status, response.ErrorConvert(*loginTypeErr))
	}

	id := static.ROLE_ADMIN_USER_READ
	if loginType.LoginType == static.LOGIN_TYPE_MANAGEMENT {
		id = static.ROLE_MANAGEMENT_USER_READ
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.HashKey,
		},
		ID: id,
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

// 検索_同一企業
func (c *UserController) SearchByCompany(e echo.Context) error {
	req := request.SearchUserByCompany{}
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

	// ログイン種別取得
	loginType, loginTypeErr := c.login.GetLoginType(&request.GetLoginType{
		User: ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: req.HashKey,
			},
		},
	})
	if loginTypeErr != nil {
		return e.JSON(loginTypeErr.Status, response.ErrorConvert(*loginTypeErr))
	}

	id := static.ROLE_ADMIN_USER_READ
	if loginType.LoginType == static.LOGIN_TYPE_MANAGEMENT {
		id = static.ROLE_MANAGEMENT_USER_READ
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.HashKey,
		},
		ID: id,
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

// 応募者ステータス変更
func (c *UserController) UpdateStatus(e echo.Context) error {
	req := request.UpdateStatus{}
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

	if err := c.a.UpdateStatus(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// ステータスイベントマスタ一覧
func (c *UserController) ListStatusEvent(e echo.Context) error {
	res, err := c.s.ListStatusEvent()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// アサイン関連マスタ取得
func (c *UserController) AssignMaster(e echo.Context) error {
	res, err := c.s.AssignMaster()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 面接官割り振り方法更新
func (c *UserController) UpdateAssignMethod(e echo.Context) error {
	req := request.UpdateAssignMethod{}
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

	if err := c.s.UpdateAssignMethod(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// 書類提出ルールマスタ取得
func (c *UserController) DocumentRuleMaster(e echo.Context) error {
	res, err := c.s.DocumentRuleMaster()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 職種マスタ取得
func (c *UserController) OccupationMaster(e echo.Context) error {
	res, err := c.s.OccupationMaster()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}
