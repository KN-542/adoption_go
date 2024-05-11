package controller

import (
	"api/resources/static"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/service"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ICompanyController interface {
	// 登録
	Create(e echo.Context) error
}

type CompanyController struct {
	company service.ICompanyService
	login   service.ILoginService
}

func NewCompanyController(
	company service.ICompanyService,
	login service.ILoginService,
) ICompanyController {
	return &CompanyController{company, login}
}

func (c *CompanyController) GetLoginService() service.ILoginService {
	return c.login
}

// 登録
func (c *CompanyController) Create(e echo.Context) error {
	req := request.CompanyCreate{}
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
	); err != nil {
		return err
	}

	user, err := c.company.Create(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, user)
}
