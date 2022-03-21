package crypto

import (
	"testing"
)

func TestNewJWTHS256(t *testing.T) {
	s := "VERY_SECURE_JWT_SECRET_L0L"
	j := NewJWTHS256(s)
	if s != string(j.secret) {
		t.Errorf("incorrect NewJWTHS256 secret, got %s, want %s", j.secret, s)
		t.FailNow()
	}
}
