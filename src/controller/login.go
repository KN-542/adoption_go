package controller

import (
	"api/resources/static"
	"api/src/infra"
	"api/src/model"
	enum "api/src/model/enum"
	"api/src/service"
	"fmt"
	"log"
	"net/http"
	"time"

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
	// セッション存在確認
	SessionConfirm(e echo.Context) error
	// パスワード変更
	PasswordChange(e echo.Context) error
	// ログアウト
	Logout(e echo.Context) error
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
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// ログイン
	user, err := c.s.Login(&req)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	// MFA JWT
	if err := infra.JWTMFAToken(e); err != nil {
		return e.JSON(http.StatusOK, user)
	}

	// パスワード変更 必要性
	passChangeFlg, err := c.s.PasswordChangeCheck(&model.User{HashKey: user.HashKey})
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	if *passChangeFlg == int8(enum.PASSWORD_CHANGE_UNREQUIRED) {
		// JWT＆Cookie
		cookie, err := c.s.JWT(&user.HashKey, 1*time.Hour, "jwt_token", "JWT_SECRET")
		if err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
		e.SetCookie(cookie)

		// JWT(MFA)＆Cookie
		cookie2, err := c.s.JWT(&user.HashKey, 2*time.Hour, "jwt_token2", "JWT_SECRET2")
		if err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
		e.SetCookie(cookie2)

		user.MFA = int8(enum.MFA_AUTHENTICATED)
		return e.JSON(http.StatusOK, user)
	}

	user.MFA = int8(enum.MFA_AUTHENTICATED)
	user.PasswordChange = int8(enum.PASSWORD_CHANGE_REQUIRED)
	return e.JSON(http.StatusOK, user)
}

// MFA 認証コード生成
func (c *LoginController) CodeGenerate(e echo.Context) error {
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.CodeGenerate(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// MFA
func (c *LoginController) MFA(e echo.Context) error {
	req := model.UserMFA{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// MFA
	if err := c.s.MFA(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	// パスワード変更 必要性
	passChangeFlg, err := c.s.PasswordChangeCheck(&model.User{HashKey: req.HashKey})
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	if *passChangeFlg == int8(enum.PASSWORD_CHANGE_UNREQUIRED) {
		// JWT＆Cookie
		cookie, err := c.s.JWT(&req.HashKey, 1*time.Hour, "jwt_token", "JWT_SECRET")
		if err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
		e.SetCookie(cookie)

		// JWT(MFA)＆Cookie
		cookie2, err := c.s.JWT(&req.HashKey, 2*time.Hour, "jwt_token2", "JWT_SECRET2")
		if err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
		e.SetCookie(cookie2)

		return e.JSON(
			http.StatusOK,
			model.LoginStatusResponse{},
		)
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
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// ログインJWT
	status, err := infra.JWTLoginToken(e)
	if err != nil {
		log.Printf("%v", err)
		return e.JSON(status, model.ErrorCodeResponse{Code: static.CODE_LOGIN_REQUIRED})
	}

	// MFA JWT
	if err := infra.JWTMFAToken(e); err != nil {
		return e.JSON(
			http.StatusOK,
			model.LoginStatusResponse{},
		)
	}

	// ユーザーが削除されていないかの確認
	user, err2 := c.s.UserCheck(&req)
	if err2 != nil {
		return e.JSON(err2.Status, model.ErrorConvert(*err2))
	}

	// Redis 有効期限更新
	if err := c.s.SessionConfirm(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	// JWT＆Cookie 更新
	cookie, err3 := c.s.JWT(&req.HashKey, 1*time.Hour, "jwt_token", "JWT_SECRET")
	if err3 != nil {
		return e.JSON(err3.Status, model.ErrorConvert(*err3))
	}
	e.SetCookie(cookie)

	user.MFA = int8(enum.MFA_AUTHENTICATED)
	return e.JSON(http.StatusOK, user)
}

// パスワード変更
func (c *LoginController) PasswordChange(e echo.Context) error {
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// パスワード変更
	if err := c.s.PasswordChange(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	// JWT＆Cookie
	cookie, err := c.s.JWT(&req.HashKey, 1*time.Hour, "jwt_token", "JWT_SECRET")
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	// JWT(MFA)＆Cookie
	cookie2, err := c.s.JWT(&req.HashKey, 2*time.Hour, "jwt_token2", "JWT_SECRET2")
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	e.SetCookie(cookie2)

	return e.JSON(http.StatusOK, "OK")
}

// セッション存在確認
func (c *LoginController) SessionConfirm(e echo.Context) error {
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.SessionConfirm(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, "OK")
}

// ログアウト
func (c *LoginController) Logout(e echo.Context) error {
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	cookie, err := c.s.Logout(&req)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}
