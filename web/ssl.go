package web

import (
	"autossl/common/exception"
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
	e.DELETE("/:uuid", remove)
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

	certs := service.GetCerts()
	existingCerts := make(map[string]struct{})
	for _, cert := range certs {
		existingCerts[cert.Name] = struct{}{}
	}
	if _, exists := existingCerts[certCommand.Domain]; exists {
		err := exception.CertificateExistsErr(certCommand.Domain)
		_ = c.JSON(http.StatusBadRequest, response.Message(err.Error()))
		return nil
	}

	id := util.GenerateID()

	err := os.MkdirAll(service.CertPath, 0755)
	if err != nil {
		fmt.Println("错误:", err)
	}

	err = middleware.Issue(certCommand.Domain)
	if err != nil {
		return err
	}

	err = middleware.Install(certCommand.Domain, id)
	if err != nil {
		return err
	}
	service.SaveUuid(certCommand.Domain, id)
	if err != nil {
		return err
	}

	_ = c.JSON(http.StatusOK, response.Message(certCommand.Domain+" Certificate generated successfully."))
	return nil
}

func remove(c echo.Context) error {
	uuid := c.Param("uuid")
	err := service.RemoveUUID(uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Message(err.Error()))
	}
	err = service.RemoveFiles(uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Message(err.Error()))
	}
	return c.JSON(http.StatusOK, response.Message("Certificate removed successfully."))
}
