package web

import (
	"autossl/internal/domain"
	"autossl/internal/service"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"path/filepath"
)

func SSLRoutes(e *echo.Echo) {
	e.POST("/import", upload)
	e.GET("/dl/:uuid", download)
	e.HEAD("/dl/:uuid", downloadHead)
	e.GET("/list", list)
}

func upload(c echo.Context) error {
	certFile, err := c.FormFile("cert")
	if err != nil {
		return err
	}

	keyFile, err := c.FormFile("key")
	if err != nil {
		return err
	}

	certName := c.FormValue("name")

	err = service.AddCert(certName, certFile, keyFile)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.String(http.StatusOK, fmt.Sprintf("Files %s and %s uploaded successfully.", certFile.Filename, keyFile.Filename))
}

func download(c echo.Context) error {
	uuid := c.Param("uuid")
	filePath := filepath.Join(service.CertPath, uuid)

	etag, err := service.Etag(filePath)
	if err != nil {
		return err
	}
	c.Response().Header().Set("ETag", etag)

	return c.File(filePath)
}

func downloadHead(c echo.Context) error {
	uuid := c.Param("uuid")
	filePath := filepath.Join(service.CertPath, uuid)

	etag, err := service.Etag(filePath)
	if err != nil {
		return err
	}
	c.Response().Header().Set("ETag", etag)

	return c.NoContent(http.StatusOK)
}

func list(c echo.Context) error {
	certs := service.GetCerts()
	dm := os.Getenv("DOMAIN")
	url := dm + "/dl/"
	list := make([]domain.CertDTO, 0, len(certs))
	for _, cert := range certs {
		certLink := url + cert.Id + ".crt"
		keyLink := url + cert.Id + ".key"
		certDTO := domain.CertDTO{Name: cert.Name, Id: cert.Id, Cert: certLink, Key: keyLink}
		list = append(list, certDTO)
	}
	_ = c.JSON(http.StatusOK, list)
	return nil
}
