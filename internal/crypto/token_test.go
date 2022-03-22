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
		name  string
		time  time.Time
		token string
		user  *data.User
		jwt   Token
		err   error
	}{
		{
			"HS256",
			time.UnixMicro(timestamp),
			tokenHS256,
			u,
			NewJWTHS256(secret),
			nil,
		},
		{
			"HS512",
			time.UnixMicro(timestamp),
			tokenHS512,
			u,
			NewJWTHS512(secret),
			nil,
		},
		{
			"HS256 user missing address",
			time.UnixMicro(timestamp),
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50IiwiaXNzIjoiZmFpcmhpdmUuaW8iLCJleHAiOjE2NDg1NTIsIm5iZiI6MTY0Nzk1MiwiaWF0IjoxNjQ3OTUyfQ.fU_Vi8s1vxU59oRSAnTGj3bN4veMeNuFYBLNWKuOnJE",
			&data.User{
				Email: email,
				Type:  utype,
			},
			NewJWTHS256(secret),
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ss, err := tc.jwt.Create(tc.user, tc.time)
			if err != tc.err {
				t.Errorf("incorrect error, got %v, want %v", err, tc.err)
				t.Errorf("error creating user %v : %v", tc.user, err)
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
	tt := []struct {
		name                  string
		jwt                   Token
		ss                    string
		address, email, utype string
		err                   error
	}{
		{
			"valid token",
			j,
			ss,
			address, email, utype,
			err,
		},
		{
			"expired token",
			j,
			tokenHS256,
			"", "", "",
			ErrInvalidToken,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			user, err := tc.jwt.Extract(tc.ss)
			if err != tc.err {
				t.Errorf("incorrect error, got %v, want %v", err, tc.err)
				t.FailNow()
			}
			if user != nil {
				t.Logf("valid token : %v\n", tc.ss)
				if user.Address != tc.address {
					t.Errorf("incorrect address, got %v, want %v", user.Address, tc.address)
					t.FailNow()
				}
				if user.Email != tc.email {
					t.Errorf("incorrect email, got %v, want %v", user.Email, tc.email)
					t.FailNow()
				}
				if user.Type != tc.utype {
					t.Errorf("incorrect type, got %v, want %v", user.Type, tc.utype)
					t.FailNow()
				}
				if user.UUID == "" {
					t.Errorf("UUID is incorrect, cannot be empty string")
					t.FailNow()
				}
				if user.Timestamp == 0 {
					t.Errorf("Timestamp is incorrect, cannot be 0")
					t.FailNow()
				}
			}
		})
	}

	t.Run("token without address", func(t *testing.T) {
		ss, _ = j.Create(&data.User{
			Email: email,
			Type:  utype,
		}, now)
		_, err = j.Extract(ss)
		if err != ErrInvalidToken {
			t.Errorf("incorrect error, got %v, want %v", err, ErrInvalidToken)
			t.Errorf(ss)
			t.FailNow()
		}
	})

}

func TestHash(t *testing.T) {
	tt := []struct {
		name  string
		token string
		hash  string
	}{
		{
			"hash HS256 token",
			tokenHS256,
			"95257FA88B8C5A9DF58D83B8145446AF2E05FA9CB74733B78B275123AFB84743BF260ADBEAA19B134DEC9A25A3BB27E07BF00E68907A98116B23A40B81CB1E5B",
		},
		{
			"hash HS512 token",
			tokenHS512,
			"0525C1001A93DA4F9249985B79B7F47029F8AE601A024577EF6EE1060694712A7761CFBE2DBD2CFE19A6E6E97D5C90695117481B8DFDECEC94C81B941EA6A34F",
		},
		{
			"hash DATA",
			"DATA",
			"084E310EDCFBD2591B9997B55870D1AE49BCF1AEE7C74EFB4236CE8A9F28A6CE5FBF3394742969DFE578031822975EA44DE0C2AE68163368C8AA0185263FC874",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			h := hash(tc.token)
			if h != tc.hash {
				t.Errorf("incorrect hash, got %s, want %s", h, tc.hash)
				t.FailNow()
			}
		})
	}
}

func TestName(t *testing.T) {
	j := NewJWTHS256(secret)
	n := "HS256"

	if j.method.Alg() != n {
		t.Errorf("incorrect method name, got %s, want %s", j.method.Alg(), n)
		t.FailNow()
	}
}
