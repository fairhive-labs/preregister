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
}

func TestNewUSer(t *testing.T) {
	address := "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	email := "john.doe@mailservice.com"
	utype := "talent"
	u := NewUser(address, email, utype)

	if u.Address != address {
		t.Errorf("Address is incorrect, got %s, want %s", u.Address, address)
		t.FailNow()
	}

	if u.Email != email {
		t.Errorf("Email is incorrect, got %s, want %s", u.Email, email)
		t.FailNow()
	}

	if u.Type != "talent" && u.Type != "agent" && u.Type != "mentor" && u.Type != "client" {
		t.Errorf("Type is incorrect, got %s, want %s", u.Type, utype)
		t.FailNow()
	}

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

}
