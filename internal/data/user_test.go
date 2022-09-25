package data

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

const sponsor = "0xD01efFE216E16a85Fc529db66c26aBeCf4D885f8" // real address but empty balance

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

	tt := []struct {
		name string
		u    *User
		err  *struct {
			field string
			tag   string
			value interface{}
			param string
		}
	}{
		{"valid_user", NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doe@mailservice.com", "talent", sponsor), nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			u := tc.u

			err := validate.Struct(u)
			if err != nil {
				if tc.err != nil {
					for _, err := range err.(validator.ValidationErrors) {

						fmt.Println(err.Field())
						fmt.Println(err.Tag())
						fmt.Println(err.Value())
						fmt.Println(err.Param())
					}
				} else {
					t.Errorf("user %v is not valid, %v ", *u, err)
					t.FailNow()
				}
			}

			if _, err := uuid.Parse(u.UUID); err != nil {
				t.Errorf("UUID is incorrect, cannot be parsed: %v", err)
				t.FailNow()
			}

			if u.Timestamp == 0 {
				t.Errorf("Timestamp is incorrect, cannot be 0")
				t.FailNow()
			}
		})
	}

}

func TestMarshalling(t *testing.T) {
	address := "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
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
