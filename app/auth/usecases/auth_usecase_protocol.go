package usecases

import (
	"TodoAPI/app/auth/entities"
	"context"

	"github.com/labstack/echo/v4"
)

type AuthUsecaseProtocol interface {
	Registration(c context.Context, user *entities.User) error
	Login(c echo.Context, username, password string) error
	RefreshToken(c echo.Context) error
	// SetCookie(c echo.Context, name, token string, expiration time.Time)
	// GenerateAccessToken(user *entities.User) (string, time.Time, error)
	// GenerateRefreshToken(user *entities.User) (string, time.Time, error)
}
