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
	tokenHS256 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50Iiwic3BvbnNvciI6IjB4RDAxZWZGRTIxNkUxNmE4NUZjNTI5ZGI2NmMyNmFCZUNmNEQ4ODVmOCIsImlzcyI6ImZhaXJoaXZlLmlvIiwiZXhwIjoxNjQ4NTUyLCJuYmYiOjE2NDc5NTIsImlhdCI6MTY0Nzk1Mn0.2E-pDxS9wak4prF00WYMM0aLYJSUejYqHi1DdgW0hXo"
	tokenHS512 = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50Iiwic3BvbnNvciI6IjB4RDAxZWZGRTIxNkUxNmE4NUZjNTI5ZGI2NmMyNmFCZUNmNEQ4ODVmOCIsImlzcyI6ImZhaXJoaXZlLmlvIiwiZXhwIjoxNjQ4NTUyLCJuYmYiOjE2NDc5NTIsImlhdCI6MTY0Nzk1Mn0.uLJ02V1kBcE8Dqy1Fer6F94FJDSbsIDpuIYtap7OZvr3FhPK3MOVb-Pxt2bK-jHuh6dLNa4ftyQSpIjLro8IBA"
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
			"3807327CC5167677D06A0A9E710D7AC84C20129C0CA053E103ACD27E08C963D08DBD376B29BD03EFEE4BCEB1BB6E70C3F5C45D966E5FF1C08BD185A44B54F102",
		},
		{
			"hash HS512 token",
			tokenHS512,
			"5ED93B53FA32AA5DC7A2D7BB10B746EB43A1F330C97A12E7B0C465738612EA7D37707411FE729EF1E1191D7F09DB1654ED502BFDBAFC02BB917AB6A8529A81E5",
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
