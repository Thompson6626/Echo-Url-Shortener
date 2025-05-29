package main

import (
	"Url-Shortener/internal/store"
	"github.com/labstack/echo/v4"
)

func getUserFromContext(c echo.Context) *store.User {
	user, _ := c.Get("user").(*store.User)
	return user
}
