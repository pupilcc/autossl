package application

import (
	"autossl/domain/model"
	"autossl/domain/service"
	"autossl/infrastructure/exception"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Generate(c echo.Context) error {
	var certCommand *model.CertCommand
	if err := c.Bind(&certCommand); err != nil {
		return err
	}

	// Find the certificate by domain
	exist := service.ExistCert(certCommand.Domain)
	if exist {
		err := exception.CertificateExistsErr(certCommand.Domain)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create ssl
	err := service.CreateCert(certCommand.Domain)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func ListCert(c echo.Context) error {
	list, err := service.ListCert()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, list)
}

func DeleteCert(c echo.Context) error {
	code := c.Param("code")

	err := service.DeleteCert(code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func RotateDownloadCode(c echo.Context) error {
	fileType := c.Param("fileType")
	if fileType != "crt" && fileType != "key" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid file type")
	}

	if err := service.RotateDownloadCode(c.Param("code"), fileType); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func Upload(c echo.Context) error {
	certFile, err := c.FormFile("cert")
	if err != nil {
		return err
	}

	keyFile, err := c.FormFile("key")
	if err != nil {
		return err
	}

	domainName := c.FormValue("domain")

	err = service.ImportCert(domainName, certFile, keyFile)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func Download(c echo.Context) error {
	return download(c, false)
}

func DownloadHead(c echo.Context) error {
	return download(c, true)
}

func download(c echo.Context, head bool) error {
	filePath, private, err := service.DownloadFile(c.Param("file"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if private {
		c.Response().Header().Set(echo.HeaderCacheControl, "no-store")
	}
	etag, err := service.Etag(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	c.Response().Header().Set("ETag", etag)

	if head {
		return c.NoContent(http.StatusOK)
	}
	return c.File(filePath)
}
