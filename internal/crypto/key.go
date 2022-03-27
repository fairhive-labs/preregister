package crypto

import (
	"crypto/rand"
	"encoding/hex"
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
