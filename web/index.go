package web

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func IndexRoutes(e *echo.Echo) {
	e.GET("/", index)
	e.HEAD("/", indexHead)
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func indexHead(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return c.NoContent(http.StatusOK)
}
