package echoaddons

import (
	"encoding/base64"
	"github.com/labstack/echo"
	"net/http"
)

func AuthHandler(secret string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-API-KEY")
			bytes, err := base64.StdEncoding.DecodeString(key)
			if err != nil || string(bytes) != secret {
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}

			return next(c)
		}
	}
}
