package crypto

import (
	"fmt"
	"time"

	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/golang-jwt/jwt/v4"
)

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
	ss, err := token.SignedString(j.secret)
	if err != nil {
		fmt.Printf("error creating token for user %v : %v", *user, err)
		err = ErrSigningToken
	}
	return ss, err
}

func (j JWTHS) Extract(token string) (u *data.User, err error) {
	uclaims := &UserClaims{}
	tk, err := jwt.ParseWithClaims(token, uclaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return j.secret, nil
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

func (JWTHS) Hash(token string) string {
	return hash(token)
}
