package main

import "github.com/labstack/echo/v4"


// BindAndValidate binds the request body to a struct of type T and validates it using the `Validate` instance.
//
// T must be a struct type. This function is intended to reduce repetitive code in Echo handlers.
// It performs two steps:
//   1. Binds the incoming JSON request body to the struct of type T using Echo's Context.Bind method.
//   2. Validates the struct using the configured validator (e.g., github.com/go-playground/validator).
//
// If binding or validation fails, the function returns a non-nil error.
// On success, it returns a pointer to the validated struct.
//
// Example:
//
//	type ShortURLPayload struct {
//	    URL string `json:"url" validate:"required,url"`
//	}
//
//	func (app *Application) createShortURLHandler(c echo.Context) error {
//	    payload, err := BindAndValidate[ShortURLPayload](c)
//	    if err != nil {
//	        return app.badRequestResponse(c, err)
//	    }
//	    // Use payload.URL ...
//	}
func BindAndValidate[T any](c echo.Context) (*T, error) {
	var payload T
	if err := c.Bind(&payload); err != nil {
		return nil, err
	}
	if err := Validate.Struct(payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
