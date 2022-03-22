package crypto

import (
	"testing"
	"time"

	"github.com/fairhive-labs/preregister/internal/data"
)

const (
	secret     = "VERY_SECURE_JWT_SECRET_L0L"
	address    = "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	email      = "john.doe@mailservice.com"
	utype      = "talent"
	timestamp  = 1647952128425
	tokenHS256 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50IiwiaXNzIjoiZmFpcmhpdmUuaW8iLCJleHAiOjE2NDg1NTIsIm5iZiI6MTY0Nzk1MiwiaWF0IjoxNjQ3OTUyfQ.VHl3lWyBfXvH0q80gwPO6OgrUnKen_07yk63IS6u6XM"
	tokenHS512 = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50IiwiaXNzIjoiZmFpcmhpdmUuaW8iLCJleHAiOjE2NDg1NTIsIm5iZiI6MTY0Nzk1MiwiaWF0IjoxNjQ3OTUyfQ.Dw_v9S0akP4xU8-3N2AcElUP_uhAnso619Ri5EKLESKWGIL0j_Cbb8JY7RrefWCe_Y4RkWdyhPY5VgHbGqPTug"
)

var (
	u = &data.User{
		Address: address,
		Email:   email,
		Type:    utype,
	}
)

func TestNewJWTHS256(t *testing.T) {
	j := NewJWTHS256(secret)
	if secret != string(j.secret) {
		t.Errorf("incorrect NewJWTHS256 secret, got %s, want %s", j.secret, secret)
		t.FailNow()
	}
}

func TestCreate(t *testing.T) {

	tt := []struct {
		time  time.Time
		token string
		jwt   Token
	}{
		{
			time.UnixMicro(timestamp),
			tokenHS256,
			NewJWTHS256(secret),
		},
		{
			time.UnixMicro(timestamp),
			tokenHS512,
			NewJWTHS512(secret),
		},
	}

	for _, tc := range tt {
		t.Run(tc.jwt.Name(), func(t *testing.T) {
			ss, err := tc.jwt.Create(u, tc.time)
			if err != nil {
				t.Errorf("error creating user %v : %v", *u, err)
				t.FailNow()
			}

			if ss != tc.token {
				t.Errorf("incorrect token, got %s, want %s", ss, tc.token)
				t.FailNow()
			}
		})
	}
}

func TestExtract(t *testing.T) {
	j := NewJWTHS256(secret)
	now := time.Now()
	ss, err := j.Create(u, now)
	if err != nil {
		t.Errorf("error creating user %v : %v", *u, err)
		t.FailNow()
	}

	a, e, ut, err := j.Extract(ss)
	if err != nil {
		t.Errorf("error extracting JWT %v : %v", ss, err)
		t.FailNow()
	}
	if a != address {
		t.Errorf("incorrect address, got %v, want %v", a, address)
		t.FailNow()
	}
	if e != email {
		t.Errorf("incorrect email, got %v, want %v", e, email)
		t.FailNow()
	}
	if ut != utype {
		t.Errorf("incorrect type, got %v, want %v", ut, utype)
		t.FailNow()
	}
}

func TestHash(t *testing.T) {

}

func TestName(t *testing.T) {
	j := NewJWTHS256(secret)
	n := "HS256"

	if j.method.Alg() != n {
		t.Errorf("incorrect method name, got %s, want %s", j.method.Alg(), n)
		t.FailNow()
	}
}
