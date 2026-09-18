package repository

import (
	"autossl/domain/model"
	"autossl/infrastructure/database"
	"autossl/infrastructure/ent"
	"autossl/infrastructure/ent/cert"
	"context"
)

type CertRepo struct {
	Db *ent.Client
}

var certRepo *CertRepo

func GetCertRepo() *CertRepo {
	if certRepo == nil {
		certRepo = &CertRepo{
			Db: database.GetDBClient(),
		}
	}
	return certRepo
}

func (repo *CertRepo) Create(cert *model.Cert) error {
	_, err := repo.Db.Cert.
		Create().
		SetCode(cert.Code).
		SetCertCode(cert.CertCode).
		SetKeyCode(cert.KeyCode).
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

func (repo *CertRepo) FindByCertCode(code string) (*model.Cert, error) {
	first, err := repo.Db.Cert.Query().Where(cert.Or(
		cert.CertCodeEQ(code),
		cert.And(cert.CodeEQ(code), cert.CertCodeIsNil()),
	)).Only(context.Background())
	if err != nil {
		return nil, err
	}
	return ToModelCert(first), nil
}

func (repo *CertRepo) FindByKeyCode(code string) (*model.Cert, error) {
	first, err := repo.Db.Cert.Query().Where(cert.Or(
		cert.KeyCodeEQ(code),
		cert.And(cert.CodeEQ(code), cert.KeyCodeIsNil()),
	)).Only(context.Background())
	if err != nil {
		return nil, err
	}
	return ToModelCert(first), nil
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

func (repo *CertRepo) UpdateCertCode(code string, certCode string) error {
	_, err := repo.Db.Cert.Update().Where(cert.CodeEQ(code)).SetCertCode(certCode).Save(context.Background())
	return err
}

func (repo *CertRepo) UpdateKeyCode(code string, keyCode string) error {
	_, err := repo.Db.Cert.Update().Where(cert.CodeEQ(code)).SetKeyCode(keyCode).Save(context.Background())
	return err
}

func ToModelCert(c *ent.Cert) *model.Cert {
	return &model.Cert{
		Code:     c.Code,
		CertCode: c.CertCode,
		KeyCode:  c.KeyCode,
		Domain:   c.Domain,
	}
}
