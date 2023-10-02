package controller

import (
	"api/src/model"
	"api/src/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IApplicantController interface {
	/*
		OAuth2.0用(削除予定)
	*/
	// 認証URL作成
	GetOauthURL(e echo.Context) error
	// シート取得
	GetSheets(e echo.Context) error
	/*
		txt、csvダウンロード用
	*/
	// 応募者ダウンロード
	Download(e echo.Context) error
}

type applicantController struct {
	s service.IApplicantService
}

func NewApplicantController(s service.IApplicantService) IApplicantController {
	return &applicantController{s}
}

/*
	OAuth2.0用(削除予定)
*/
// 認証URL作成
func (c *applicantController) GetOauthURL(e echo.Context) error {
	res, err := c.s.GetOauthURL()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, res)
}

// シート取得
func (c *applicantController) GetSheets(e echo.Context) error {
	request := model.ApplicantSearch{}
	if err := e.Bind(&request); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := c.s.GetSheets(request)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, res)
}

/*
	txt、csvダウンロード用
*/
// シート取得
func (c *applicantController) Download(e echo.Context) error {
	request := model.ApplicantsDownload{}
	if err := e.Bind(&request); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := c.s.Download(&request)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, res)
}
