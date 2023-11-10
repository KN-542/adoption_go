package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		cookie, err := e.Cookie("jwt_token")
		if err != nil {
			log.Printf("%v", err)
			return e.JSON(http.StatusUnauthorized, err.Error())
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
			return e.JSON(http.StatusUnauthorized, err.Error())
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// トークンが有効であれば、クレームを利用することができます
		} else {
			log.Print("Invalid token")
			return e.JSON(http.StatusUnauthorized, "Invalid token")
		}

		return next(e)
	}
}
