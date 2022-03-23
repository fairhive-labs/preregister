package crypto

import (
	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/golang-jwt/jwt/v4"
)

type JWTHMAC struct {
	JWTBase[[]byte]
}

func NewJWTHS256(s string) *JWTHMAC {
	return &JWTHMAC{JWTBase[[]byte]{jwt.SigningMethodHS256, []byte(s)}}
}

func NewJWTHS512(s string) *JWTHMAC {
	return &JWTHMAC{JWTBase[[]byte]{jwt.SigningMethodHS512, []byte(s)}}
}

func (j JWTHMAC) Extract(token string) (u *data.User, err error) {
	return extract[*jwt.SigningMethodHMAC](token, j.k)
}
