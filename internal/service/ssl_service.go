package service

import (
	"autossl/internal/domain"
	"autossl/util"
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
				return fmt.Errorf("cert %s already exist", name)
			}
		}
	}

	id := util.GenerateID()

	saveUuid(id, name)

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

	err = os.MkdirAll(CertPath, 0755)
	if err != nil {
		fmt.Println("错误:", err)
	}
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

	// 将源文件内容复制到目标文件
	if _, err = io.Copy(certDst, certSrc); err != nil {
		return err
	}

	if _, err = io.Copy(keyDst, keySrc); err != nil {
		return err
	}
	return nil
}

func saveUuid(id string, name string) {
	certs := GetCerts()
	if certs == nil {
		certs = []domain.Cert{{Name: name, Id: id}}
	} else {
		cert := domain.Cert{Name: name, Id: id}
		certs = append(certs, cert)
	}

	jsonData, err := json.Marshal(certs)
	if err != nil {
		fmt.Println("转换为 JSON 失败：", err)
	}

	err = os.WriteFile(UuidPath, jsonData, 0644)
	if err != nil {
		fmt.Println("写入文件失败：", err)
	}
}

func GetCerts() []domain.Cert {
	// 判断文件是否存在
	_, err := os.Stat(UuidPath)
	if os.IsNotExist(err) {
		return nil
	}

	// 存在则读取 map
	_ = fmt.Sprintf("File exists")
	// 从文件中读取 JSON
	jsonDataFromFile, err := os.ReadFile(UuidPath)
	if err != nil {
		fmt.Println("读取文件失败：", err)
	}

	// 将 JSON 解析回 map
	certs := make([]domain.Cert, 1)
	err = json.Unmarshal(jsonDataFromFile, &certs)
	if err != nil {
		fmt.Println("解析 JSON 失败：", err)
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
