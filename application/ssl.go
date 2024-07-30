package application

import (
	"autossl/domain/model"
	"autossl/domain/service"
	"autossl/infrastructure/acme"
	"autossl/infrastructure/exception"
	"autossl/infrastructure/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"path/filepath"
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
	code := util.GenerateID()
	err := service.CreateCert(certCommand.Domain, code)
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
	code := c.Param("code")
	filePath := filepath.Join(acme.CertPath, code)

	etag, err := service.Etag(filePath)
	if err != nil {
		return err
	}
	c.Response().Header().Set("ETag", etag)

	return c.File(filePath)
}

func DownloadHead(c echo.Context) error {
	code := c.Param("code")
	filePath := filepath.Join(acme.CertPath, code)

	etag, err := service.Etag(filePath)
	if err != nil {
		return err
	}
	c.Response().Header().Set("ETag", etag)

	return c.NoContent(http.StatusOK)
}
