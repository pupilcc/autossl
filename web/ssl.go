package web

import (
	"autossl/internal/domain"
	"autossl/internal/service"
	"crypto/md5"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func SSLRoutes(e *echo.Echo) {
	r := e.Group("/ssl")
	r.POST("/import", upload)
	r.GET("/dl/:uuid", download)
	r.GET("", list)
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
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	etag := fmt.Sprintf("%x", hash.Sum(nil))
	c.Response().Header().Set("ETag", etag)

	return c.File(filePath)
}

func list(c echo.Context) error {
	certs := service.GetCerts()
	dm := os.Getenv("DOMAIN")
	url := dm + "/ssl/dl/"
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
