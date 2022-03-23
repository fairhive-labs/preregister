package crypto

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestNewJWTES256(t *testing.T) {
	j := NewJWTES256()
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

}
