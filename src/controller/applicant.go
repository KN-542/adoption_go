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
	// 応募者ダウンロード
	Download(e echo.Context) error
	// 応募者取得(1件)
	Get(e echo.Context) error
	// 検索
	Search(e echo.Context) error
	// 書類アップロード
	DocumentsUpload(e echo.Context) error
	// 書類ダウンロード
	DocumentDownload(e echo.Context) error
	// 面接希望日登録
	InsertDesiredAt(e echo.Context) error
	// 応募者ステータス一覧取得
	GetApplicantStatus(e echo.Context) error
	// サイト一覧取得
	GetSites(e echo.Context) error
	// Google Meet Url 発行
	GetGoogleMeetUrl(e echo.Context) error
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
	request := model.ApplicantAndUser{}
	if err := e.Bind(&request); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	res, err := c.s.GetOauthURL(&request)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// 応募者ダウンロード
func (c *ApplicantController) Download(e echo.Context) error {
	request := model.ApplicantsDownload{}
	if err := e.Bind(&request); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.Download(&request); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, "OK")
}

// 応募者取得(1件)
func (c *ApplicantController) Get(e echo.Context) error {
	request := model.Applicant{}
	if err := e.Bind(&request); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	res, err := c.s.Get(&request)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// 検索
func (c *ApplicantController) Search(e echo.Context) error {
	request := model.ApplicantSearchRequest{}
	if err := e.Bind(&request); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	res, err := c.s.Search(&request)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// 書類アップロード
func (c *ApplicantController) DocumentsUpload(e echo.Context) error {
	hashKey := e.FormValue("hash_key")

	resumeExtension := e.FormValue("resume_extension")
	if resumeExtension != "" {
		resume, err := e.FormFile("resume")
		if err != nil {
			log.Printf("%v", err)
			return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
		}

		if err := c.s.S3Upload(&model.FileUpload{
			HashKey:   hashKey,
			Extension: resumeExtension,
			NamePre:   "resume",
		}, resume); err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
	}

	curriculumVitaeExtension := e.FormValue("curriculum_vitae_extension")
	if curriculumVitaeExtension != "" {
		curriculumVitae, err := e.FormFile("curriculum_vitae")
		if err != nil {
			log.Printf("%v", err)
			return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
		}

		if err := c.s.S3Upload(&model.FileUpload{
			HashKey:   hashKey,
			Extension: curriculumVitaeExtension,
			NamePre:   "curriculum_vitae",
		}, curriculumVitae); err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
	}

	return e.JSON(http.StatusOK, "OK")
}

// 書類ダウンロード
func (c *ApplicantController) DocumentDownload(e echo.Context) error {
	request := model.FileDownload{}
	if err := e.Bind(&request); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	file, fileName, err := c.s.S3Download(&request)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	e.Response().Header().Set("Content-Disposition", "attachment; filename="+*fileName)
	e.Response().Header().Set("Content-Type", "application/octet-stream")
	return e.Blob(http.StatusOK, "application/octet-stream", file)
}

// 面接希望日登録
func (c *ApplicantController) InsertDesiredAt(e echo.Context) error {
	request := model.ApplicantDesired{}
	if err := e.Bind(&request); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.InsertDesiredAt(&request); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, "OK")
}

// 応募者ステータス一覧取得
func (c *ApplicantController) GetApplicantStatus(e echo.Context) error {
	res, err := c.s.GetApplicantStatus()
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// サイト一覧取得
func (c *ApplicantController) GetSites(e echo.Context) error {
	res, err := c.s.GetSites()
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// Google Meet Url 発行
func (c *ApplicantController) GetGoogleMeetUrl(e echo.Context) error {
	request := model.ApplicantAndUser{}
	if err := e.Bind(&request); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	res, err := c.s.GetGoogleMeetUrl(&request)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}
