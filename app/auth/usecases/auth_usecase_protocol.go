package usecases

import (
	"TodoAPI/app/auth/entities"
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthUsecaseProtocol interface {
	Registration(c context.Context, user *entities.User) error
	Login(c context.Context, username, password string) (map[string]interface{}, error)
	RefreshToken(c echo.Context) (map[string]string, error)
	SetCookie(c echo.Context, name, token string, expiration time.Time)
	GenerateAccessToken(user *entities.User) (string, time.Time, error)
	GenerateRefreshToken(user *entities.User) (string, time.Time, error)
}
