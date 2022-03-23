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
