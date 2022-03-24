package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fairhive-labs/preregister/internal/data"
)

func TestRegister(t *testing.T) {
	app := *NewApp(data.MockDB)
	r := setupRouter(app)
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
			http.StatusAccepted,
			"",
		},
		{"valid client",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"client",
			http.StatusAccepted,
			"",
		},
		{"valid agent",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"agent",
			http.StatusAccepted,
			"",
		},
		{"valid mentor",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"mentor",
			http.StatusAccepted,
			"",
		},
		{"valid advisor",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"advisor",
			http.StatusAccepted,
			"",
		},
		{"valid contributor",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"contributor",
			http.StatusAccepted,
			"",
		},
		{"valid investor",
			"0x8ba1f109551bD432803012645Ac136ddd64DBA72",
			"john.doe@mailservice.com",
			"investor",
			http.StatusAccepted,
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

			jsonUser, _ := json.Marshal(data.User{
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
			case http.StatusAccepted:
				if w.Code != http.StatusAccepted {
					t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusAccepted)
					t.FailNow()
				}

				var res struct {
					Hash string
				}
				err := json.NewDecoder(w.Body).Decode(&res)
				if err != nil {
					t.Errorf("Cannot decode response body %v, %v", w.Body, err)
					t.FailNow()
				}
				if res.Hash == "" {
					t.Errorf("incorrect hash, cannot be empty string")
				}

			default:
				t.Errorf("status %d not supported", tc.status)
				t.FailNow()
			}
		})
	}

}

func TestActivate(t *testing.T) {
	app := *NewApp(data.MockDB)
	r := setupRouter(app)

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
			req, _ := http.NewRequest("POST", fmt.Sprintf("/activate/%s", token), nil)
			r.ServeHTTP(w, req)

			switch tc.status {
			case http.StatusOK:
				if w.Code != http.StatusOK {
					t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusOK)
					t.Errorf(w.Body.String())
					t.FailNow()
				}

				var res struct {
					Activated bool
					Token     string
				}

				json.NewDecoder(w.Body).Decode(&res)

				if !res.Activated {
					t.Errorf("Activated is incorrect, got %v, want %v", res.Activated, true)
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
