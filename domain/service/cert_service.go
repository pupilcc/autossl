package service

import (
	"autossl/domain/model"
	"autossl/infrastructure/acme"
	"autossl/infrastructure/exception"
	"autossl/infrastructure/repository"
	"autossl/infrastructure/util"
	"crypto/md5"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ponytail: process-local locks only; use a distributed lock for multiple replicas.
var createCertLocks sync.Map

func ExistCert(domain string) bool {
	domain = baseDomain(domain)
	repo := repository.GetCertRepo()

	byDomain, err := repo.FindByDomain(domain)
	if err != nil {
		return false
	}
	return byDomain != nil
}

func CreateCert(domain string) error {
	domain = baseDomain(domain)
	lock, _ := createCertLocks.LoadOrStore(domain, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	if ExistCert(domain) {
		return exception.CertificateExistsErr(domain)
	}

	code := util.GenerateID()
	keyCode := util.GenerateID()
	for keyCode == code {
		keyCode = util.GenerateID()
	}

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

	repo := repository.GetCertRepo()
	cert := &model.Cert{
		Code:     code,
		CertCode: code,
		KeyCode:  keyCode,
		Domain:   domain,
	}
	err = repo.Create(cert)
	if err != nil {
		return err
	}

	return nil
}

func baseDomain(domain string) string {
	domain = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(domain)), ".")
	return strings.TrimPrefix(domain, "*.")
}

func ImportCert(domainName string, certFile *multipart.FileHeader, keyFile *multipart.FileHeader) error {
	// checked cert exist
	if ExistCert(domainName) {
		return exception.CertificateExistsErr(domainName)
	}

	err := CreateCert(domainName)
	if err != nil {
		return err
	}

	cert, err := repository.GetCertRepo().FindByDomain(baseDomain(domainName))
	if err != nil {
		return err
	}

	err = uploadFile(cert.Code, certFile, keyFile)
	if err != nil {
		return err
	}

	return nil
}

func ListCert() ([]*model.Cert, error) {
	repo := repository.GetCertRepo()
	list, err := repo.List()
	if err != nil {
		return nil, err
	}

	dm := strings.TrimRight(os.Getenv("DOMAIN"), "/")
	url := dm + "/dl/"
	for _, cert := range list {
		cert.Domain = baseDomain(cert.Domain)
		cert.DNSNames = certificateDNSNames(filepath.Join(acme.CertPath, cert.Code+".crt"))
		certCode := cert.CertCode
		if certCode == "" {
			certCode = cert.Code
		}
		keyCode := cert.KeyCode
		if keyCode == "" {
			keyCode = cert.Code
		}
		certLink := url + certCode + ".crt"
		keyLink := url + keyCode + ".key"
		cert.Cert = certLink
		cert.Key = keyLink
	}

	return list, nil
}

func DownloadFile(file string) (string, bool, error) {
	code, extension, private, ok := downloadTarget(file)
	if !ok {
		return "", false, os.ErrNotExist
	}

	repo := repository.GetCertRepo()
	var cert *model.Cert
	var err error
	if private {
		cert, err = repo.FindByKeyCode(code)
	} else {
		cert, err = repo.FindByCertCode(code)
	}
	if err != nil {
		return "", private, err
	}
	return filepath.Join(acme.CertPath, cert.Code+extension), private, nil
}

func downloadTarget(file string) (code string, extension string, private bool, ok bool) {
	if code, ok = strings.CutSuffix(file, ".crt"); ok && code != "" {
		return code, ".crt", false, true
	}
	if code, ok = strings.CutSuffix(file, ".key"); ok && code != "" {
		return code, ".key", true, true
	}
	return "", "", false, false
}

func RotateDownloadCode(code string, fileType string) error {
	repo := repository.GetCertRepo()
	cert, err := repo.FindByCode(code)
	if err != nil {
		return err
	}

	certCode := cert.CertCode
	if certCode == "" {
		certCode = cert.Code
	}
	keyCode := cert.KeyCode
	if keyCode == "" {
		keyCode = cert.Code
	}

	newCode := util.GenerateID()
	for newCode == certCode || newCode == keyCode {
		newCode = util.GenerateID()
	}

	if fileType == "crt" {
		return repo.UpdateCertCode(code, newCode)
	}
	return repo.UpdateKeyCode(code, newCode)
}

func certificateDNSNames(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	for {
		block, rest := pem.Decode(data)
		if block == nil {
			certificate, err := x509.ParseCertificate(data)
			if err == nil {
				return certificate.DNSNames
			}
			return nil
		}
		data = rest
		if block.Type != "CERTIFICATE" {
			continue
		}
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err == nil && len(certificate.DNSNames) > 0 {
			return certificate.DNSNames
		}
	}
}

func DeleteCert(code string) error {
	var err error
	repo := repository.GetCertRepo()

	byCode, err := repo.FindByCode(code)
	if err != nil {
		return err
	}

	err = repo.Delete(code)
	if err != nil {
		return err
	}

	err = acme.Remove(byCode.Domain)
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
