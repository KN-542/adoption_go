package router

import (
	"api/src/controller"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(
	login controller.ILoginController,
	user controller.IUserController,
	applicant controller.IApplicantController,
) *echo.Echo {
	e := echo.New()

	// CORSミドルウェアの設定。認証情報を含むリクエストを許可
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FE_CSR_URL"), os.Getenv("FE_SSR_URL")},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true, // 認証情報を含むリクエストを許可
	}))

	e.POST("/login", login.Login)
	e.POST("/code_gen", login.CodeGenerate)
	e.POST("/mfa", login.MFA)
	e.POST("/decode", login.JWTDecode)
	e.POST("/password_change", login.PasswordChange)
	e.POST("/session_confirm", login.SessionConfirm)

	e.POST("/user/list", user.List)
	e.POST("/user/create", user.Create)
	e.POST("/user/role_list", user.RoleList)

	e.POST("/applicant/get_url", applicant.GetOauthURL)
	e.POST("/applicant/get_sheets", applicant.GetSheets)
	e.POST("/applicant/download", applicant.Download)
	e.POST("/applicant/search", applicant.Search)

	return e
}
