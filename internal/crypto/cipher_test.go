package crypto

import (
	"crypto/aes"
	"errors"
	"fmt"
	"testing"
)

const plaintext = "this is a secret message"

var (
	keys = map[int]string{
		32: "5a2c2e716f228ebdc8fa23d77e5a9ce8f3d46fe7944f439e8ea034832ab3a449",
		24: "42c12ae3b1f3bc00bb95ae635b4abbf2d18c5fb3b5e3093c",
		16: "a95c3bc19469a9cd8b0cf4d09dc04818",
	}
	ctexts = map[int]string{
		32: "f51b129d562d59cc7739aba99c62a91b55d5e5fd6f27652d322e6ce55eacf53fa49c435dd3620d7c9019390d00dce7e74d7b0311",
		24: "46757b2949d67efe480bd7893ede2d269b3896122d33f338060c3f1384d9eb74d264e9ffec17338fe767633ce4a09dc765f96ea2",
		16: "1f13498abd3eff31af5dc7fca99d112de541fd34cae974f468018f3c2a991f8330df70563216344cd5b335a2711c11474b9f98d9",
	}
)

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
	tt := []int{16, 24, 32}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%d", tc), func(t *testing.T) {
			_, err := Encrypt(plaintext, keys[tc])
			if err != nil {
				t.Errorf("incorrect error, got %v, want %v", err, nil)
				t.FailNow()
			}
		})
	}

	t.Run("invalid key", func(t *testing.T) {
		n := 12
		ks, _ := GenerateKey(n)
		_, err := Encrypt(plaintext, ks)
		if !errors.Is(err, aes.KeySizeError(n)) {
			t.Errorf("incorrect error, got %v, want %v", err, aes.KeySizeError(n))
			t.FailNow()
		}
	})
}

func TestDecrypt(t *testing.T) {
	tt := []int{16, 24, 32}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("valid key %d bytes", tc), func(t *testing.T) {
			txt, err := Decrypt(ctexts[tc], keys[tc])
			if err != nil {
				t.Errorf("incorrect error, got %v, want %v", err, nil)
				t.FailNow()
			}
			if txt != plaintext {
				t.Errorf("incorrect decrypted text, got %q, want %q", txt, plaintext)
				t.FailNow()
			}
		})
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("wrong key %d bytes", tc), func(t *testing.T) {
			ks, _ := GenerateKey(tc)
			txt, err := Decrypt(ctexts[tc], ks)
			if err == nil {
				t.Errorf("incorrect error, should be nil")
				t.FailNow()
			}
			if txt == plaintext {
				t.Errorf("decrypted text and plaintext cannot be equal")
				t.FailNow()
			}
		})
	}
}

func TestEncryptDecryptRotation(t *testing.T) {
	n := 1000
	s := []int{16, 24, 32}
	for i := 0; i < n; i++ {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			ks, _ := GenerateKey(s[i%len(s)])
			ctext, _ := Encrypt(plaintext, ks)
			txt, _ := Decrypt(ctext, ks)
			if txt != plaintext {
				t.Errorf("incorrect decrypted text, got %q, want %q", txt, plaintext)
				t.FailNow()
			}
		})
	}
}
