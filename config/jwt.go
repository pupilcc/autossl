package config

import "github.com/labstack/echo/v4"

var skipTokenPaths = []string{
	"/",
	"/health",
	"/login",
	"/dl/:code",
}

func Skip(c echo.Context) bool {
	path := c.Path()
	for _, p := range skipTokenPaths {
		if path == p {
			return true
		}
	}

	return false
}
