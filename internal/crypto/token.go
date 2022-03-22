package crypto

import (
	"time"

	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/golang-jwt/jwt/v4"
)

type Token interface {
	Create(user *data.User, t time.Time) (string, error)
	Extract(token string) (string, string, string)
	Hash(token string) string
	Name() string
}

type UserClaims struct {
	data.User
	jwt.RegisteredClaims
}

type JWTBase struct {
	method jwt.SigningMethod
}

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
	claims := UserClaims{
		*user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t.Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(t),
			NotBefore: jwt.NewNumericDate(t),
			Issuer:    "fairhive.io",
		},
	}
	token := jwt.NewWithClaims(j.method, claims)
	return token.SignedString(j.secret)
}

func (j JWTHS) Extract(token string) (string, string, string) {
	return "", "", ""
}

func (j JWTHS) Hash(token string) string {
	return "hash"
}

func (j JWTHS) Name() string {
	return j.method.Alg()
}
