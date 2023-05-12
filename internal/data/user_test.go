package data

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const sponsor = "0xE3C3691DB5f5185F37A3f98e5ec76403B2d10c3E" // trendev eth address

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

// Test NewUser() + defacto IsValid() & IsSet()
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
		name       string
		u          *User
		err        *errorDetails
		valid, set bool
	}{
		{"valid_user", validUser, nil, true, true},
		{"invalid_user_address",
			NewUser("0x8bz1f109551bD432803012645Ac136ddd64DBA73", "john.doe@mailservice.com", "talent", sponsor),
			&errorDetails{"Address", "eth_addr", "0x8bz1f109551bD432803012645Ac136ddd64DBA73"},
			false, false,
		},
		{"missing_user_address",
			NewUser("", "john.doe@mailservice.com", "talent", sponsor),
			&errorDetails{"Address", "required", ""},
			false, false,
		},
		{"invalid_email",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemailservice.com", "talent", sponsor),
			&errorDetails{"Email", "email", "john.doemailservice.com"},
			false, false,
		},
		{"missing_email",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "", "talent", sponsor),
			&errorDetails{"Email", "required", ""},
			false, false,
		},
		{"invalid_sponsor",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemail@service.com", "talent", "0x8bz1f109551bD432803012645Ac136ddd64DBA73"),
			&errorDetails{"Sponsor", "eth_addr", "0x8bz1f109551bD432803012645Ac136ddd64DBA73"},
			false, false,
		},
		{"missing_sponsor",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemail@service.com", "talent", ""),
			&errorDetails{"Sponsor", "required", ""},
			false, false,
		},
		{"missing_type",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemail@service.com", "", sponsor),
			&errorDetails{"Type", "required", ""},
			false, false,
		},
		{"invalid_type",
			NewUser("0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doemail@service.com", "unsupported_type", sponsor),
			&errorDetails{"Type", "oneof", "unsupported_type"},
			false, false,
		},
		{"missing_uuid",
			&noUUIDUser,
			&errorDetails{"UUID", "required", ""},
			false, true,
		},
		{"invalid_uuid",
			&invalidUUIDUser,
			&errorDetails{"UUID", "uuid", "fakeUUID"},
			false, true,
		},
		{"invalid_timestamp",
			&invalidTimestampUser,
			&errorDetails{"Timestamp", "gt", int64(0)},
			false, true,
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

			if u.IsValid() != tc.valid {
				t.Errorf("incorrect complete validation, got %v, want %v", u.IsValid(), tc.valid)
			}

			if u.IsSet() != tc.set {
				t.Errorf("incorrect partial validation, got %v, want %v", u.IsSet(), tc.set)
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

func TestString(t *testing.T) {
	a1, a2 := "0xaD51c5ac7612DB8dD1611c6B2e317E4950c40942", "0x9C93c71065ea9101F252dE2e0f277437f473ac04"
	e1, e2 := "user1@domain.com", "user2@domain.com"
	id1, id2 := "4a8e9808-563e-4761-a8fa-305fef099a3e", "942a5811-926d-4014-baff-ef707f38407e"
	tm1, tm2 := 1683907220519, 1683807190432
	ty1, ty2 := "talent", "client"
	s1, s2 := "0x095cb719f8f69952599c15af31c80Ccb825E15d4", "0x233F858EaF43AFFE5DDFBD3AD69ACc6f5de6C529"

	tt := []struct {
		name string
		u    *User
		exp  string
	}{
		{
			"valid_user1",
			&User{a1, e1, id1, int64(tm1), ty1, s1},
			"{\"address\":\"0xaD51c5ac7612DB8dD1611c6B2e317E4950c40942\",\"email\":\"user1@domain.com\",\"uuid\":\"4a8e9808-563e-4761-a8fa-305fef099a3e\",\"type\":\"talent\",\"sponsor\":\"0x095cb719f8f69952599c15af31c80Ccb825E15d4\",\"timestamp\":\"2023-05-12T18:00:20.519+02:00\"}",
		},
		{
			"valid_user2",
			&User{a2, e2, id2, int64(tm2), ty2, s2},
			"{\"address\":\"0x9C93c71065ea9101F252dE2e0f277437f473ac04\",\"email\":\"user2@domain.com\",\"uuid\":\"942a5811-926d-4014-baff-ef707f38407e\",\"type\":\"client\",\"sponsor\":\"0x233F858EaF43AFFE5DDFBD3AD69ACc6f5de6C529\",\"timestamp\":\"2023-05-11T14:13:10.432+02:00\"}",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var un struct { // same struct than User String()
				*User
				Timestamp string `json:"timestamp"`
			}
			json.Unmarshal([]byte(tc.exp), &un)

			d, err := time.Parse(time.RFC3339Nano, un.Timestamp)
			if err != nil {
				t.Errorf("cannot parse time %s: %v", tc.exp, err)
				t.FailNow()
			}

			u := &User{
				Address:   un.Address,
				Email:     un.Email,
				UUID:      un.UUID,
				Timestamp: d.UnixMilli(),
				Type:      un.Type,
				Sponsor:   un.Sponsor,
			}

			if u.Timestamp != tc.u.Timestamp {
				t.Errorf("timestamp is incorrect, got %d, want %d", u.Timestamp, tc.u.Timestamp)
				t.FailNow()
			}
		})
	}

}
