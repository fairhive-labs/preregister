package crypto

import (
	"fmt"
	"testing"
)

var ks = map[int]string{
	32: "5a2c2e716f228ebdc8fa23d77e5a9ce8f3d46fe7944f439e8ea034832ab3a449",
	24: "42c12ae3b1f3bc00bb95ae635b4abbf2d18c5fb3b5e3093c",
	16: "a95c3bc19469a9cd8b0cf4d09dc04818",
}

func TestGenerateKey(t *testing.T) {
	tt := []struct {
		s int
		l int
	}{
		{32, 64},
		{24, 48},
		{16, 32},
		{10, 20},
		{12, 24},
		{128, 256},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%d", tc.s), func(t *testing.T) {
			k, err := GenerateKey(tc.s)
			t.Logf("GenerateKey(%d)=%s\n", tc.s, k)
			if err != nil {
				t.Errorf("incorrect error, got %v, want nil", err)
				t.FailNow()
			}
			if len(k) != tc.l {
				t.Errorf("incorrect length, got %d, want %d", len(k), tc.l)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	plaintext := "this is a secret message"

	tt := []struct {
		ks  string
		err error
	}{
		{
			ks[32],
			nil,
		},
		{
			ks[24],
			nil,
		},
		{
			ks[16],
			nil,
		},
	}

	for i, tc := range tt {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			es, err := Encrypt(plaintext, tc.ks)
			if err != tc.err {
				t.Errorf("incorrect error, got %v, want %v", err, tc.err)
				t.FailNow()
			}
			if "" == es {
				t.Errorf("incorrect encrypted string, cannot be empty string")
			}
		})
	}
}
