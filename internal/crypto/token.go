package crypto

import "github.com/fairhive-labs/preregister/internal/data"

type Token interface {
	Create(user *data.User) string
	Extract(token string) (string, string, string)
	Hash(token string) string
}

type JWTHS256 struct {
	secret []byte
}

func NewJWTHS256(s string) *JWTHS256 {
	return &JWTHS256{[]byte(s)}
}

func (j *JWTHS256) Create(user *data.User) string {
	return ""
}

func (j *JWTHS256) Extract(token string) (string, string, string) {
	return "", "", ""
}

func (j *JWTHS256) Hash(token string) string {
	return "hash"
}
