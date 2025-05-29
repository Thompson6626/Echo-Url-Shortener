package main

import (
	"Url-Shortener/internal/store"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"strconv"
	"strings"
)

func (app *application) AuthTokenMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return app.unauthorizedErrorResponse(c, fmt.Errorf("authorization header is missing"))
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return app.unauthorizedErrorResponse(c, fmt.Errorf("authorization header is malformed"))
			}

			token := parts[1]
			jwtToken, err := app.authenticator.ValidateToken(token)
			if err != nil {
				return app.unauthorizedErrorResponse(c, err)
			}

			claims, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				return app.unauthorizedErrorResponse(c, fmt.Errorf("invalid token claims"))
			}

			userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
			if err != nil {
				return app.unauthorizedErrorResponse(c, err)
			}

			user, err := app.getUser(c.Request().Context(), userID)
			if err != nil {
				return app.unauthorizedErrorResponse(c, err)
			}

			c.Set("user", user)

			return next(c)
		}
	}

}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	// No Redis
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetById(ctx, userID)
	}

	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// If no user found in cache get user then save it in cache
	if user == nil {
		user, err = app.store.Users.GetById(ctx, userID)
		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (app *application) RateLimiterMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if app.config.rateLimiter.Enabled {
			remoteIP := c.RealIP()
			if allow, retryAfter := app.rateLimiter.Allow(remoteIP); !allow {
				return app.rateLimitExceededResponse(c, retryAfter.String());
			}
		}

		return next(c)
	}
}
