package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSONError(c echo.Context, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return c.JSON(status, envelope{Error: message})
}

func (app *application) jsonResponse(c echo.Context, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return c.JSON(status, envelope{Data: data})
}

func readJSON(c echo.Context, data any) error {
	maxBytes := int64(1_048_578) // 1 MB
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxBytes)

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}
