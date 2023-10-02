package router

import (
	"api/src/controller"

	"github.com/labstack/echo/v4"
)

func NewRouter(
	interviewer controller.IInterviewerController,
	applicant controller.IApplicantController,
) *echo.Echo {
	e := echo.New()
	e.POST("interviewer/create", interviewer.RegisterInterviewer)
	e.POST("interviewer/list", interviewer.Interviewers)
	e.POST("applicant/get_url", applicant.GetOauthURL)
	e.POST("applicant/get_sheets", applicant.GetSheets)
	e.POST("applicant/download", applicant.Download)
	return e
}
