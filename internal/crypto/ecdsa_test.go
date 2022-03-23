package crypto

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

const privateKey = `-----BEGIN PRIVATE KEY-----
MHcCAQEEIAwRtGkYqi732qh84HafnKE7YkW0CNpvvNseNGbxpsgGoAoGCCqGSM49
AwEHoUQDQgAE+xszAkYoKJP5CEvCaLuCGxAGDKIWecgPQxYElRWn/183SnpMHDRE
fXa4/Jzadq8dmo4taNQucoOLjD7IaN5OcA==
-----END PRIVATE KEY-----`

func TestReadPrivateKey(t *testing.T) {
	_, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Errorf("error parsing PEM ECDSA private key: %v", err)
		t.FailNow()
	}
}

func TestNewJWTECDSA(t *testing.T) {
	j, err := NewJWTECDSA(privateKey, jwt.SigningMethodES256)
	if err != nil {
		t.Errorf("error creating NewJWTECDSA: %v", err)
		t.FailNow()
	}
	if j == nil {
		t.Errorf("jwt cannot be nil")
		t.FailNow()
	}
	if j.k == nil {
		t.Errorf("incorrect key, cannot be nil")
		t.FailNow()
	}
	if j.method != jwt.SigningMethodES256 {
		t.Errorf("incorrect method, got %v, want %v", j.method, jwt.SigningMethodES256)
		t.FailNow()
	}
}

func TestNewJWTES256(t *testing.T) {
	j, err := NewJWTES256()
	if err != nil {
		t.Errorf("error creating NewJWTES256: %v", err)
		t.FailNow()
	}
	if j == nil {
		t.Errorf("jwt cannot be nil")
		t.FailNow()
	}
	if j.k == nil {
		t.Errorf("incorrect key, cannot be nil")
		t.FailNow()
	}
	if j.method != jwt.SigningMethodES256 {
		t.Errorf("incorrect method, got %v, want %v", j.method, jwt.SigningMethodES256)
		t.FailNow()
	}

	x509pvk, err := x509.MarshalECPrivateKey(j.k)
	if err != nil {
		t.Errorf("error marshaling ECDSA private key: %v", err)
		t.FailNow()
	}

	var pvkbuf bytes.Buffer
	if err := pem.Encode(&pvkbuf, &pem.Block{Type: "PRIVATE KEY", Bytes: x509pvk}); err != nil {
		t.Errorf("error encoding ECDSA private key: %v", err)
		t.FailNow()
	}
	if pvkbuf.Len() == 0 {
		t.Errorf("error encoding ECDSA private key: buffer cannot be empty")
		t.FailNow()
	}

	pk, err := jwt.ParseECPrivateKeyFromPEM(pvkbuf.Bytes())
	if err != nil {
		t.Errorf("error parsing PEM ECDSA private key: %v", err)
		t.FailNow()
	}

	if !j.k.Equal(pk) {
		t.Errorf("incorrect private key, got %v, want %v", pk, j.k)
		t.FailNow()
	}

	fmt.Print(string(pvkbuf.Bytes()))

}
