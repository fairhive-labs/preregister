package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

var ErrTooShortCipherText = errors.New("ciphertext too short")

func GenerateKey(size int) (s string, err error) {
	b := make([]byte, size)
	if _, err = io.ReadFull(rand.Reader, b); err != nil {
		return s, err
	}
	s = hex.EncodeToString(b)
	return
}

func Encrypt(text, ks string) (string, error) {
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
	return fmt.Sprintf("%x", gcm.Seal(nonce, nonce, []byte(text), nil)), nil
}

func Decrypt(ctext, ks string) (string, error) {
	k, err := hex.DecodeString(ks)
	if err != nil {
		return "", err
	}
	enc, err := hex.DecodeString(ctext)
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

	nonceSize := gcm.NonceSize()
	if len(ctext) < nonceSize {
		return "nil", ErrTooShortCipherText
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	b, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
