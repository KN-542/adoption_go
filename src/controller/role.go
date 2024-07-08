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

type IRoleController interface {
	// 検索_企業ID
	SearchByCompanyID(e echo.Context) error
}

type RoleController struct {
	role  service.IRoleService
	login service.ILoginService
}

func NewRoleController(
	role service.IRoleService,
	login service.ILoginService,
) IRoleController {
	return &RoleController{role, login}
}

func (c *RoleController) GetLoginService() service.ILoginService {
	return c.login
}

// 検索_企業ID
func (c *RoleController) SearchByCompanyID(e echo.Context) error {
	req := request.SearchRoleByComapny{}
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
		ID: static.ROLE_MANAGEMENT_ROLE_READ,
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

	// 検索
	res, sErr := c.role.SearchRoleByComapny(&req)
	if sErr != nil {
		return e.JSON(sErr.Status, response.ErrorConvert(*sErr))
	}

	return e.JSON(http.StatusOK, res)
}
