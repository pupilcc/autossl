package service

import (
	"autossl/domain/model"
	"autossl/infrastructure/acme"
	"autossl/infrastructure/database"
	"autossl/infrastructure/exception"
	"autossl/infrastructure/repository"
	"autossl/infrastructure/util"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func ExistCert(domain string) bool {
	repo := &repository.CertRepo{
		Db: database.Init(),
	}

	byDomain, err := repo.FindByDomain(domain)
	if err != nil {
		return false
	}
	return byDomain != nil
}

func CreateCert(domain string, code string) error {
	err := os.MkdirAll(acme.CertPath, 0755)
	if err != nil {
		return err
	}

	err = acme.Issue(domain)
	if err != nil {
		return err
	}

	err = acme.Install(domain, code)
	if err != nil {
		return err
	}

	repo := &repository.CertRepo{
		Db: database.Init(),
	}
	cert := &model.Cert{
		Code:   code,
		Domain: domain,
	}
	err = repo.Create(cert)
	if err != nil {
		return err
	}

	return nil
}

func ImportCert(domainName string, certFile *multipart.FileHeader, keyFile *multipart.FileHeader) error {
	// checked cert exist
	if ExistCert(domainName) {
		return exception.CertificateExistsErr(domainName)
	}

	code := util.GenerateID()
	err := CreateCert(domainName, code)
	if err != nil {
		return err
	}

	err = uploadFile(code, certFile, keyFile)
	if err != nil {
		return err
	}

	return nil
}

func ListCert() ([]*model.Cert, error) {
	repo := &repository.CertRepo{
		Db: database.Init(),
	}
	list, err := repo.List()
	if err != nil {
		return nil, err
	}

	dm := os.Getenv("DOMAIN")
	url := dm + "/dl/"
	for _, cert := range list {
		certLink := url + cert.Code + ".crt"
		keyLink := url + cert.Code + ".key"
		cert.Cert = certLink
		cert.Key = keyLink
	}

	return list, nil
}

func DeleteCert(code string) error {
	repo := &repository.CertRepo{
		Db: database.Init(),
	}
	err := repo.Delete(code)
	if err != nil {
		return err
	}

	err = deleteFiles(code)
	if err != nil {
		return err
	}

	return nil
}

func Etag(filePath string) (string, error) {
	etag := ""
	file, err := os.Open(filePath)
	if err != nil {
		return etag, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return etag, err
	}
	etag = fmt.Sprintf("%x", hash.Sum(nil))
	return etag, nil
}

func deleteFiles(code string) error {
	certPath := filepath.Join(acme.CertPath, code+".crt")
	keyPath := filepath.Join(acme.CertPath, code+".key")

	filePaths := []string{certPath, keyPath}

	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func uploadFile(code string, certFile *multipart.FileHeader, keyFile *multipart.FileHeader) error {
	certSrc, err := certFile.Open()
	if err != nil {
		return err
	}
	defer certSrc.Close()

	keySrc, err := keyFile.Open()
	if err != nil {
		return err
	}
	defer keySrc.Close()

	_ = os.MkdirAll(acme.CertPath, 0755)
	certDst, err := os.Create(filepath.Join(acme.CertPath, code+".crt"))
	if err != nil {
		return err
	}
	defer certDst.Close()

	keyDst, err := os.Create(filepath.Join(acme.CertPath, code+".key"))
	if err != nil {
		return err
	}
	defer keyDst.Close()

	if _, err = io.Copy(certDst, certSrc); err != nil {
		return err
	}

	if _, err = io.Copy(keyDst, keySrc); err != nil {
		return err
	}
	return nil
}
