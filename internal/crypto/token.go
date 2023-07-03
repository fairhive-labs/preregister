package crypto

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/sha3"
)

type Token interface {
	Create(user *data.User, t time.Time) (string, error)
	Extract(token string) (*data.User, error)
	Hash(token string) string
}

type KeyConstraint interface {
	[]byte | *ecdsa.PrivateKey
}

type JWTBase[K KeyConstraint] struct {
	method jwt.SigningMethod
	k      K
}

type UserClaims struct {
	data.User
	jwt.RegisteredClaims
}

var (
	ErrSigningToken = errors.New("cannot sign token")
	ErrInvalidToken = errors.New("invalid token")
)

func hash(token string) string {
	return fmt.Sprintf("%X", sha3.Sum512([]byte(token)))
}

func (JWTBase[K]) Hash(token string) string {
	return hash(token)
}

func create(user *data.User, t time.Time, m jwt.SigningMethod, k interface{}) (string, error) {
	claims := UserClaims{
		*user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t.Add(10 * time.Minute)), // seconds
			IssuedAt:  jwt.NewNumericDate(t),                       // seconds
			NotBefore: jwt.NewNumericDate(t),                       // seconds
			Issuer:    "poln.org",
		},
	}
	token := jwt.NewWithClaims(m, claims)
	ss, err := token.SignedString(k)
	if err != nil {
		fmt.Printf("error creating token for user %v : %v", *user, err)
		err = ErrSigningToken
	}
	return ss, err
}

func (j JWTBase[K]) Create(user *data.User, t time.Time) (string, error) {
	return create(user, t, j.method, j.k)
}

func extract[SM *jwt.SigningMethodHMAC | *jwt.SigningMethodECDSA](token string, k interface{}) (u *data.User, err error) {
	uclaims := &UserClaims{}
	tk, _ := jwt.ParseWithClaims(token, uclaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(SM); !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return k, nil
	})

	if tk.Valid && uclaims.IsSet() {
		return data.NewUser(uclaims.Address, uclaims.Email, uclaims.Type, uclaims.Sponsor), nil
	}
	//fmt.Printf("Error extracting JWT: %v\n", err)
	err = ErrInvalidToken
	return
}
