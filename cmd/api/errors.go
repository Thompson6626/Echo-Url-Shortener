package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *application) internalServerError(c echo.Context, err error) error {
	app.logger.Errorw("internal error", "method", c.Request().Method, "path", c.Path(), "error", err.Error())
	return writeJSONError(c, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) forbiddenResponse(c echo.Context) error {
	app.logger.Warnw("forbidden", "method", c.Request().Method, "path", c.Path())
	return writeJSONError(c, http.StatusForbidden, "forbidden")
}

func (app *application) badRequestResponse(c echo.Context, err error) error {
	app.logger.Warnw("bad request", "method", c.Request().Method, "path", c.Path(), "error", err.Error())
	return writeJSONError(c, http.StatusBadRequest, err.Error())
}

func (app *application) conflictResponse(c echo.Context, err error) error {
	app.logger.Errorw("conflict response", "method", c.Request().Method, "path", c.Path(), "error", err.Error())
	return writeJSONError(c, http.StatusConflict, err.Error())
}

func (app *application) notFoundResponse(c echo.Context, err error) error {
	app.logger.Warnw("not found error", "method", c.Request().Method, "path", c.Path(), "error", err.Error())
	return writeJSONError(c, http.StatusNotFound, "not found")
}

func (app *application) unauthorizedErrorResponse(c echo.Context, err error) error {
	app.logger.Warnw("unauthorized error", "method", c.Request().Method, "path", c.Path(), "error", err.Error())
	return writeJSONError(c, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedBasicErrorResponse(c echo.Context, err error) error {
	app.logger.Warnw("unauthorized basic error", "method", c.Request().Method, "path", c.Path(), "error", err.Error())
	c.Response().Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	return writeJSONError(c, http.StatusUnauthorized, "unauthorized")
}

func (app *application) rateLimitExceededResponse(c echo.Context, retryAfter string) error {
	app.logger.Warnw("rate limit exceeded", "method", c.Request().Method, "path", c.Path())
	c.Response().Header().Set("Retry-After", retryAfter)
	return writeJSONError(c, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
	