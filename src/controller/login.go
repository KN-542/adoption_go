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

// Hello World Test By Batch
func (c *LoginController) HelloWorld(e echo.Context) error {
	fmt.Println("Hello World!")

	return e.JSON(http.StatusOK, "OK")
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
		cookie, err := c.s.JWT(&user.HashKey, "jwt_token", "JWT_SECRET")
		if err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
		e.SetCookie(cookie)

		// JWT(MFA)＆Cookie
		cookie2, err := c.s.JWT(&user.HashKey, "jwt_token2", "JWT_SECRET2")
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
		cookie, err := c.s.JWT(&req.HashKey, "jwt_token", "JWT_SECRET")
		if err != nil {
			return e.JSON(err.Status, model.ErrorConvert(*err))
		}
		e.SetCookie(cookie)

		// JWT(MFA)＆Cookie
		cookie2, err := c.s.JWT(&req.HashKey, "jwt_token2", "JWT_SECRET2")
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
	status, err := infra.JWTLoginToken(e, "jwt_token", "JWT_SECRET")
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
	cookie, err3 := c.s.JWT(&req.HashKey, "jwt_token", "JWT_SECRET")
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
	cookie, err := c.s.JWT(&req.HashKey, "jwt_token", "JWT_SECRET")
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	// JWT(MFA)＆Cookie
	cookie2, err := c.s.JWT(&req.HashKey, "jwt_token2", "JWT_SECRET2")
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

// ログイン(応募者)
func (c *LoginController) LoginApplicant(e echo.Context) error {
	req := model.Applicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// ログイン認証
	applicant, err := c.s.LoginApplicant(&req)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	return e.JSON(http.StatusOK, applicant)
}

// MFA Applicant
func (c *LoginController) MFAApplicant(e echo.Context) error {
	req := model.UserMFA{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// MFA
	if err := c.s.MFA(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	// S3 Name Redisに登録
	if err := c.s.S3NamePreInsert(&model.Applicant{HashKey: req.HashKey}); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	// JWT＆Cookie
	cookie, err := c.s.JWT(&req.HashKey, "jwt_token3", "JWT_SECRET3")
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}

// JWT 検証(応募者)
func (c *LoginController) JWTDecodeApplicant(e echo.Context) error {
	req := model.Applicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	// ログインJWT
	status, err := infra.JWTLoginToken(e, "jwt_token3", "JWT_SECRET3")
	if err != nil {
		log.Printf("%v", err)
		return e.JSON(status, model.ErrorCodeResponse{Code: static.CODE_LOGIN_REQUIRED})
	}

	// 応募者が削除されていないかの確認
	applicant, err2 := c.s.UserCheckApplicant(&req)
	if err2 != nil {
		return e.JSON(err2.Status, model.ErrorConvert(*err2))
	}

	// Redis 有効期限更新
	if err := c.s.SessionConfirmApplicant(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}

	// JWT＆Cookie 更新
	cookie, err3 := c.s.JWT(&req.HashKey, "jwt_token3", "JWT_SECRET3")
	if err3 != nil {
		return e.JSON(err3.Status, model.ErrorConvert(*err3))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, applicant)
}

// MFA 認証コード生成(応募者)
func (c *LoginController) CodeGenerateApplicant(e echo.Context) error {
	req := model.Applicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	if err := c.s.CodeGenerateApplicant(&req); err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	return e.JSON(http.StatusOK, "OK")
}

// ログアウト(応募者)
func (c *LoginController) LogoutApplicant(e echo.Context) error {
	req := model.Applicant{}
	if err := e.Bind(&req); err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusBadRequest, fmt.Errorf(static.MESSAGE_BAD_REQUEST))
	}

	cookie, err := c.s.LogoutApplicant(&req)
	if err != nil {
		return e.JSON(err.Status, model.ErrorConvert(*err))
	}
	e.SetCookie(cookie)

	return e.JSON(http.StatusOK, "OK")
}
