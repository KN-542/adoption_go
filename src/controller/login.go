package controller

import (
	"api/resources/static"
	"api/src/infra"
	"api/src/model/ddl"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/service"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

const JWT_TOKEN string = "jwt_token"
const JWT_SECRET string = "JWT_SECRET"

type ILoginController interface {
	// Hello World Test By Batch
	HelloWorld(e echo.Context) error
	// ログイン
	Login(e echo.Context) error
	// MFA 認証コード生成
	CodeGenerate(e echo.Context) error
	// MFA
	MFA(e echo.Context) error
	// JWT 検証
	JWTDecode(e echo.Context) error
	// セッション存在確認
	SessionConfirm(e echo.Context) error
	// パスワード変更
	PasswordChange(e echo.Context) error
	// ログアウト
	Logout(e echo.Context) error
	// ログイン(応募者)
	LoginApplicant(e echo.Context) error
	// MFA Applicant
	MFAApplicant(e echo.Context) error
	// JWT 検証(応募者)
	JWTDecodeApplicant(e echo.Context) error
	// MFA 認証コード生成(応募者)
	CodeGenerateApplicant(e echo.Context) error
	// ログアウト(応募者)
	LogoutApplicant(e echo.Context) error
}

type LoginController struct {
	s service.ILoginService
}

func NewLoginController(s service.ILoginService) ILoginController {
	return &LoginController{s}
}

func (c *LoginController) GetLoginService() service.ILoginService {
	return c.s
}

// Hello World Test By Batch
func (c *LoginController) HelloWorld(e echo.Context) error {
	fmt.Println("Hello World!")

	return e.JSON(http.StatusOK, "OK")
}

// ログイン
func (c *LoginController) Login(e echo.Context) error {
	req := request.Login{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// ログイン
	user, err := c.s.Login(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, user)
}

// MFA 認証コード生成
func (c *LoginController) CodeGenerate(e echo.Context) error {
	req := request.CodeGenerate{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.CodeGenerate(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// MFA
func (c *LoginController) MFA(e echo.Context) error {
	req := request.MFA{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// MFA
	mfa, err := c.s.MFA(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	if !mfa.IsPasswordChange {
		// JWT＆Cookie
		cookie, err := c.s.JWT(&req.HashKey, JWT_TOKEN, JWT_SECRET)
		if err != nil {
			return e.JSON(err.Status, response.ErrorConvert(*err))
		}
		e.SetCookie(cookie)
	}

	return e.JSON(http.StatusOK, mfa)
}

// JWT 検証
func (c *LoginController) JWTDecode(e echo.Context) error {
	req := request.JWTDecode{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(c, e, req.HashKey, JWT_TOKEN, JWT_SECRET); err != nil {
		return err
	}

	return e.JSON(http.StatusOK, "OK")
}

// パスワード変更
func (c *LoginController) PasswordChange(e echo.Context) error {
	req := request.PasswordChange{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// パスワード変更
	if err := c.s.PasswordChange(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	// JWT＆Cookie
	cookie, err := c.s.JWT(&req.HashKey, JWT_TOKEN, JWT_SECRET)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}

// セッション存在確認
func (c *LoginController) SessionConfirm(e echo.Context) error {
	req := ddl.User{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.SessionConfirm(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, "OK")
}

// ログアウト
func (c *LoginController) Logout(e echo.Context) error {
	req := request.Logout{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	cookie, err := c.s.Logout(&req, JWT_TOKEN)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}

// ログイン(応募者)
func (c *LoginController) LoginApplicant(e echo.Context) error {
	req := ddl.Applicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// ログイン認証
	applicant, err := c.s.LoginApplicant(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, applicant)
}

// MFA Applicant
func (c *LoginController) MFAApplicant(e echo.Context) error {
	req := request.MFA{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// MFA
	_, err := c.s.MFA(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	// S3 Name Redisに登録
	if err := c.s.S3NamePreInsert(&ddl.Applicant{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	}); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	// JWT＆Cookie
	cookie, err := c.s.JWT(&req.HashKey, "jwt_token3", "JWT_SECRET3")
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}

// JWT 検証(応募者)
func (c *LoginController) JWTDecodeApplicant(e echo.Context) error {
	req := ddl.Applicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// ログインJWT
	status, err := infra.JWTLoginToken(e, "jwt_token3", "JWT_SECRET3")
	if err != nil {
		log.Printf("%v", err)
		return e.JSON(status, response.ErrorCode{Code: static.CODE_LOGIN_REQUIRED})
	}

	// 応募者が削除されていないかの確認
	applicant, err2 := c.s.UserCheckApplicant(&req)
	if err2 != nil {
		return e.JSON(err2.Status, response.ErrorConvert(*err2))
	}

	// Redis 有効期限更新
	if err := c.s.SessionConfirmApplicant(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	// JWT＆Cookie 更新
	cookie, err3 := c.s.JWT(&req.HashKey, "jwt_token3", "JWT_SECRET3")
	if err3 != nil {
		return e.JSON(err3.Status, response.ErrorConvert(*err3))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, applicant)
}

// MFA 認証コード生成(応募者)
func (c *LoginController) CodeGenerateApplicant(e echo.Context) error {
	req := ddl.Applicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.CodeGenerateApplicant(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// ログアウト(応募者)
func (c *LoginController) LogoutApplicant(e echo.Context) error {
	req := ddl.Applicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	cookie, err := c.s.LogoutApplicant(&req)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}
