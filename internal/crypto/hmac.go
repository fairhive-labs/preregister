package crypto

import (
	"time"

	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/golang-jwt/jwt/v4"
)

type JWTHMAC struct {
	k []byte
	JWTBase
}

func NewJWTHS256(s string) *JWTHMAC {
	return &JWTHMAC{[]byte(s), JWTBase{jwt.SigningMethodHS256}}
}

func NewJWTHS512(s string) *JWTHMAC {
	return &JWTHMAC{[]byte(s), JWTBase{jwt.SigningMethodHS512}}
}

func (j JWTHMAC) Create(user *data.User, t time.Time) (string, error) {
	return create(user, t, j.method, j.k)
}

func (j JWTHMAC) Extract(token string) (u *data.User, err error) {
	return extract[*jwt.SigningMethodHMAC](token, j.k)
}
