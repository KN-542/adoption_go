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
		AllowOrigins: []string{
			os.Getenv("BATCH_URL"),
			os.Getenv("FE_CSR_URL"),
			os.Getenv("FE_SSG_URL"),
			os.Getenv("FE_APPLICANT_CSR_URL"),
			os.Getenv("FE_APPLICANT_SSG_URL"),
		},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true, // 認証情報を含むリクエストを許可
	}))

	e.POST("/hello", login.HelloWorld)
	e.POST("/login", login.Login)
	e.POST("/logout", login.Logout)
	e.POST("/code_gen", login.CodeGenerate)
	e.POST("/mfa", login.MFA)
	e.POST("/decode", login.JWTDecode)
	e.POST("/password_change", login.PasswordChange)
	e.POST("/session_confirm", login.SessionConfirm)
	e.POST("/login_applicant", login.LoginApplicant)
	e.POST("/mfa_applicant", login.MFAApplicant)
	e.POST("/decode_applicant", login.JWTDecodeApplicant)
	e.POST("/code_gen_applicant", login.CodeGenerateApplicant)
	e.POST("/logout_applicant", login.LogoutApplicant)

	e.POST("/user/list", user.List)
	e.POST("/user/create", user.Create)
	e.POST("/user/create_group", user.InsertGroup)
	e.POST("/user/role_list", user.RoleList)
	e.POST("/user/search_group", user.SearchGroups)
	e.POST("/user/schedule_type", user.ListScheduleType)
	e.POST("/user/create_schedule", user.InsertSchedules)
	e.POST("/user/schedules", user.Schedules)
	e.POST("/user/delete_schedule", user.DeleteSchedule)
	e.POST("/user/reserve_table", user.DispReserveTable)

	e.POST("/applicant/get_url", applicant.GetOauthURL)
	e.POST("/applicant/download", applicant.Download)
	e.POST("/applicant/get", applicant.Get)
	e.POST("/applicant/search", applicant.Search)
	e.POST("/applicant/documents", applicant.DocumentsUpload)
	e.POST("/applicant/documents_download", applicant.DocumentDownload)
	e.POST("/applicant/desired", applicant.InsertDesiredAt)
	e.POST("/applicant/status", applicant.GetApplicantStatus)
	e.POST("/applicant/sites", applicant.GetSites)
	e.POST("/applicant/get_google_meet_url", applicant.GetGoogleMeetUrl)

	return e
}
