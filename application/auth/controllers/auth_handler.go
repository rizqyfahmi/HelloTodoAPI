package controllers

import (
	"TodoAPI/application/auth/entities"
	"TodoAPI/application/auth/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase usecases.AuthUsecaseProtocol
}

func InitAuthHandler(e *echo.Echo, authUsecase usecases.AuthUsecaseProtocol) {
	handler := &AuthHandler{
		authUsecase: authUsecase,
	}

	e.POST("/registration", handler.Registration)
	e.POST("/login", handler.Login)
	e.POST("/logout", handler.Logout)
	e.POST("/refresh-token", handler.RefreshToken)
	e.POST("/validate-token", handler.ValidateToken)
}

func (r *AuthHandler) Registration(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	user := entities.User{
		Username: username,
		Password: []byte(password),
	}

	ctx := c.Request().Context()

	err := r.authUsecase.Registration(ctx, &user)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Store successfully",
	})
}

func (r *AuthHandler) Login(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	err := r.authUsecase.Login(c, username, password)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successfully",
	})
}

func (r *AuthHandler) Logout(c echo.Context) error {
	r.authUsecase.Logout(c)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Logout successfully",
	})
}

func (r *AuthHandler) RefreshToken(c echo.Context) error {

	err := r.authUsecase.RefreshToken(c)

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Refresh token successfully",
	})
}

func (r *AuthHandler) ValidateToken(c echo.Context) error {
	err := r.authUsecase.ValidateToken(c)

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Validate token successfully",
	})
}
