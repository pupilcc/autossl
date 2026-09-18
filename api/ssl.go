package api

import (
	"autossl/application"
	"crypto/subtle"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

func SSLRoutes(e *echo.Echo) error {
	downloads := e.Group("/dl")
	if strings.EqualFold(os.Getenv("DOWNLOAD_AUTH_ENABLED"), "true") {
		token := os.Getenv("DOWNLOAD_AUTH_TOKEN")
		if token == "" {
			return errors.New("DOWNLOAD_AUTH_TOKEN is required when DOWNLOAD_AUTH_ENABLED=true")
		}
		downloads.Use(downloadAuth(token))
	}

	e.POST("/import", upload)
	downloads.GET("/:file", download)
	downloads.HEAD("/:file", downloadHead)
	e.GET("/list", list)
	e.POST("/generate", generate)
	e.POST("/:code/rotate/:fileType", rotateDownloadCode)
	e.DELETE("/:code", remove)
	return nil
}

func downloadAuth(expected string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := strings.CutPrefix(c.Request().Header.Get(echo.HeaderAuthorization), "Bearer ")
			if !ok || subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
				c.Response().Header().Set(echo.HeaderWWWAuthenticate, "Bearer")
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid download token")
			}
			return next(c)
		}
	}
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

func rotateDownloadCode(c echo.Context) error {
	return application.RotateDownloadCode(c)
}
