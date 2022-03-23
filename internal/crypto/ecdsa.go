package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/golang-jwt/jwt/v4"
)

type JWTECDSA struct {
	k *ecdsa.PrivateKey
	JWTBase
}

func NewJWTES256() (*JWTECDSA, error) {
	pvk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &JWTECDSA{pvk, JWTBase{jwt.SigningMethodES256}}, nil
}

func NewJWTECDSA(k string, m *jwt.SigningMethodECDSA) (*JWTECDSA, error) {
	pvk, err := jwt.ParseECPrivateKeyFromPEM([]byte(k))
	if err != nil {
		return nil, err
	}
	return &JWTECDSA{pvk, JWTBase{m}}, nil
}
