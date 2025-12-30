package middleware

import (
	"net/http"

	"github.com/bencleary/uploader/internal/encryption"
	"github.com/labstack/echo/v4"
)

func ValidateEncryptionKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Request().Header.Get("key")

		if key == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Encryption key is required")
		}

		if !encryption.IsValidKey(key) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Encryption key is invalid")
		}

		return next(c)
	}
}
