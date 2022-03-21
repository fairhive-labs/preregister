package data

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Address   string `json:"address" binding:"required,eth_addr"`
	Email     string `json:"email" binding:"required,email"`
	UUID      string `json:"uuid"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type" binding:"required,oneof=advisor agent client contributor investor mentor talent"`
}

func (u *User) Setup() {
	u.UUID = uuid.New().String()
	u.Timestamp = time.Now().UnixMilli()
}

func NewUser(a, e, t string) *User {
	u := &User{
		Address: a,
		Email:   e,
		Type:    t,
	}
	u.Setup()
	return u
}

func (u *User) IsSupportedUser() bool {
	switch u.Type {
	case "advisor", "agent", "client", "contributor", "investor", "mentor", "talent":
		return true
	default:
		return false
	}
}
