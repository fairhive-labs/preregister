package data

import (
	"testing"

	"github.com/google/uuid"
)

func TestSetup(t *testing.T) {

	u := User{}
	u.Setup()

	if u.UUID == "" {
		t.Errorf("UUID is incorrect, cannot be empty")
		t.FailNow()
	}

	if _, err := uuid.Parse(u.UUID); err != nil {
		t.Errorf("UUID is incorrect, cannot be parsed: %v", err)
		t.FailNow()
	}

	if u.Timestamp == 0 {
		t.Errorf("Timestamp is incorrect, cannot be set")
		t.FailNow()
	}

	if u.Validated {
		t.Errorf("Validated is incorrect, got %v, want %v", u.Validated, false)
		t.FailNow()
	}
}
