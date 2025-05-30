package main

import (
	"Url-Shortener/internal/store"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateUrlPayload struct {
	OriginalUrl string `json:"url"`
}

func (app *application) createUrlHandler(c echo.Context) error {

	payload, err := BindAndValidate[CreateUrlPayload](c)

	if err != nil {
		return app.badRequestResponse(c, err)
	}

	user := getUserFromContext(c)

	url := &store.ShortURL{
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

	shortenedUrl, err := app.store.Urls.GetByShortCode(context, shortCode)

	if err != nil {
		return app.internalServerError(c, err)
	}

	return c.Redirect(http.StatusFound, shortenedUrl.OriginalURL)
}

func (app *application) getAllUrlsByUserHandler(c echo.Context) error {
	user := getUserFromContext(c)

	context := c.Request().Context()

	urls, err := app.store.Urls.GetAllUrlsByUser(context, user.ID)

	if err != nil {
		return app.internalServerError(c, err)
	}

	return app.jsonResponse(c, http.StatusOK, urls)
}

func (app *application) deleteUrlHandler(c echo.Context) error {
	shortURL := c.Get("shortURL").(*store.ShortURL) // Get from context set by middleware
	ctx := c.Request().Context()

	if err := app.store.Urls.Delete(ctx, shortURL.ShortCode); err != nil {
		return app.internalServerError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
