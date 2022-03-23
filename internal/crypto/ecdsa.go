package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type JWTECDSA struct {
	k *ecdsa.PrivateKey
	JWTBase
}

func NewJWTES256() *JWTECDSA {
	rng := rand.Reader
	pvk, err := ecdsa.GenerateKey(elliptic.P256(), rng)
	if err != nil {
		fmt.Println("cannot generate ECDSA key")
		os.Exit(1)
	}
	return &JWTECDSA{pvk, JWTBase{jwt.SigningMethodES256}}
}
