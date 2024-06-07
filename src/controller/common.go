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

type ICommonController interface {
	// サイドバー表示
	Sidebar(e echo.Context) error
	// 使用可能ロール一覧
	Roles(e echo.Context) error
	// チーム切り替え
	// ChangeTeam(e echo.Context) error
}

type CommonController struct {
	common service.ICommonService
	login  service.ILoginService
}

func NewCommonController(
	common service.ICommonService,
	login service.ILoginService,
) ICommonController {
	return &CommonController{common, login}
}

func (c *CommonController) GetLoginService() service.ILoginService {
	return c.login
}

// サイドバー表示
func (c *CommonController) Sidebar(e echo.Context) error {
	req := request.Sidebar{}
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
	); err != nil {
		return err
	}

	res, sErr := c.common.Sidebar(&req)
	if sErr != nil {
		return e.JSON(sErr.Status, response.ErrorConvert(*sErr))
	}

	return e.JSON(http.StatusOK, res)
}

// 使用可能ロール一覧
func (c *CommonController) Roles(e echo.Context) error {
	req := request.Roles{}
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
	); err != nil {
		return err
	}

	res, sErr := c.common.Roles(&req)
	if sErr != nil {
		return e.JSON(sErr.Status, response.ErrorConvert(*sErr))
	}

	return e.JSON(http.StatusOK, res)
}
