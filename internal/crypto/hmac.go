package crypto

import (
	"time"

	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/golang-jwt/jwt/v4"
)

type JWTHS struct {
	secret []byte
	JWTBase
}

func NewJWTHS256(s string) *JWTHS {
	return &JWTHS{[]byte(s), JWTBase{jwt.SigningMethodHS256}}
}

func NewJWTHS512(s string) *JWTHS {
	return &JWTHS{[]byte(s), JWTBase{jwt.SigningMethodHS512}}
}

func (j JWTHS) Create(user *data.User, t time.Time) (string, error) {
	return create(user, t, j.method, j.secret)
}

func (j JWTHS) Extract(token string) (u *data.User, err error) {
	return extract[*jwt.SigningMethodHMAC](token, j.secret)
}
