package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBaseDomain(t *testing.T) {
	for input, want := range map[string]string{
		"latticeway.com":      "latticeway.com",
		"*.latticeway.com":    "latticeway.com",
		" *.LATTICEWAY.COM. ": "latticeway.com",
	} {
		if got := baseDomain(input); got != want {
			t.Errorf("baseDomain(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestCertificateDNSNames(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		DNSNames:     []string{"example.com", "*.example.com"},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "certificate.crt")
	if err := os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0600); err != nil {
		t.Fatal(err)
	}

	if got := certificateDNSNames(path); !reflect.DeepEqual(got, template.DNSNames) {
		t.Fatalf("certificateDNSNames() = %q, want %q", got, template.DNSNames)
	}
}
