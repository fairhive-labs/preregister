package crypto

import "crypto/ecdsa"

type JWTECDSA struct {
	k *ecdsa.PrivateKey
	JWTBase
}
