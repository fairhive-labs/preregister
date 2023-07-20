package data

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	Address   string `json:"address" binding:"required,eth_addr" validate:"required,eth_addr"`
	Email     string `json:"email" binding:"required,email" validate:"required,email"`
	UUID      string `json:"uuid,omitempty" validate:"required,uuid"`
	Timestamp int64  `json:"timestamp,omitempty" validate:"gt=0"`
	Type      string `json:"type" binding:"required,oneof=advisor agent initiator contributor investor mentor contractor" validate:"required,oneof=advisor agent initiator contributor investor mentor contractor"`
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

// IsValid tests if all fields are valid
func (u *User) IsValid() bool {
	return nil == validate.Struct(u)
}

// IsSet tests if only required fields are valid
func (u *User) IsSet() bool {
	err := validate.StructExcept(u, "UUID", "Timestamp")
	if err != nil {
		log.Print(err)
	}
	return nil == err
}

func (u User) String() string {
	r, _ := json.Marshal(&struct {
		*User
		Timestamp string `json:"timestamp"`
	}{
		User:      &u,
		Timestamp: time.UnixMilli(u.Timestamp).Format(time.RFC3339Nano),
	})
	return string(r)
}
