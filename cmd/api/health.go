package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *application) healthCheckHandler(c echo.Context) error {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(c, http.StatusOK, data); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}
