package exception

import "fmt"

const (
	CertificateExists = "certificate for %s already exists"
)

func CertificateExistsErr(name string) error {
	return fmt.Errorf(CertificateExists, name)
}
