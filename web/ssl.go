package web

import (
	"autossl/internal/domain"
	"autossl/internal/service"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func SSLRoutes(e *echo.Echo) {
	r := e.Group("/ssl")
	r.POST("/import", upload)
	r.GET("/download/:uuid", download)
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
	return c.Attachment(service.CertPath+"/"+uuid, uuid)
}

func list(c echo.Context) error {
	certs := service.GetCerts()
	url := os.Getenv("DOMAIN")
	list := make([]domain.CertDTO, 0, len(certs))
	for _, cert := range certs {
		certLink := url + "/download/" + cert.Id + ".crt"
		keyLink := url + "/download/" + cert.Id + ".key"
		certDTO := domain.CertDTO{Name: cert.Name, Id: cert.Id, Cert: certLink, Key: keyLink}
		list = append(list, certDTO)
	}
	_ = c.JSON(http.StatusOK, list)
	return nil
}
