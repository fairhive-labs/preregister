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
	utype      = "contractor"
	timestamp  = 1647952128425
	tokenHS256 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoiY29udHJhY3RvciIsInNwb25zb3IiOiIweEQwMWVmRkUyMTZFMTZhODVGYzUyOWRiNjZjMjZhQmVDZjREODg1ZjgiLCJpc3MiOiJwb2xuLm9yZyIsImV4cCI6MTY0ODU1MiwibmJmIjoxNjQ3OTUyLCJpYXQiOjE2NDc5NTJ9.fcTXYvP9bAN934OjZnftZchCS6238J9TH9MkVT5c1sQ"
	tokenHS512 = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiMHg4YmExZjEwOTU1MWJENDMyODAzMDEyNjQ1QWMxMzZkZGQ2NERCQTcyIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoiY29udHJhY3RvciIsInNwb25zb3IiOiIweEQwMWVmRkUyMTZFMTZhODVGYzUyOWRiNjZjMjZhQmVDZjREODg1ZjgiLCJpc3MiOiJwb2xuLm9yZyIsImV4cCI6MTY0ODU1MiwibmJmIjoxNjQ3OTUyLCJpYXQiOjE2NDc5NTJ9.fQqiz03S8DBEkoolEBFJ4eS09KBBdaTm7hSgiUUXeRUoptgxB279wU2sA4BtDDn8NieTyxxT-WZDECLEXTrkOA"
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
			"E6D97C289D1D0CA200CAC40112FEFE044CFFCE5D3BF1F122C598ED6E035D46295BBEB16E966260BE16628CEBBB3F6457D9DB22A7A7F4BE668C3D3755A36C94DC",
		},
		{
			"hash HS512 token",
			tokenHS512,
			"A3DDCFFF0A67F9461B6200BD17814633BF5C355174F54C5E534A758008297F500A94034D01BE54EECDBB45878BEA5C69F36EB22F618911295A5E9BBE388B4065",
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
