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
	"os"

	"github.com/labstack/echo/v4"
)

// 任意の構造体からILoginServiceフィールドを取得
func getServiceFromController[T any](c *T) (service.ILoginService, error) {
	type HasLoginService interface {
		GetLoginService() service.ILoginService
	}

	ctrl, ok := any(c).(HasLoginService)
	if !ok {
		return nil, fmt.Errorf(static.MESSAGE_NOT_FOUND_LOGIN_SERVICE)
	}

	return ctrl.GetLoginService(), nil
}

// JWT検証_共通化
func JWTDecodeCommon[T any](c *T, e echo.Context, hash_key string, token string, secret string, isUser bool) *response.Error {
	// Go単体で動作確認したい場合はGO_ENVをlocalに
	if os.Getenv("GO_ENV") == "local" {
		return nil
	}

	s, err := getServiceFromController(c)
	if err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusUnauthorized,
		}
	}

	cookie, err2 := e.Cookie(token)
	if err2 != nil {
		log.Printf(static.MESSAGE_UNEXPECTED_COOKIE)
		log.Printf("%v", err2)
		return &response.Error{
			Status: http.StatusUnauthorized,
		}
	}
	if err := s.JWTDecode(cookie, secret); err != nil {
		log.Printf("%v", err)
		return err
	}

	if isUser {
		// ユーザーが削除されていないかの確認
		if err := s.UserCheck(&request.JWTDecode{
			User: ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: hash_key,
				},
			},
		}); err != nil {
			return err
		}
	} else {
		// 応募者チェック
		if err := s.CheckApplicant(&request.CheckApplicant{
			Applicant: ddl.Applicant{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: hash_key,
				},
			},
		}); err != nil {
			return err
		}
	}

	// JWT＆Cookie 更新
	cookie, err3 := s.JWT(&hash_key, token, secret)
	if err3 != nil {
		return err3
	}
	e.SetCookie(cookie)

	return nil
}
