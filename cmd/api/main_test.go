package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestRegister(t *testing.T) {
	r := setupRouter()
	tt := []struct {
		name    string
		address string
		email   string
		utype   string
		status  int
		err     string
	}{
		{"valid talent",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"talent",
			http.StatusCreated,
			"",
		},
		{"valid client",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"client",
			http.StatusCreated,
			"",
		},
		{"valid agent",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"agent",
			http.StatusCreated,
			"",
		},
		{"valid mentor",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"mentor",
			http.StatusCreated,
			"",
		},
		{"empty address",
			"",
			"john.doe@mailservice.com",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Address' Error:Field validation for 'Address' failed on the 'required' tag"}`,
		},
		{"0x address",
			"0x",
			"john.doe@mailservice.com",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Address' Error:Field validation for 'Address' failed on the 'eth_addr' tag"}`,
		},
		{"0x0000 address",
			"0x0000",
			"john.doe@mailservice.com",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Address' Error:Field validation for 'Address' failed on the 'eth_addr' tag"}`,
		},
		{"non hexadecimal address",
			"0xYZ25EF3F5B8A186998338A2ADA83795FBA2D695E",
			"john.doe@mailservice.com",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Address' Error:Field validation for 'Address' failed on the 'eth_addr' tag"}`,
		},
		{"too short address",
			"0xDC25EF3F5B8A186998338A2ADA83795FBA2D69",
			"john.doe@mailservice.com",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Address' Error:Field validation for 'Address' failed on the 'eth_addr' tag"}`,
		},
		{"too long address",
			"0xDC25EF3F5B8A186998338A2ADA83795FBA2D695E5E5E5E",
			"john.doe@mailservice.com",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Address' Error:Field validation for 'Address' failed on the 'eth_addr' tag"}`,
		},
		{"empty email",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
		},
		{"malformated email unsupported special characters",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john/doe@email_^me.fr",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
		},
		{"malformated email no @",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe.email.fr",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
		},
		{"malformated email no user",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"@ovh.com",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
		},
		{"malformated email no domain",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@",
			"talent",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
		},
		{"empty type of user",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Type' Error:Field validation for 'Type' failed on the 'required' tag"}`,
		},
		{"unsupported type of user",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"dev",
			http.StatusBadRequest,
			`{"error":"Key: 'User.Type' Error:Field validation for 'Type' failed on the 'oneof' tag"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			address := tc.address
			email := tc.email
			utype := tc.utype

			jsonUser, _ := json.Marshal(User{
				Address: address,
				Email:   email,
				Type:    utype,
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonUser))
			r.ServeHTTP(w, req)

			switch tc.status {
			case http.StatusBadRequest:
				if w.Code != http.StatusBadRequest {
					t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusBadRequest)
					t.FailNow()
				}

				if w.Body == nil {
					t.Errorf("Response body cannot be nil")
					t.FailNow()
				}

				if w.Body.String() != tc.err {
					t.Errorf("Error is incorrect, got %s, want %s", w.Body.String(), tc.err)
					t.FailNow()
				}
			case http.StatusCreated:
				if w.Code != http.StatusCreated {
					t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusCreated)
					t.FailNow()
				}

				m := map[string]interface{}{}
				if err := json.Unmarshal(jsonUser, &m); err != nil {
					t.Errorf("Cannot unmarshal marshaled user: %v", err)
					t.FailNow()
				}

				if len(m) != 6 {
					t.Errorf("Incorrect number of field in json, got %d, want %d", len(m), 6)
					t.FailNow()
				}

				var u User
				err := json.NewDecoder(w.Body).Decode(&u)
				if err != nil {
					t.Errorf("Cannot decode response body %v, %v", w.Body, err)
					t.FailNow()
				}

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

				if u.Validated {
					t.Errorf("Validated is incorrect, got %v, want %v", u.Validated, false)
					t.FailNow()
				}
			default:
				t.Errorf("status %d not supported", tc.status)
				t.FailNow()
			}
		})
	}

}

func TestValidate(t *testing.T) {
	r := setupRouter()

	tt := []struct {
		name   string
		token  string
		status int
		err    string
	}{
		{
			"valid token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG4uZG9lQG1haWxzZXJ2aWNlLmNvbSIsInV1aWQiOiI4MGM3MWVlZS01ZGIyLTRlODctOWIwNC02ODZhYWRkOGQ1N2EifQ.96OWDzW6mL78KjAuVOwa4erGSeMeusWTYFv6Wsnv5-k",
			http.StatusOK,
			"",
		},
		{
			"no token",
			"",
			http.StatusNotFound,
			"",
		},
		{
			"malformated token",
			"eyJhbGciOiJIUzI1N.ZZZZZ.dczrracv.de",
			http.StatusUnauthorized,
			`{"error":"Unauthorized"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			token := tc.token
			req, _ := http.NewRequest("GET", fmt.Sprintf("/validate/%s", token), nil)
			r.ServeHTTP(w, req)

			switch tc.status {
			case http.StatusOK:
				if w.Code != http.StatusOK {
					t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusOK)
					t.Errorf(w.Body.String())
					t.FailNow()
				}

				var res struct {
					Validated bool
					Token     string
				}

				json.NewDecoder(w.Body).Decode(&res)

				if !res.Validated {
					t.Errorf("Validated is incorrect, got %v, want %v", res.Validated, true)
					t.FailNow()
				}

				if res.Token != token {
					t.Errorf("Token is incorrect, got %s, want %s", res.Token, token)
					t.FailNow()
				}
			case http.StatusNotFound:
				if w.Code != http.StatusNotFound {
					t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusNotFound)
					t.Errorf(w.Body.String())
					t.FailNow()
				}
			case http.StatusUnauthorized:
				if w.Code != http.StatusUnauthorized {
					t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusUnauthorized)
					t.Errorf(w.Body.String())
					t.FailNow()
				}

				if w.Body == nil {
					t.Errorf("Response body cannot be nil")
					t.FailNow()
				}

				if w.Body.String() != tc.err {
					t.Errorf("Error is incorrect, got %s, want %s", w.Body.String(), tc.err)
					t.FailNow()
				}
			default:
				t.Errorf("status %d not supported", tc.status)
				t.FailNow()
			}
		})
	}
}
