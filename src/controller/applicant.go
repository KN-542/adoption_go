package controller

import (
	"api/resources/static"
	"api/src/model"
	"api/src/service"
	"fmt"
	"log"
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
	// 検索
	Search(e echo.Context) error
}

type ApplicantController struct {
	s service.IApplicantService
}

func NewApplicantController(s service.IApplicantService) IApplicantController {
	return &ApplicantController{s}
}

/*
	OAuth2.0用(削除予定)
*/
// 認証URL作成
func (c *ApplicantController) GetOauthURL(e echo.Context) error {
	res, err := c.s.GetOauthURL()
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// シート取得
func (c *ApplicantController) GetSheets(e echo.Context) error {
	request := model.ApplicantSearch{}
	if err := e.Bind(&request); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	res, err := c.s.GetSheets(request)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

/*
	txt、csvダウンロード用
*/
// 応募者ダウンロード
func (c *ApplicantController) Download(e echo.Context) error {
	request := model.ApplicantsDownload{}
	if err := e.Bind(&request); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	res, err := c.s.Download(&request)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// 検索
func (c *ApplicantController) Search(e echo.Context) error {
	res, err := c.s.Search()
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}
