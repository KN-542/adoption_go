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
	team controller.ITeamController,
	schedule controller.IScheduleController,
	company controller.ICompanyController,
	applicant controller.IApplicantController,
	manuscript controller.IManuscriptController,
	role controller.IRoleController,
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
	e.POST("/confirm_team_applicant", login.ConfirmTeamApplicant)
	e.POST("/login_applicant", login.LoginApplicant)
	e.POST("/mfa_applicant", login.MFAApplicant)
	e.POST("/decode_applicant", login.JWTDecodeApplicant)
	e.POST("/code_gen_applicant", login.CodeGenerateApplicant)
	e.POST("/logout_applicant", login.LogoutApplicant)

	// 共通
	e.GET("/health", common.HealthCheck)
	e.POST("/sidebar", common.Sidebar)
	e.POST("/roles", common.Roles)
	e.POST("/change_team", common.ChangeTeam)

	// ユーザー
	e.POST("/user/search_company", user.SearchByCompany)
	e.POST("/user/search", user.Search)
	e.POST("/user/create", user.Create)
	e.POST("/user/delete", user.Delete)

	// チーム
	e.POST("/team/create", team.Create)
	e.POST("/team/update", team.Update)
	e.POST("/team/delete", team.Delete)
	e.POST("/team/search", team.Search)
	e.POST("/team/search_company", team.SearchByCompany)
	e.POST("/team/get", team.Get)

	// 予定
	e.POST("/schedule/type", schedule.SearchScheduleType)
	e.POST("/schedule/create", schedule.Insert)
	e.POST("/schedule/update", schedule.Update)
	e.POST("/schedule/search", schedule.Search)
	e.POST("/schedule/delete", schedule.Delete)

	// 企業
	e.POST("/company/create", company.Create)
	e.POST("/company/search", company.Search)

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
	e.POST("/applicant/assign_user", applicant.AssignUser)
	e.POST("/applicant/check_assign_user", applicant.CheckAssignableUser)
	e.POST("/applicant/types", applicant.ListApplicantTypeByTeam)
	e.POST("/applicant/update_type", applicant.CreateApplicantTypeAssociation)
	e.POST("/applicant/update_status", applicant.UpdateSelectStatus)
	e.POST("/applicant/result", applicant.InputResult)

	// ロール
	e.POST("/role/search_company", role.SearchByCompanyID)

	// 原稿
	e.POST("/manuscript/search", manuscript.Search)
	e.POST("/manuscript/search_by_team", manuscript.SearchManuscriptByTeam)
	e.POST("/manuscript/create", manuscript.Create)
	e.POST("/manuscript/assign_applicant", manuscript.CreateApplicantAssociation)
	e.POST("/manuscript/delete", manuscript.Delete)

	// 設定
	e.POST("/setting/get_team", team.GetOwn)
	e.POST("/setting/update_team", team.UpdateBasic)
	e.POST("/setting/team", user.UpdateStatus)
	e.POST("/setting/status_events", user.ListStatusEvent)
	e.POST("/setting/status_events_of_team", team.StatusEvents)
	e.POST("/setting/assign_masters", user.AssignMaster)
	e.POST("/setting/processing_list", team.ListInterviewProcessing)
	e.POST("/setting/update_assign", user.UpdateAssignMethod)
	e.POST("/setting/document_rules", user.DocumentRuleMaster)
	e.POST("/setting/occupations", user.OccupationMaster)
	e.POST("/setting/create_applicant_type", applicant.CreateApplicantType)
	e.POST("/setting/applicant_types", applicant.ListApplicantType)

	return e
}
