package main

import (
	"Url-Shortener/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) registerUserHandle(c echo.Context) error {

	payload, err := BindAndValidate[RegisterUserPayload](c)

	if err != nil {
		return app.badRequestResponse(c, err)
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		return app.internalServerError(c, err)
	}

	ctx := c.Request().Context()

	err = app.store.Users.Create(ctx, user)

	if err != nil {
		switch err {
		case store.ErrDuplicateEmail, store.ErrDuplicateUsername:
			return app.badRequestResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	return c.JSON(http.StatusCreated, user)
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) loginUserHandler(c echo.Context) error {

	payload, err := BindAndValidate[LoginUserPayload](c)

	if err != nil {
		return app.badRequestResponse(c, err)
	}

	user, err := app.store.Users.GetByEmail(c.Request().Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			return app.unauthorizedErrorResponse(c, err)
		default:
			return app.internalServerError(c, err)
		}
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		return app.unauthorizedErrorResponse(c, err)
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		return app.internalServerError(c, err)
	}

	if err := app.jsonResponse(c, http.StatusCreated, token); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}
