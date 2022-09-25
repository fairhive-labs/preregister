package data

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	Address   string `json:"address" binding:"required,eth_addr" validate:"required,eth_addr"`
	Email     string `json:"email" binding:"required,email" validate:"required,email"`
	UUID      string `json:"uuid,omitempty" validate:"required,uuid"`
	Timestamp int64  `json:"timestamp,omitempty" validate:"gt=0"`
	Type      string `json:"type" binding:"required,oneof=advisor agent client contributor investor mentor talent" validate:"required,oneof=advisor agent client contributor investor mentor talent"`
	Sponsor   string `json:"sponsor" binding:"required,eth_addr" validate:"required,eth_addr"`
}

var validate = validator.New()

func (u *User) Setup() {
	u.UUID = uuid.New().String()
	u.Timestamp = time.Now().UnixMilli()
}

func NewUser(a, e, t, s string) *User {
	u := &User{
		Address: a,
		Email:   e,
		Type:    t,
		Sponsor: s,
	}
	u.Setup()
	return u
}

func (u *User) HasSupportedType() bool {
	switch u.Type {
	case "advisor", "agent", "client", "contributor", "investor", "mentor", "talent":
		return true
	default:
		return false
	}
}

// IsValid tests if all fields are valid
func (u *User) IsValid() bool {
	return nil == validate.Struct(u)
}

// IsSet tests if only required fields are valid
func (u *User) IsSet() bool {
	return nil == validate.StructExcept(u, "UUID", "Timestamp")
}
