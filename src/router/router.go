package router

import (
	"api/src/controller"
	m "api/src/middleware"
	"os"

	// echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(
	login controller.ILoginController,
	user controller.IUserController,
	applicant controller.IApplicantController,
) *echo.Echo {
	e := echo.New()

	// CORSミドルウェアの設定。認証情報を含むリクエストを許可する。
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FE_CSR_URL"), os.Getenv("FE_SSR_URL")},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true, // 認証情報を含むリクエストを許可
	}))

	e.POST("/login", login.Login)

	u := e.Group("/user")
	u.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FE_CSR_URL"), os.Getenv("FE_SSR_URL")},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))
	u.Use(m.JWTMiddleware)
	u.POST("/list", user.List)
	u.POST("/create", user.Create)
	u.POST("/role_list", user.RoleList)

	a := e.Group("/applicant")
	a.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FE_CSR_URL"), os.Getenv("FE_SSR_URL")},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))
	a.Use(m.JWTMiddleware)
	// a.Use(echojwt.WithConfig(echojwt.Config{
	// 	SigningKey:  []byte(os.Getenv("JWT_SECRET")),
	// 	TokenLookup: "cookie:token",
	// }))
	a.POST("/get_url", applicant.GetOauthURL)
	a.POST("/get_sheets", applicant.GetSheets)
	a.POST("/download", applicant.Download)
	a.POST("/search", applicant.Search)

	return e
}
