package data

type User struct {
	Address   string `json:"address" binding:"required,eth_addr"`
	Email     string `json:"email" binding:"required,email"`
	UUID      string `json:"uuid"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type" binding:"required,oneof=client talent agent mentor"`
	Validated bool   `json:"validated"`
}
