package repository

import (
	"autossl/domain/model"
)

type CertRepository interface {
	Create(cert *model.Cert) error
	FindByDomain(domain string) (*model.Cert, error)
	FindByCode(code string) (*model.Cert, error)
	FindByCertCode(code string) (*model.Cert, error)
	FindByKeyCode(code string) (*model.Cert, error)
	List() ([]*model.Cert, error)
	UpdateCertCode(code string, certCode string) error
	UpdateKeyCode(code string, keyCode string) error
	Delete(code string) error
}
