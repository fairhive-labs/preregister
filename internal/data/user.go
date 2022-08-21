package data

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Address   string `json:"address" binding:"required,eth_addr"`
	Email     string `json:"email" binding:"required,email"`
	UUID      string `json:"uuid,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Type      string `json:"type" binding:"required,oneof=advisor agent client contributor investor mentor talent"`
	Sponsor   string `json:"sponsor" binding:"required,eth_addr"`
}

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
