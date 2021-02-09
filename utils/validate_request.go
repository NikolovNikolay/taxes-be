package util

import (
	"errors"
	"github.com/labstack/echo"
	"net/http"
)

func ValidateRequest(c echo.Context, data interface{}) error {
	if err := c.Bind(data); err != nil {
		return err
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	return nil
}

func ValidateRequestBody(c echo.Context, data interface{}) error {
	if c.Request().Body == http.NoBody {
		return errors.New("no body")
	}

	return ValidateRequest(c, data)
}
