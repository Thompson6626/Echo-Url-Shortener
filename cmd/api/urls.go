package main

import (
	"Url-Shortener/internal/store"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateUrlPayload struct {
	ShortCode   string `json:"short_code"`
	OriginalUrl string `json:"original_url"`
}

func (app *application) createUrlHandler(c echo.Context) error {

	payload, err := BindAndValidate[CreateUrlPayload](c)

	if err != nil {
		return app.badRequestResponse(c, err)
	}

	user := getUserFromContext(c)

	url := &store.ShortURL{
		ShortCode:   payload.ShortCode,
		OriginalURL: payload.OriginalUrl,
		UserID:      user.ID,
	}

	context := c.Request().Context()

	if err := app.store.Urls.Create(context, url); err != nil {
		return app.internalServerError(c, err)
	}

	if err := app.jsonResponse(c, http.StatusCreated, url); err != nil {
		return app.internalServerError(c, err)
	}

	return nil
}

func (app *application) getUrlHandler(c echo.Context) error {
	shortCode := c.Param("shortCode")

	if shortCode == "" {
		return writeJSONError(c, http.StatusBadRequest, "shortCode is required")
	}

	context := c.Request().Context()

	originalUrl, err := app.store.Urls.GetByShortCode(context, shortCode)

	if err != nil {
		return app.internalServerError(c, err)
	}

	return c.Redirect(http.StatusFound, originalUrl)
}
