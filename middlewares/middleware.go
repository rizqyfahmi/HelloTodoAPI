package middlewares

import "github.com/labstack/echo/v4"

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
