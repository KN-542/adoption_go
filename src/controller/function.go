package controller

import (
	"api/resources/static"
	"api/src/model/ddl"
	"api/src/model/request"
	"api/src/model/response"
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
func JWTDecodeCommon[T any](c *T, e echo.Context, hash_key string, token string, secret string) error {
	// Go単体で動作確認したい場合はGO_ENVをlocalに
	if os.Getenv("GO_ENV") == "local" {
		return nil
	}

	s, err := getServiceFromController(c)
	if err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusUnauthorized, err)
	}

	cookie, err := e.Cookie(token)
	if err != nil {
		log.Printf("%v", err)
		return e.JSON(http.StatusUnauthorized, fmt.Errorf(static.MESSAGE_UNEXPECTED_COOKIE))
	}
	if err := s.JWTDecode(cookie, secret); err != nil {
		log.Printf("%v", err)
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	// ユーザーが削除されていないかの確認
	if err := s.UserCheck(&request.JWTDecode{
		User: ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: hash_key,
			},
		},
	}); err != nil {
		return e.JSON(err.Status, response.ErrorConvert(*err))
	}

	// JWT＆Cookie 更新
	cookie, err2 := s.JWT(&hash_key, token, secret)
	if err2 != nil {
		return e.JSON(err2.Status, response.ErrorConvert(*err2))
	}
	e.SetCookie(cookie)

	return nil
}
