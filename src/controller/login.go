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

const JWT_TOKEN string = "jwt_token"
const JWT_SECRET string = "JWT_SECRET"
const JWT_TOKEN2 string = "jwt_token2"
const JWT_SECRET2 string = "JWT_SECRET2"

type ILoginController interface {
	// ログイン
	Login(e echo.Context) error
	// MFA 認証コード生成
	CodeGenerate(e echo.Context) error
	// MFA
	MFA(e echo.Context) error
	// JWT 検証
	JWTDecode(e echo.Context) error
	// パスワード変更
	PasswordChange(e echo.Context) error
	// ログアウト
	Logout(e echo.Context) error
	// ログイン(応募者)
	LoginApplicant(e echo.Context) error
	// MFA 認証コード生成(応募者)
	CodeGenerateApplicant(e echo.Context) error
	// MFA(応募者)
	MFAApplicant(e echo.Context) error
	// JWT 検証(応募者)
	JWTDecodeApplicant(e echo.Context) error
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

// ログイン
func (c *LoginController) Login(e echo.Context) error {
	req := request.Login{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// ログイン
	user, sErr := c.s.Login(&req)
	if sErr != nil {
		return e.JSON(sErr.Status, response.ErrorConvert(*sErr))
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
	mfa, sErr := c.s.MFA(&req)
	if sErr != nil {
		return e.JSON(sErr.Status, response.ErrorConvert(*sErr))
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
	if err := JWTDecodeCommon(c, e, req.HashKey, JWT_TOKEN, JWT_SECRET, true); err != nil {
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
	req := request.LoginApplicant{}
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

// MFA 認証コード生成(応募者)
func (c *LoginController) CodeGenerateApplicant(e echo.Context) error {
	req := request.CodeGenerateApplicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.CodeGenerateApplicant(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// MFA(応募者)
func (c *LoginController) MFAApplicant(e echo.Context) error {
	req := request.MFAApplicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// MFA
	if err := c.s.MFAApplicant(&req); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	// JWT＆Cookie
	cookie, err := c.s.JWT(&req.HashKey, JWT_TOKEN2, JWT_SECRET2)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}

// JWT 検証(応募者)
func (c *LoginController) JWTDecodeApplicant(e echo.Context) error {
	req := request.JWTDecodeApplicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// JWT検証
	if err := JWTDecodeCommon(c, e, req.HashKey, JWT_TOKEN2, JWT_SECRET2, false); err != nil {
		return err
	}

	return e.JSON(http.StatusOK, "OK")
}

// ログアウト(応募者)
func (c *LoginController) LogoutApplicant(e echo.Context) error {
	req := request.LogoutApplicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	cookie, err := c.s.LogoutApplicant(&req, JWT_TOKEN2)
	if err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}
