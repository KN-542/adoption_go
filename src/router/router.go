package router

import (
	"api/src/controller"

	"github.com/labstack/echo/v4"
)

func NewRouter(
	login controller.ILoginController,
	user controller.IUserController,
	applicant controller.IApplicantController,
) *echo.Echo {
	e := echo.New()
	e.POST("login", login.Login)
	e.POST("user/list", user.List)
	e.POST("user/create", user.Create)
	e.POST("applicant/get_url", applicant.GetOauthURL)
	e.POST("applicant/get_sheets", applicant.GetSheets)
	e.POST("applicant/download", applicant.Download)
	e.POST("applicant/search", applicant.Search)
	return e
}
