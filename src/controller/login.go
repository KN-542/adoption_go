package controller

import (
	"api/src/infra"
	"api/src/model"
	enum "api/src/model/enum"
	"api/src/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
}

type LoginController struct {
	s service.ILoginService
}

func NewLoginController(s service.ILoginService) ILoginController {
	return &LoginController{s}
}

// ログイン
func (c *LoginController) Login(e echo.Context) error {
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	// ログイン
	if err := c.s.Login(&req); err != nil {
		return e.JSON(err.Status, err.Error.Error())
	}

	// MFA JWT
	if err := infra.JWTMFAToken(e); err != nil {
		return e.JSON(
			http.StatusOK,
			model.LoginStatusResponse{
				MFA: int8(enum.MFA_UNAUTHENTICATED),
			},
		)
	}

	// パスワード変更 必要性
	passChangeFlg, err := c.s.PasswordChangeCheck()
	if err != nil {
		return e.JSON(err.Status, err.Error.Error())
	}

	if *passChangeFlg == int8(enum.PASSWORD_CHANGE_UNREQUIRED) {
		// JWT＆Cookie
		cookie, err := c.s.JWT(&req.Email, 0.1, "jwt_token")
		if err != nil {
			return e.JSON(err.Status, err.Error.Error())
		}
		e.SetCookie(cookie)
	}

	return e.JSON(
		http.StatusOK,
		model.LoginStatusResponse{
			MFA:            int8(enum.MFA_UNAUTHENTICATED),
			PasswordChange: *passChangeFlg,
		},
	)
}

// MFA 認証コード生成
func (c *LoginController) CodeGenerate(e echo.Context) error {
	if err := c.s.CodeGenerate(); err != nil {
		return e.JSON(err.Status, err.Error.Error())
	}
	return e.JSON(http.StatusOK, "OK")
}

// MFA
func (c *LoginController) MFA(e echo.Context) error {
	req := model.UserMFA{}
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	// MFA
	if err := c.s.MFA(&req); err != nil {
		return e.JSON(err.Status, err.Error.Error())
	}

	// パスワード変更 必要性
	passChangeFlg, err := c.s.PasswordChangeCheck()
	if err != nil {
		return e.JSON(err.Status, err.Error.Error())
	}

	if *passChangeFlg == int8(enum.PASSWORD_CHANGE_UNREQUIRED) {
		// JWT＆Cookie
		cookie, err := c.s.JWT(nil, 0.1, "jwt_token")
		if err != nil {
			return e.JSON(err.Status, err.Error.Error())
		}
		e.SetCookie(cookie)

		// JWT(MFA)＆Cookie
		cookie2, err := c.s.JWT(nil, 24, "jwt_token2")
		if err != nil {
			return e.JSON(err.Status, err.Error.Error())
		}
		e.SetCookie(cookie2)
	}

	return e.JSON(
		http.StatusOK,
		model.LoginStatusResponse{
			PasswordChange: *passChangeFlg,
		},
	)
}

// JWT 検証
func (c *LoginController) JWTDecode(e echo.Context) error {
	// ログインJWT
	if err := infra.JWTLoginToken(e); err != nil {
		return err
	}

	// MFA JWT
	if err := infra.JWTMFAToken(e); err != nil {
		return e.JSON(
			http.StatusOK,
			model.LoginStatusResponse{
				MFA: int8(enum.MFA_UNAUTHENTICATED),
			},
		)
	}

	// ユーザーが削除されていないかの確認
	if err := c.s.UserCheck(); err != nil {
		return e.JSON(err.Status, err.Error.Error())
	}

	return e.JSON(
		http.StatusOK,
		model.LoginStatusResponse{
			MFA: int8(enum.MFA_UNAUTHENTICATED),
		},
	)
}

// パスワード変更
func (c *LoginController) PasswordChange(e echo.Context) error {
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	// パスワード変更
	if err := c.s.PasswordChange(&req); err != nil {
		return e.JSON(err.Status, err.Error.Error())
	}

	return e.JSON(http.StatusOK, "OK")
}
