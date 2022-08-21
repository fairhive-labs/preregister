package crypto

import (
	"testing"

	"github.com/fairhive-labs/preregister/internal/data"
)

const (
	secret     = "VERY_SECURE_JWT_SECRET_L0L"
	address    = "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	sponsor    = "0x9ba1f109551bD432803012645Ac136ddd64DBA73"
	email      = "john.doe@mailservice.com"
	utype      = "talent"
	timestamp  = 1647952128425
	tokenHS256 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50Iiwic3BvbnNvciI6IjB4OWJhMWYxMDk1NTFiRDQzMjgwMzAxMjY0NUFjMTM2ZGRkNjREQkE3MyIsImlzcyI6ImZhaXJoaXZlLmlvIiwiZXhwIjoxNjQ4NTUyLCJuYmYiOjE2NDc5NTIsImlhdCI6MTY0Nzk1Mn0.-W9t_rVilW4cSts8k2bhrn2VicGAW2SYnF1ae8VyRe4"
	tokenHS512 = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50Iiwic3BvbnNvciI6IjB4OWJhMWYxMDk1NTFiRDQzMjgwMzAxMjY0NUFjMTM2ZGRkNjREQkE3MyIsImlzcyI6ImZhaXJoaXZlLmlvIiwiZXhwIjoxNjQ4NTUyLCJuYmYiOjE2NDc5NTIsImlhdCI6MTY0Nzk1Mn0.uICiapkj1hH9i4wXh9ztyhlf_teZ_6annIO4RQSvhUG45eZDb5eo5eyeq0AdhnxMrB64Vnm_96vKVR4qqkbVmA"
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
			"D886086FAAABBE8A4B35F4A0993142FFA05999AB0BDCE36D96A14EC43CD958ED5A20218CC195D89276FFE09A647DB0C92AD9B3AA93EA7BB98835B500EDB43370",
		},
		{
			"hash HS512 token",
			tokenHS512,
			"680F81D3D4B0452BF6B8E0CCAEE6A3C0B0AA39B00C9BCE508BB3DC76A9128D8907751D885439660E0119066CE2B1D2BD6B4AB5D9388386669D0EF91D20007C09",
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
