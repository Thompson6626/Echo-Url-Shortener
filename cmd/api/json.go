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

func readJSON(c echo.Context, data any) error {
	maxBytes := int64(1_048_578) // 1 MB
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxBytes)

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}
// Wrapping the response with error
func writeJSONError(c echo.Context, status int, message string) error {
    errResp := struct {
        Error struct {
            Message string `json:"message"`
            Code    int    `json:"code,omitempty"`
        } `json:"error"`
    }{}

    errResp.Error.Message = message
    errResp.Error.Code = status

    return c.JSON(status, errResp)
}

func (app *application) jsonResponse(c echo.Context, status int, data any) error {
	wrapResp := struct {
		Data any `json:"data"`
	}{}

	wrapResp.Data = data

	return c.JSON(status, wrapResp)
}

