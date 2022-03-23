package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/golang-jwt/jwt/v4"
)

type JWTECDSA struct {
	JWTBase[*ecdsa.PrivateKey]
}

func NewJWTES256() (*JWTECDSA, error) {
	pvk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &JWTECDSA{JWTBase[*ecdsa.PrivateKey]{jwt.SigningMethodES256, pvk}}, nil
}

func NewJWTECDSA(k string, m *jwt.SigningMethodECDSA) (*JWTECDSA, error) {
	pvk, err := jwt.ParseECPrivateKeyFromPEM([]byte(k))
	if err != nil {
		return nil, err
	}
	return &JWTECDSA{JWTBase[*ecdsa.PrivateKey]{m, pvk}}, nil
}

func (j JWTECDSA) Extract(token string) (u *data.User, err error) {
	return extract[*jwt.SigningMethodECDSA](token, j.k)
}
