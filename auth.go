package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Err struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func loginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		auth := c.Request().Header.Get("Authorization")

		if auth != "" {
			_, err := time.Parse("January 2, 2006", auth)
			if err != nil {
				c.JSON(http.StatusUnauthorized, Err{
					Message: "Authorization header not provided",
					Error:   err.Error(),
				})
			}

			return next(c)

		}

		return c.JSON(http.StatusUnauthorized, Err{
			Error:   "Unauthorized",
			Message: "Please provide valid authorization in the header",
		})
	}
}
