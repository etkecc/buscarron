package controllers

import (
	"crypto/subtle"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func basicAuth(auth Auth, log *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(login, password string, c echo.Context) (bool, error) {
		allowedIP := true
		if len(auth.IPs) != 0 {
			allowedIP = slices.Contains(auth.IPs, c.RealIP())
		}
		match := equals(auth.Login, login) && equals(auth.Password, password)
		log.
			Info().
			Str("from", c.RealIP()).
			Str("path", c.Request().URL.Path).
			Bool("allowed_ip", allowedIP).
			Bool("allowed_credentials", match).
			Msg("authorization attempt")

		return match && allowedIP, nil
	})
}

// equals performs equality check in constant time
func equals(str1, str2 string) bool {
	b1 := []byte(str1)
	b2 := []byte(str2)
	return subtle.ConstantTimeEq(int32(len(b1)), int32(len(b2))) == 1 && subtle.ConstantTimeCompare(b1, b2) == 1
}
