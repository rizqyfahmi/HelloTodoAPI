package controllers

import (
	"TodoAPI/app/auth/entities"
	"TodoAPI/app/auth/usecases"
	"net/http"
	"time"

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
	e.POST("/refresh-token", handler.RefreshToken)
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

	ctx := c.Request().Context()

	result, err := r.authUsecase.Login(ctx, username, password)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
		})
	}

	accessToken := result["accessToken"].(map[string]interface{})
	r.authUsecase.SetCookie(c, accessToken["name"].(string), accessToken["token"].(string), accessToken["expiration"].(time.Time))

	refreshToken := result["refreshToken"].(map[string]interface{})
	r.authUsecase.SetCookie(c, refreshToken["name"].(string), refreshToken["token"].(string), refreshToken["expiration"].(time.Time))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successfully",
	})
}

func (r *AuthHandler) RefreshToken(c echo.Context) error {

	_, err := r.authUsecase.RefreshToken(c)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Refresh token successfully",
	})
}
