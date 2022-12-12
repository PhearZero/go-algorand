package sync

import (
	"github.com/labstack/echo/v4"
	"net/url"
	"strings"
)

// Unescape removes URL encoding from the path
func Unescape(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Raw URL
		rawPath := c.Request().URL.RawPath

		//println("Current url: " + rawPath)
		// Check for possible escaped characters
		if strings.Contains(rawPath, "%2F") {
			fixedPath, _ := url.QueryUnescape(rawPath)
			c.Request().URL.RawPath = strings.TrimSuffix(fixedPath, "/")
			//println("New url: " + c.Request().URL.String())
		}

		return next(c)
	}
}
