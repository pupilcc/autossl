package repository

import (
	"autossl/infrastructure/ent/enttest"
	"context"
	"testing"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
)

func TestLegacyDownloadCodes(t *testing.T) {
	client := enttest.Open(t, dialect.SQLite, "file:legacy?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	ctx := context.Background()
	if _, err := client.Cert.Create().SetCode("legacyCode").SetDomain("example.com").Save(ctx); err != nil {
		t.Fatal(err)
	}

	repo := &CertRepo{Db: client}
	if _, err := repo.FindByCertCode("legacyCode"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.FindByKeyCode("legacyCode"); err != nil {
		t.Fatal(err)
	}

	if err := repo.UpdateKeyCode("legacyCode", "newKeyCode"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.FindByKeyCode("legacyCode"); err == nil {
		t.Fatal("legacy private-key URL still works after rotation")
	}
	if _, err := repo.FindByKeyCode("newKeyCode"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.FindByCertCode("legacyCode"); err != nil {
		t.Fatal(err)
	}
}
