package router

import (
	"api/src/controller"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(
	common controller.ICommonController,
	login controller.ILoginController,
	user controller.IUserController,
	company controller.ICompanyController,
	applicant controller.IApplicantController,
) *echo.Echo {
	e := echo.New()

	// CORSミドルウェアの設定。認証情報を含むリクエストを許可
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			os.Getenv("FE_CSR_URL"),
			os.Getenv("FE_SSG_URL"),
			os.Getenv("FE_APPLICANT_CSR_URL"),
			os.Getenv("FE_APPLICANT_SSG_URL"),
		},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true, // 認証情報を含むリクエストを許可
	}))

	// ログイン
	e.POST("/login", login.Login)
	e.POST("/logout", login.Logout)
	e.POST("/code_gen", login.CodeGenerate)
	e.POST("/mfa", login.MFA)
	e.POST("/decode", login.JWTDecode)
	e.POST("/password_change", login.PasswordChange)
	e.POST("/login_applicant", login.LoginApplicant)
	e.POST("/mfa_applicant", login.MFAApplicant)
	e.POST("/decode_applicant", login.JWTDecodeApplicant)
	e.POST("/code_gen_applicant", login.CodeGenerateApplicant)
	e.POST("/logout_applicant", login.LogoutApplicant)

	// 共通
	e.POST("/sidebar", common.Sidebar)
	e.POST("/roles", common.Roles)

	// ユーザー
	e.POST("/user/search", user.Search)
	e.POST("/user/create", user.Create)
	e.POST("/user/create_team", user.InsertTeam)
	e.POST("/user/search_team", user.SearchTeam)
	e.POST("/user/schedule_type", user.SearchScheduleType)
	e.POST("/user/create_schedule", user.InsertSchedules)
	e.POST("/user/update_schedule", user.UpdateSchedule)
	e.POST("/user/schedules", user.SearchSchedule)
	e.POST("/user/delete_schedule", user.DeleteSchedule)

	// 企業
	e.POST("/company/create", company.Create)

	// 応募者
	e.POST("/applicant/get_url", applicant.GetOauthURL)
	e.POST("/applicant/download", applicant.Download)
	e.POST("/applicant/get", applicant.Get)
	e.POST("/applicant/search", applicant.Search)
	e.POST("/applicant/documents", applicant.DocumentsUpload)
	e.POST("/applicant/documents_download", applicant.DocumentDownload)
	e.POST("/applicant/desired", applicant.InsertDesiredAt)
	e.POST("/applicant/status", applicant.GetStatusList)
	e.POST("/applicant/sites", applicant.GetSites)
	e.POST("/applicant/get_google_meet_url", applicant.GetGoogleMeetUrl)
	e.POST("/applicant/reserve_table", applicant.ReserveTable)

	return e
}
