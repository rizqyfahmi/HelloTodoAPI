package usecases

import (
	"TodoAPI/application/auth/entities"
	"context"

	"github.com/labstack/echo/v4"
)

type AuthUsecaseProtocol interface {
	Registration(c context.Context, user *entities.User) error
	Login(c echo.Context, username, password string) error
	Logout(c echo.Context) error
	RefreshToken(c echo.Context) error
	ValidateToken(c echo.Context) error
}
