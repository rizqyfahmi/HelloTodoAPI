package middlewares

import (
	"TodoAPI/app/auth/usecases"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/labstack/echo/v4"
)

type Middleware struct{}

func InitMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) HandleCORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

func (m *Middleware) VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		cookie, errCookie := c.Cookie(usecases.GetAccessTokenName())

		if errCookie != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": errCookie.Error(),
			})
		}

		token, errParse := jwt.ParseWithClaims(cookie.Value, &usecases.AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(usecases.GetJWTSecret()), nil
		})

		if errParse != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": errParse.Error(),
			})
		}

		_, ok := token.Claims.(*usecases.AuthClaims)

		if !ok || !token.Valid {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "invalid token",
			})
		}

		return next(c)
	}
}
