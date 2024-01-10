package service

import (
	"autossl/common/exception"
	"autossl/common/util"
	"autossl/internal/domain"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var CertPath = "./data/cert"
var UuidPath = "./data/uuid.json"

func AddCert(name string, certFile *multipart.FileHeader, keyFile *multipart.FileHeader) error {
	// checked cert exist
	certs := GetCerts()
	if certs != nil {
		for _, cert := range certs {
			if cert.Name == name {
				return exception.CertificateExistsErr(name)
			}
		}
	}

	id := util.GenerateID()

	SaveUuid(name, id)

	err := uploadFile(id, certFile, keyFile)
	if err != nil {
		return err
	}

	return nil
}

func uploadFile(id string, certFile *multipart.FileHeader, keyFile *multipart.FileHeader) error {
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

	_ = os.MkdirAll(CertPath, 0755)
	certDst, err := os.Create(filepath.Join(CertPath, id+".crt"))
	if err != nil {
		return err
	}
	defer certDst.Close()

	keyDst, err := os.Create(filepath.Join(CertPath, id+".key"))
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

func SaveUuid(name string, id string) {
	certs := GetCerts()
	if certs == nil {
		certs = []domain.Cert{{Name: name, Id: id}}
	} else {
		cert := domain.Cert{Name: name, Id: id}
		certs = append(certs, cert)
	}

	err := writeFile(certs)
	if err != nil {
		fmt.Println("Write file fail: ", err)
	}
}

func GetCerts() []domain.Cert {
	_, err := os.Stat(UuidPath)
	if os.IsNotExist(err) {
		return nil
	}

	_ = fmt.Sprintf("File exists")
	jsonDataFromFile, err := os.ReadFile(UuidPath)
	if err != nil {
		fmt.Println("Read file fail: ", err)
	}

	certs := make([]domain.Cert, 1)
	err = json.Unmarshal(jsonDataFromFile, &certs)
	if err != nil {
		fmt.Println("parse json fail: ", err)
	}

	return certs
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

func RemoveUUID(id string) error {
	certs := GetCerts()
	if certs == nil {
		return nil
	}
	for i, cert := range certs {
		if cert.Id == id {
			certs[i] = certs[len(certs)-1]
			certs[len(certs)-1] = domain.Cert{}
			certs = certs[:len(certs)-1]
			break
		}
	}
	err := writeFile(certs)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFiles(id string) error {
	certPath := filepath.Join(CertPath, id+".crt")
	keyPath := filepath.Join(CertPath, id+".key")

	filePaths := []string{certPath, keyPath}

	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(certs []domain.Cert) error {
	jsonData, err := json.Marshal(certs)
	if err != nil {
		return err
	}

	err = os.WriteFile(UuidPath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
