package usecases

import (
	"TodoAPI/app/auth/entities"
	"TodoAPI/app/auth/repositories"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
	jwtSecretKey           = "some-secret-key"
	jwtRefreshSecretKey    = "some-refresh-secret-key"
)

type AuthClaims struct {
	User *entities.User `json:"user"`
	jwt.StandardClaims
}

type AuthUsecase struct {
	repo    repositories.AuthRepositoryProtocol
	timeout time.Duration
}

func InitAuthUsecase(repo repositories.AuthRepositoryProtocol, timeout time.Duration) AuthUsecaseProtocol {
	return &AuthUsecase{
		repo:    repo,
		timeout: timeout,
	}
}

func (usecase *AuthUsecase) Registration(c context.Context, user *entities.User) error {

	ctx, cancel := context.WithTimeout(c, usecase.timeout)
	defer cancel()

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	const min int64 = 1000
	const max int64 = 9999

	user.Id = rand.Int63n(max-min) + min
	user.Password = password
	user.UpdatedAt = time.Now()
	user.CreatedAt = time.Now()

	if err := usecase.repo.Store(ctx, user); err != nil {
		return err
	}

	return nil
}

func (usecase *AuthUsecase) Login(c context.Context, username, password string) (map[string]interface{}, error) {

	ctx, cancel := context.WithTimeout(c, usecase.timeout)
	defer cancel()

	user, err := usecase.repo.GetUser(ctx, username)

	if err != nil {
		return nil, err
	}

	errBcrypt := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if errBcrypt != nil {
		return nil, errBcrypt
	}

	accessToken, expAccessToken, errAccessToken := usecase.GenerateAccessToken(&user)
	if errAccessToken != nil {
		return nil, errAccessToken
	}

	refreshToken, expRefreshToken, errRefreshToken := usecase.GenerateRefreshToken(&user)
	if errRefreshToken != nil {
		return nil, errRefreshToken
	}

	result := map[string]interface{}{
		"accessToken": map[string]interface{}{
			"name":       accessTokenCookieName,
			"token":      accessToken,
			"expiration": expAccessToken,
		},
		"refreshToken": map[string]interface{}{
			"name":       refreshTokenCookieName,
			"token":      refreshToken,
			"expiration": expRefreshToken,
		},
	}

	return result, nil
}

func (usecase *AuthUsecase) RefreshToken(c echo.Context) (map[string]string, error) {

	cookie, errCookie := c.Cookie(GetRefreshTokenName())

	if errCookie != nil {
		return nil, errCookie
	}

	token, errParse := jwt.ParseWithClaims(cookie.Value, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(GetRefreshJWTSecret()), nil
	})

	if errParse != nil {
		fmt.Print("Error Parse")
		return nil, errParse
	}

	claims, ok := token.Claims.(*AuthClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	rctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(rctx, usecase.timeout)
	defer cancel()

	user, err := usecase.repo.GetUser(ctx, claims.User.Username)

	if err != nil {
		return nil, err
	}

	accessToken, expAccessToken, errAccessToken := usecase.GenerateAccessToken(&user)
	if errAccessToken != nil {
		return nil, errAccessToken
	}

	usecase.SetCookie(c, accessTokenCookieName, accessToken, expAccessToken)

	refreshToken, expRefreshToken, errRefreshToken := usecase.GenerateRefreshToken(&user)
	if errRefreshToken != nil {
		return nil, errRefreshToken
	}

	usecase.SetCookie(c, refreshTokenCookieName, refreshToken, expRefreshToken)

	return nil, nil
}

func (usecase *AuthUsecase) SetCookie(c echo.Context, name, token string, expiration time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func (usecase *AuthUsecase) GenerateAccessToken(user *entities.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

func (usecase *AuthUsecase) GenerateRefreshToken(user *entities.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * 7 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func GetAccessTokenName() string {
	return accessTokenCookieName
}

func GetRefreshTokenName() string {
	return refreshTokenCookieName
}

func GetJWTSecret() string {
	return jwtSecretKey
}

func GetRefreshJWTSecret() string {
	return jwtRefreshSecretKey
}

func generateToken(user *entities.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	claims := AuthClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}
