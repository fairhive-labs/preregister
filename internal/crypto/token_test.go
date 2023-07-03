package crypto

import (
	"testing"

	"github.com/fairhive-labs/preregister/internal/data"
)

const (
	secret     = "VERY_SECURE_JWT_SECRET_L0L"
	address    = "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	sponsor    = "0xD01efFE216E16a85Fc529db66c26aBeCf4D885f8" // real address but empty balance
	email      = "john.doe@mailservice.com"
	utype      = "talent"
	timestamp  = 1647952128425
	tokenHS256 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50Iiwic3BvbnNvciI6IjB4RDAxZWZGRTIxNkUxNmE4NUZjNTI5ZGI2NmMyNmFCZUNmNEQ4ODVmOCIsImlzcyI6InBvbG4ub3JnIiwiZXhwIjoxNjQ4NTUyLCJuYmYiOjE2NDc5NTIsImlhdCI6MTY0Nzk1Mn0.6k7_ZPyuAekpHrZ4x4gpBNcbx88aLbRKToLYSlUAY0I"
	tokenHS512 = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50Iiwic3BvbnNvciI6IjB4RDAxZWZGRTIxNkUxNmE4NUZjNTI5ZGI2NmMyNmFCZUNmNEQ4ODVmOCIsImlzcyI6InBvbG4ub3JnIiwiZXhwIjoxNjQ4NTUyLCJuYmYiOjE2NDc5NTIsImlhdCI6MTY0Nzk1Mn0.I4zymcstzyrI9Q51FS3ejCbHbk88EgjiRLVi9osN2W1xYxUUXnjOdMYCZgm7GAGrpDTyRLHsyNql2V6yMh6JZg"
)

var (
	u = &data.User{
		Address: address,
		Email:   email,
		Type:    utype,
		Sponsor: sponsor,
	}
)

func TestHash(t *testing.T) {
	tt := []struct {
		name  string
		token string
		hash  string
	}{
		{
			"hash HS256 token",
			tokenHS256,
			"14D40231970546854D1F94BEFBB7BD4C6C73D622AB7D44B7C223D249CAFE404A392C94665FCC5472685EA873D34807F088216AAA0FA969E25D13E5CBF4138BF6",
		},
		{
			"hash HS512 token",
			tokenHS512,
			"0DCEDEE5D4B79F4CEF97DCE4A43E2B8D445E8AED627A0201CAB34FAEA287E2C3A9E8A57E9EB025326000E1AEF8F6F2C5C9F748094477E4A1AA446E4603A747B3",
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
