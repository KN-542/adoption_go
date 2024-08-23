package controller

import (
	"api/src/model/ddl"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/service"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IApplicantController interface {
	// 検索
	Search(e echo.Context) error
	// サイト一覧取得
	GetSites(e echo.Context) error
	// 応募者ステータス一覧取得
	GetStatusList(e echo.Context) error
	// 応募者ダウンロード
	Download(e echo.Context) error
	// 予約表表示
	ReserveTable(e echo.Context) error
	// 書類アップロード
	DocumentsUpload(e echo.Context) error
	// 面接希望日登録
	InsertDesiredAt(e echo.Context) error
	// 認証URL作成
	GetOauthURL(e echo.Context) error
	// 応募者取得(1件)
	Get(e echo.Context) error
	// 書類ダウンロード
	DocumentDownload(e echo.Context) error
	// Google Meet Url 発行
	GetGoogleMeetUrl(e echo.Context) error
	// 面接官割り振り
	AssignUser(e echo.Context) error
	// 面接官割り振り可能判定
	CheckAssignableUser(e echo.Context) error
}

type ApplicantController struct {
	s     service.IApplicantService
	user  service.IUserService
	login service.ILoginService
	role  service.IRoleService
}

func NewApplicantController(
	s service.IApplicantService,
	user service.IUserService,
	login service.ILoginService,
	role service.IRoleService,
) IApplicantController {
	return &ApplicantController{s, user, login, role}
}

func (c *ApplicantController) GetLoginService() service.ILoginService {
	return c.login
}

// 検索
func (c *ApplicantController) Search(e echo.Context) error {
	req := request.SearchApplicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_APPLICANT_READ,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusNoContent,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	// 検索
	res, sErr := c.s.Search(&req)
	if sErr != nil {
		return e.JSON(sErr.Status, response.ErrorConvert(*sErr))
	}

	return e.JSON(http.StatusOK, res)
}

// サイト一覧取得
func (c *ApplicantController) GetSites(e echo.Context) error {
	res, err := c.s.GetSites()
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 応募者ステータス一覧取得
func (c *ApplicantController) GetStatusList(e echo.Context) error {
	req := request.ApplicantStatusList{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_APPLICANT_READ,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusNoContent,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	res, err := c.s.GetStatusList(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// 応募者ダウンロード
func (c *ApplicantController) Download(e echo.Context) error {
	req := request.ApplicantDownload{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_APPLICANT_DOWNLOAD,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusForbidden,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	res, err := c.s.Download(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// 予約表表示
func (c *ApplicantController) ReserveTable(e echo.Context) error {
	req := request.ReserveTable{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(c, e, req.HashKey, JWT_TOKEN2, JWT_SECRET2, false); err != nil {
		return err
	}

	res, err := c.s.ReserveTable(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// 認証URL作成
func (c *ApplicantController) GetOauthURL(e echo.Context) error {
	req := request.GetOauthURL{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_APPLICANT_READ,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusNoContent,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	res, err := c.s.GetOauthURL(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// 応募者取得(1件)
func (c *ApplicantController) Get(e echo.Context) error {
	req := request.GetApplicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_APPLICANT_READ,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusNoContent,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	res, err := c.s.Get(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, res)
}

// 書類アップロード
func (c *ApplicantController) DocumentsUpload(e echo.Context) error {
	hashKey := e.FormValue("hash_key")

	// JWT検証
	if err := JWTDecodeCommon(c, e, hashKey, JWT_TOKEN2, JWT_SECRET2, false); err != nil {
		return err
	}

	resumeExtension := e.FormValue("resume_extension")
	if resumeExtension != "" {
		resume, err := e.FormFile("resume")
		if err != nil {
			log.Printf("%v", err)
			return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
		}

		if err := c.s.S3Upload(&request.FileUpload{
			Applicant: ddl.Applicant{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: hashKey,
				},
			},
			Extension: resumeExtension,
			NamePre:   "resume",
		}, resume); err != nil {
			return e.JSON(err.Status, response.ErrorConvert(*err))
		}
	}

	curriculumVitaeExtension := e.FormValue("curriculum_vitae_extension")
	if curriculumVitaeExtension != "" {
		curriculumVitae, err := e.FormFile("curriculum_vitae")
		if err != nil {
			log.Printf("%v", err)
			return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
		}

		if err := c.s.S3Upload(&request.FileUpload{
			Applicant: ddl.Applicant{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: hashKey,
				},
			},
			Extension: curriculumVitaeExtension,
			NamePre:   "curriculum_vitae",
		}, curriculumVitae); err != nil {
			return e.JSON(err.Status, response.ErrorConvert(*err))
		}
	}

	return e.JSON(http.StatusOK, "OK")
}

// 書類ダウンロード
func (c *ApplicantController) DocumentDownload(e echo.Context) error {
	req := request.FileDownload{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_APPLICANT_READ,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusNoContent,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	file, fileName, err := c.s.S3Download(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	e.Response().Header().Set("Content-Disposition", "attachment; filename="+*fileName)
	e.Response().Header().Set("Content-Type", "application/octet-stream")
	return e.Blob(http.StatusOK, "application/octet-stream", file)
}

// 面接希望日登録
func (c *ApplicantController) InsertDesiredAt(e echo.Context) error {
	req := request.InsertDesiredAt{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(c, e, req.HashKey, JWT_TOKEN2, JWT_SECRET2, false); err != nil {
		return err
	}

	if err := c.s.InsertDesiredAt(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, "OK")
}

// Google Meet Url 発行
func (c *ApplicantController) GetGoogleMeetUrl(e echo.Context) error {
	req := request.GetGoogleMeetUrl{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_APPLICANT_READ,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusNoContent,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	res, err := c.s.GetGoogleMeetUrl(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}

// 面接官割り振り
func (c *ApplicantController) AssignUser(e echo.Context) error {
	req := request.AssignUser{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_APPLICANT_ASSIGN_USER,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusForbidden,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	if err := c.s.AssignUser(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// 面接官割り振り可能判定
func (c *ApplicantController) CheckAssignableUser(e echo.Context) error {
	req := request.CheckAssignableUser{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey,
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey,
		},
		ID: static.ROLE_MANAGEMENT_APPLICANT_ASSIGN_USER,
	})
	if roleErr != nil {
		return e.JSON(roleErr.Status, response.ErrorConvert(*roleErr))
	}
	if !exist {
		err := &response.Error{
			Status: http.StatusForbidden,
		}
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	res, err := c.s.CheckAssignableUser(&req, false)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, res)
}
