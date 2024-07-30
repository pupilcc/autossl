package repository

import (
	"autossl/domain/model"
	"autossl/infrastructure/ent"
	"autossl/infrastructure/ent/cert"
	"context"
)

type CertRepo struct {
	Db *ent.Client
}

func (repo *CertRepo) Create(cert *model.Cert) error {
	_, err := repo.Db.Cert.
		Create().
		SetCode(cert.Code).
		SetDomain(cert.Domain).
		Save(context.Background())

	return err
}

func (repo *CertRepo) FindByDomain(domain string) (*model.Cert, error) {
	first, err := repo.Db.Cert.Query().Where(cert.DomainEQ(domain)).Only(context.Background())
	if err != nil {
		return nil, err
	}
	return ToModelCert(first), err
}

func (repo *CertRepo) FindByCode(code string) (*model.Cert, error) {
	first, err := repo.Db.Cert.Query().Where(cert.CodeEQ(code)).Only(context.Background())
	if err != nil {
		return nil, err
	}
	return ToModelCert(first), err
}

func (repo *CertRepo) List() ([]*model.Cert, error) {
	q := repo.Db.Cert.Query().
		Order(ent.Desc(cert.FieldCreatedAt))

	list, err := q.All(context.Background())
	if err != nil {
		return nil, err
	}

	l := make([]*model.Cert, len(list))

	for i, v := range list {
		l[i] = ToModelCert(v)
	}
	return l, nil
}

func (repo *CertRepo) Delete(code string) error {
	_, err := repo.Db.Cert.
		Delete().
		Where(cert.CodeEQ(code)).
		Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func ToModelCert(c *ent.Cert) *model.Cert {
	return &model.Cert{
		Code:   c.Code,
		Domain: c.Domain,
	}
}
