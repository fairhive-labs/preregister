package data

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

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

// Test NewUser() and defacto IsValid()
func TestNewUser(t *testing.T) {
	type errorDetails struct {
		field string
		tag   string
		value interface{}
	}

	validUser := NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doe@mailservice.com", "talent", sponsor)

	noUUIDUser := *validUser
	noUUIDUser.UUID = ""

	invalidUUIDUser := *validUser
	invalidUUIDUser.UUID = "fakeUUID"

	invalidTimestampUser := *validUser
	invalidTimestampUser.Timestamp = 0

	tt := []struct {
		name string
		u    *User
		err  *errorDetails
		v    bool
	}{
		{"valid_user", validUser, nil, true},
		{"invalid_user_address",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA73", "john.doe@mailservice.com", "talent", sponsor),
			&errorDetails{"Address", "eth_addr", "0x8ba1f109551bD432803012645Ac136ddd64DBA73"},
			false,
		},
		{"missing_user_address",
			NewUser("", "john.doe@mailservice.com", "talent", sponsor),
			&errorDetails{"Address", "required", ""},
			false,
		},
		{"invalid_email",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemailservice.com", "talent", sponsor),
			&errorDetails{"Email", "email", "john.doemailservice.com"},
			false,
		},
		{"missing_email",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "", "talent", sponsor),
			&errorDetails{"Email", "required", ""},
			false,
		},
		{"invalid_sponsor",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemail@service.com", "talent", "0x8ba1f109551bD432803012645Ac136ddd64DBA73"),
			&errorDetails{"Sponsor", "eth_addr", "0x8ba1f109551bD432803012645Ac136ddd64DBA73"},
			false,
		},
		{"missing_sponsor",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemail@service.com", "talent", ""),
			&errorDetails{"Sponsor", "required", ""},
			false,
		},
		{"missing_type",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemail@service.com", "", sponsor),
			&errorDetails{"Type", "required", ""},
			false,
		},
		{"invalid_type",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemail@service.com", "unsupported_type", sponsor),
			&errorDetails{"Type", "oneof", "unsupported_type"},
			false,
		},
		{"missing_uuid",
			&noUUIDUser,
			&errorDetails{"UUID", "required", ""},
			false,
		},
		{"invalid_uuid",
			&invalidUUIDUser,
			&errorDetails{"UUID", "uuid", "fakeUUID"},
			false,
		},
		{"invalid_timestamp",
			&invalidTimestampUser,
			&errorDetails{"Timestamp", "gt", int64(0)},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			u := tc.u

			err := validate.Struct(u)
			if err != nil {
				if tc.err != nil {
					for _, e := range err.(validator.ValidationErrors) {
						if e.Field() != tc.err.field {
							t.Errorf("Field is incorrect, got %v, want %v", e.Field(), tc.err.field)
						}
						if e.Tag() != tc.err.tag {
							t.Errorf("Tag is incorrect, got %v, want %v", e.Tag(), tc.err.tag)
						}
						if e.Value() != tc.err.value {
							t.Errorf("Value is incorrect, got %v, want %v", e.Value(), tc.err.value)
						}
					}
				} else {
					t.Errorf("user %v is not valid, %v ", *u, err)
					t.FailNow()
				}
			}

			if u.IsValid() != tc.v {
				t.Errorf("incorrect validation, got %v, want %v", u.IsValid(), tc.v)
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
