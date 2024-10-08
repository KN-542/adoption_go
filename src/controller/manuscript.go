package controller

import (
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/service"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IManuscriptController interface {
	// 検索
	Search(e echo.Context) error
	// 登録
	Create(e echo.Context) error
	// 応募者紐づけ登録
	CreateApplicantAssociation(e echo.Context) error
	// 検索_同一チーム
	SearchManuscriptByTeam(e echo.Context) error
	// 削除
	Delete(e echo.Context) error
}

type ManuscriptController struct {
	s     service.IManuscriptService
	login service.ILoginService
	role  service.IRoleService
}

func NewManuscriptController(
	s service.IManuscriptService,
	login service.ILoginService,
	role service.IRoleService,
) IManuscriptController {
	return &ManuscriptController{s, login, role}
}

func (c *ManuscriptController) GetLoginService() service.ILoginService {
	return c.login
}

// 検索
func (c *ManuscriptController) Search(e echo.Context) error {
	req := request.SearchManuscript{}
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
		ID: static.ROLE_MANAGEMENT_MANUSCRIPT_READ,
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

	res, sErr := c.s.Search(&req)
	if sErr != nil {
		return e.JSON(sErr.Status, response.ErrorConvert(*sErr))
	}

	return e.JSON(http.StatusOK, res)
}

// 登録
func (c *ManuscriptController) Create(e echo.Context) error {
	req := request.CreateManuscript{}
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
		ID: static.ROLE_MANAGEMENT_MANUSCRIPT_CREATE,
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

	if err := c.s.Create(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, "OK")
}

// 応募者紐づけ登録
func (c *ManuscriptController) CreateApplicantAssociation(e echo.Context) error {
	req := request.CreateApplicantAssociation{}
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
		ID: static.ROLE_MANAGEMENT_APPLICANT_SETTING_MANUSCRIPT,
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

	if err := c.s.CreateApplicantAssociation(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, "OK")
}

// 検索_同一チーム
func (c *ManuscriptController) SearchManuscriptByTeam(e echo.Context) error {
	req := request.SearchManuscriptByTeam{}
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

	res, sErr := c.s.SearchManuscriptByTeam(&req)
	if sErr != nil {
		return e.JSON(sErr.Status, response.ErrorConvert(*sErr))
	}

	return e.JSON(http.StatusOK, res)
}

// 複数の原稿IDの削除
func (c *ManuscriptController) Delete(e echo.Context) error {
	// リクエストボディをパースするための構造体
	type DeleteManuscriptRequest struct {
		UserHashKey       string   `json:"user_hash_key"`
		ManuscriptHashKey []string `json:"manuscript_hash_key"`
	}

	var req DeleteManuscriptRequest
	if err := e.Bind(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		return e.JSON(http.StatusBadRequest, "Invalid request body")
	}

	// JWT検証
	if err := JWTDecodeCommon(
		c,
		e,
		req.UserHashKey, // ここを修正
		JWT_TOKEN,
		JWT_SECRET,
		true,
	); err != nil {
		return err
	}

	// ロールチェック
	exist, roleErr := c.role.Check(&request.CheckRole{
		Abstract: request.Abstract{
			UserHashKey: req.UserHashKey, // ここを修正
		},
		ID: static.ROLE_MANAGEMENT_MANUSCRIPT_DELETE, // 削除権限のロールID
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

	// サービスで削除処理を実行
	if err := c.s.Delete(req.ManuscriptHashKey); err != nil { // req.ManuscriptIDs をそのまま使用
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, "OK")
}
