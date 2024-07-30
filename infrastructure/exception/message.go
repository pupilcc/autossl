package exception

import "fmt"

const (
	CertificateExists = "certificate for %s already exists"
)

func CertificateExistsErr(domain string) error {
	return fmt.Errorf(CertificateExists, domain)
}
