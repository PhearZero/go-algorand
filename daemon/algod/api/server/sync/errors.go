package sync

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func notfound(c echo.Context, msg string) error {
	// Write error to response
	err := c.JSON(http.StatusNotFound, DatabaseError{
		Error:  "not_found",
		Reason: msg,
	})
	// Handle writer errors
	if err != nil {
		return err
	}

	// Send caller a generic error
	return errors.New(msg)
}
func badargument(c echo.Context, msg string) error {
	err := c.JSON(http.StatusBadRequest, DatabaseError{
		Error:  "unknown_error",
		Reason: msg,
	})

	if err != nil {
		return err
	}

	return errors.New(msg)
}
