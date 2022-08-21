package data

import (
	"encoding/json"
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

func TestNewUser(t *testing.T) {
	address := "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	sponsor := "0x9ba1f109551bD432803012645Ac136ddd64DBA73"
	email := "john.doe@mailservice.com"
	utype := "talent"
	u := NewUser(address, email, utype, sponsor)

	if u.Address != address {
		t.Errorf("Address is incorrect, got %s, want %s", u.Address, address)
		t.FailNow()
	}

	if u.Email != email {
		t.Errorf("Email is incorrect, got %s, want %s", u.Email, email)
		t.FailNow()
	}

	if !u.HasSupportedType() {
		t.Errorf("Type is incorrect, got %s, want %s", u.Type, utype)
		t.FailNow()
	}

	if u.UUID == "" {
		t.Errorf("UUID is incorrect, cannot be empty string")
		t.FailNow()
	}

	if _, err := uuid.Parse(u.UUID); err != nil {
		t.Errorf("UUID is incorrect, cannot be parsed: %v", err)
		t.FailNow()
	}

	if u.Timestamp == 0 {
		t.Errorf("Timestamp is incorrect, cannot be 0")
		t.FailNow()
	}

	if u.Sponsor != sponsor {
		t.Errorf("Sponsor is incorrect, got %s, want %s", u.Sponsor, sponsor)
		t.FailNow()
	}

}

func TestMarshalling(t *testing.T) {
	address := "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	sponsor := "0x9ba1f109551bD432803012645Ac136ddd64DBA73"
	email := "john.doe@mailservice.com"
	utype := "talent"
	u := NewUser(address, email, utype, sponsor)
	jsonUser, _ := json.Marshal(u)
	n := 6
	m := make(map[string]interface{}, n)
	if err := json.Unmarshal(jsonUser, &m); err != nil {
		t.Errorf("Cannot unmarshal marshaled user: %v", err)
		t.FailNow()
	}
	if len(m) != n {
		t.Errorf("Incorrect number of field in json, got %d, want %d", len(m), n)
		t.FailNow()
	}

	u = &User{
		Address: address,
		Email:   email,
		Type:    utype,
		Sponsor: sponsor,
	}
	jsonUser, _ = json.Marshal(u)
	n = 4
	m = make(map[string]interface{}, n)
	if err := json.Unmarshal(jsonUser, &m); err != nil {
		t.Errorf("Cannot unmarshal marshaled user: %v", err)
		t.FailNow()
	}
	if len(m) != n {
		t.Errorf("Incorrect number of field in json, got %d, want %d", len(m), n)
		t.FailNow()
	}
}
