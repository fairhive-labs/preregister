package crypto

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	privateKey = `-----BEGIN PRIVATE KEY-----
MHcCAQEEIAwRtGkYqi732qh84HafnKE7YkW0CNpvvNseNGbxpsgGoAoGCCqGSM49
AwEHoUQDQgAE+xszAkYoKJP5CEvCaLuCGxAGDKIWecgPQxYElRWn/183SnpMHDRE
fXa4/Jzadq8dmo4taNQucoOLjD7IaN5OcA==
-----END PRIVATE KEY-----`
	token1 = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50IiwiaXNzIjoiZmFpcmhpdmUuaW8iLCJleHAiOjE2NDg1NTIsIm5iZiI6MTY0Nzk1MiwiaWF0IjoxNjQ3OTUyfQ.qk4RdhuwYDghwEtEkTS_3aoQF9zL9Cnb4Z0Pu7M6ALwJRZq-eCPGL-Q9_aP07oGszw8vcUZ82gw75CuX28oylQ"
)

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

func TestCreateECDSA(t *testing.T) {
	j, err := NewJWTECDSA(privateKey, jwt.SigningMethodES256)
	if err != nil {
		t.Errorf("error creating NewJWTECDSA: %v", err)
		t.FailNow()
	}

	ss, err := j.Create(u, time.UnixMicro(timestamp))
	if err != nil {
		t.Errorf("error creating ECDSA token: %v", err)
		t.FailNow()
	}
	if ss == "" {
		t.Errorf("incorrect signed string, cannot be empty")
		t.FailNow()
	}
	_, err = j.Extract(ss)
	if !errors.Is(err, ErrInvalidToken) { // expired token
		t.Errorf("incorrect error, err = %v, want %v", err, ErrInvalidToken)
		t.FailNow()
	}

	fmt.Println(ss)
}

func TestExtractECDSA(t *testing.T) {
	j, err := NewJWTECDSA(privateKey, jwt.SigningMethodES256)
	if err != nil {
		t.Errorf("error creating NewJWTECDSA: %v", err)
		t.FailNow()
	}
	_, err = j.Extract(token1)
	if !errors.Is(err, ErrInvalidToken) { // expired token
		t.Errorf("incorrect error, err = %v, want %v", err, ErrInvalidToken)
		t.FailNow()
	}
}

func TestRotationECDSA(t *testing.T) {
	j, err := NewJWTECDSA(privateKey, jwt.SigningMethodES256)
	if err != nil {
		t.Errorf("error creating NewJWTECDSA: %v", err)
		t.FailNow()
	}
	R := 1000
	now := time.Now()
	m := map[string]int{}
	for i := 0; i < R; i++ {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ss, err := j.Create(u, now)
			if err != nil {
				t.Errorf("error creating ECDSA token: %v", err)
				t.FailNow()
			}
			if ss == "" {
				t.Errorf("incorrect signed string, cannot be empty")
				t.FailNow()
			}
			v := m[ss] + 1
			if v != 1 {
				t.Errorf("%d x token %s", v, ss)
				t.FailNow()
			}
			m[ss] = v
			u2, err := j.Extract(ss)
			if errors.Is(err, ErrInvalidToken) {
				t.Errorf("incorrect error, err = %v, want nil", err)
				t.FailNow()
			}

			if u2.Address != address {
				t.Errorf("Address is incorrect, got %s, want %s", u2.Address, address)
				t.FailNow()
			}

			if u2.Email != email {
				t.Errorf("Email is incorrect, got %s, want %s", u2.Email, email)
				t.FailNow()
			}

			if !u2.HasSupportedType() {
				t.Errorf("Type is incorrect, got %s, want %s", u2.Type, utype)
				t.FailNow()
			}

			if u2.UUID == "" {
				t.Errorf("UUID is incorrect, cannot be empty string")
				t.FailNow()
			}

			if _, err := uuid.Parse(u2.UUID); err != nil {
				t.Errorf("UUID is incorrect, cannot be parsed: %v", err)
				t.FailNow()
			}

			if u2.Timestamp == 0 {
				t.Errorf("Timestamp is incorrect, cannot be 0")
				t.FailNow()
			}
		})
	}

}
