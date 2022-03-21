package crypto

import "github.com/fairhive-labs/preregister/internal/data"

type Token interface {
	Create(u *data.User) string
	Hash(t string) string
}

type HS256 struct {
	secret string
}
