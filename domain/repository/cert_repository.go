package repository

import (
	"autossl/domain/model"
)

type CertRepository interface {
	Create(cert *model.Cert) error
	FindByDomain(domain string) (*model.Cert, error)
	List() ([]*model.Cert, error)
	Delete(code string) error
}
