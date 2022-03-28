package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func GenerateKey(size int) (s string, err error) {
	b := make([]byte, size)
	if _, err = io.ReadFull(rand.Reader, b); err != nil {
		return s, err
	}
	s = hex.EncodeToString(b)
	return
}

func Encrypt(plaintext, ks string) (string, error) {
	k, err := hex.DecodeString(ks)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", gcm.Seal(nonce, nonce, []byte(plaintext), nil)), nil
}
