package crypto

import (
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

type JWTBase struct {
	method jwt.SigningMethod
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

func (JWTBase) Hash(token string) string {
	return hash(token)
}

func create(user *data.User, t time.Time, m jwt.SigningMethod, k interface{}) (string, error) {
	claims := UserClaims{
		*user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t.Add(10 * time.Minute)), // seconds
			IssuedAt:  jwt.NewNumericDate(t),                       // seconds
			NotBefore: jwt.NewNumericDate(t),                       // seconds
			Issuer:    "fairhive.io",
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

func extract[T *jwt.SigningMethodHMAC | *jwt.SigningMethodECDSA](token string, k interface{}) (u *data.User, err error) {
	uclaims := &UserClaims{}
	tk, err := jwt.ParseWithClaims(token, uclaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(T); !ok {
			fmt.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return k, nil
	})

	if tk.Valid &&
		uclaims.Address != "" &&
		uclaims.Email != "" &&
		uclaims.Type != "" { // mandatory Fields...
		return data.NewUser(uclaims.Address, uclaims.Email, uclaims.Type), nil
	}
	err = ErrInvalidToken
	return
}
