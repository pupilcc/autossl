package web

import (
	"autossl/common/response"
	"autossl/common/util"
	"autossl/internal/domain"
	"autossl/internal/service"
	"autossl/middleware"
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
	e.POST("/generate", generate)
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
		return c.JSON(http.StatusBadRequest, response.Message(err.Error()))
	}

	msg := fmt.Sprintf("Files %s and %s uploaded successfully.", certFile.Filename, keyFile.Filename)
	return c.JSON(http.StatusOK, response.Message(msg))
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

func generate(c echo.Context) error {
	var certCommand *domain.CertCommand
	if err := c.Bind(&certCommand); err != nil {
		return err
	}

	id := util.GenerateID()

	err := os.MkdirAll(service.CertPath, 0755)
	if err != nil {
		fmt.Println("错误:", err)
	}

	middleware.Issue(certCommand.Domain)
	middleware.Install(certCommand.Domain, id)
	service.SaveUuid(certCommand.Domain, id)

	return nil
}
