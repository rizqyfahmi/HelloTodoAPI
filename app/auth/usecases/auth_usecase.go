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

func (usecase *AuthUsecase) Login(c echo.Context, username, password string) error {

	rctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(rctx, usecase.timeout)
	defer cancel()

	user, err := usecase.repo.GetUser(ctx, username)

	if err != nil {
		return err
	}

	errBcrypt := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if errBcrypt != nil {
		return errBcrypt
	}

	generateAndSetToken(c, user)

	return nil
}

func (usecase *AuthUsecase) RefreshToken(c echo.Context) error {

	rctx := c.Request().Context()
	ctx, cancel := context.WithTimeout(rctx, usecase.timeout)
	defer cancel()

	cookie, errCookie := c.Cookie(GetRefreshTokenName())

	if errCookie != nil {
		return errCookie
	}

	token, errParse := jwt.ParseWithClaims(cookie.Value, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(GetRefreshJWTSecret()), nil
	})

	if errParse != nil {
		fmt.Print("Error Parse")
		return errParse
	}

	claims, ok := token.Claims.(*AuthClaims)

	if !ok || !token.Valid {
		return fmt.Errorf("invalid token")
	}

	user, err := usecase.repo.GetUser(ctx, claims.User.Username)

	if err != nil {
		return err
	}

	generateAndSetToken(c, user)

	return nil
}

func setCookie(c echo.Context, name, token string, expiration time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func generateAndSetToken(c echo.Context, user entities.User) error {
	accessToken, expAccessToken, errAccessToken := generateAccessToken(&user)
	if errAccessToken != nil {
		return errAccessToken
	}

	setCookie(c, accessTokenCookieName, accessToken, expAccessToken)

	refreshToken, expRefreshToken, errRefreshToken := generateRefreshToken(&user)
	if errRefreshToken != nil {
		return errRefreshToken
	}

	setCookie(c, refreshTokenCookieName, refreshToken, expRefreshToken)

	return nil
}

func generateAccessToken(user *entities.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

func generateRefreshToken(user *entities.User) (string, time.Time, error) {
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
