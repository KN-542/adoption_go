package infra

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTLoginToken(e echo.Context) (int, error) {
	cookie, err := e.Cookie("jwt_token")
	if err != nil {
		log.Printf("%v", err)
		return http.StatusUnauthorized, err
	}
	tokenString := cookie.Value

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// トークンが期待する署名方法（HMAC、RSAなど）を使用していることを確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Print("Unexpected jwt token")
			return nil, e.JSON(http.StatusUnauthorized, "Unexpected jwt token")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Printf("%v", err)
		return http.StatusUnauthorized, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// トークンが有効であれば、クレームを利用可能
	} else {
		log.Print("Invalid token")
		return http.StatusUnauthorized, err
	}

	return http.StatusOK, nil
}

func JWTMFAToken(e echo.Context) error {
	cookie, err := e.Cookie("jwt_token2")
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	tokenString := cookie.Value

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Print("Unexpected jwt token")
			return nil, fmt.Errorf("Unexpected jwt token")
		}
		return []byte(os.Getenv("JWT_SECRET2")), nil
	})
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	} else {
		log.Print("Invalid token")
		return fmt.Errorf("Invalid token")
	}

	return nil
}
