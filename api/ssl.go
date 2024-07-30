package api

import (
	"autossl/application"
	"github.com/labstack/echo/v4"
)

func SSLRoutes(e *echo.Echo) {
	e.POST("/import", upload)
	e.GET("/dl/:code", download)
	e.HEAD("/dl/:code", downloadHead)
	e.GET("/list", list)
	e.POST("/generate", generate)
	e.DELETE("/:code", remove)
}

func upload(c echo.Context) error {
	return application.Upload(c)
}

func download(c echo.Context) error {
	return application.Download(c)
}

func downloadHead(c echo.Context) error {
	return application.DownloadHead(c)
}

func list(c echo.Context) error {
	return application.ListCert(c)
}

func generate(c echo.Context) error {
	return application.Generate(c)
}

func remove(c echo.Context) error {
	return application.DeleteCert(c)
}
