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
	Type      string `json:"type" binding:"required,oneof=client talent agent mentor"`
	Validated bool   `json:"validated"`
}

func (u *User) Setup() {
	u.UUID = uuid.New().String()
	u.Timestamp = time.Now().UnixMilli()
	u.Validated = false
}
