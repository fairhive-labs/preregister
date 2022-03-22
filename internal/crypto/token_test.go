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

func TestNewJWTHS256(t *testing.T) {
	j := NewJWTHS256(secret)
	if secret != string(j.secret) {
		t.Errorf("incorrect NewJWTHS256 secret, got %s, want %s", j.secret, secret)
		t.FailNow()
	}
}

func TestCreate(t *testing.T) {

	tt := []struct {
		name  string
		time  time.Time
		token string
		jwt   Token
	}{
		{
			"HS256",
			time.UnixMicro(timestamp),
			tokenHS256,
			NewJWTHS256(secret),
		},
		{
			"HS512",
			time.UnixMicro(timestamp),
			tokenHS512,
			NewJWTHS512(secret),
		},
	}

	for _, tc := range tt {
		u := &data.User{
			Address: address,
			Email:   email,
			Type:    utype,
		}
		t.Run(tc.name, func(t *testing.T) {
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
}

func TestHash(t *testing.T) {

}
