package middleware

import "github.com/labstack/echo/v4"

var skipTokenPaths = map[string]bool{
	"/":             true,
	"/health":       true,
	"/login":        true,
	"/ssl/dl/:uuid": true,
}

func Skip(c echo.Context) bool {
	path := c.Path()
	if _, ok := skipTokenPaths[path]; ok {
		return true
	}

	return false
}
