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
		user.MFA = int8(enum.MFA_UNAUTHENTICATED)
		return e.JSON(http.StatusOK, user)
	}

	// パスワード変更 必要性
	passChangeFlg, err := c.s.PasswordChangeCheck(&model.User{HashKey: user.HashKey})
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	if *passChangeFlg == int8(enum.PASSWORD_CHANGE_UNREQUIRED) {
		// JWT＆Cookie
		cookie, err := c.s.JWT(&req.Email, 0.1, "jwt_token")
		if err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
		e.SetCookie(cookie)

		user.PasswordChange = int8(enum.PASSWORD_CHANGE_UNREQUIRED)
		return e.JSON(http.StatusOK, user)
	}

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
		cookie, err := c.s.JWT(&req.Email, 0.1, "jwt_token")
		if err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
		e.SetCookie(cookie)

		// JWT(MFA)＆Cookie
		cookie2, err := c.s.JWT(&req.Email, 24, "jwt_token2")
		if err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
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
	req := model.User{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

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
	if err := c.s.UserCheck(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
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
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// パスワード変更
	if err := c.s.PasswordChange(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

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
