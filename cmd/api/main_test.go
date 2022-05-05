package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fairhive-labs/preregister/internal/crypto"
	"github.com/fairhive-labs/preregister/internal/crypto/cipher"
	"github.com/fairhive-labs/preregister/internal/data"
	"github.com/fairhive-labs/preregister/internal/limiter"
	"github.com/fairhive-labs/preregister/internal/mailer"
)

func TestRegister(t *testing.T) {
	var db data.DB = data.MockDB
	k, _ := cipher.GenerateKey(32)
	app := &App{
		db,
		crypto.NewJWTHS256(k),
		&mailer.MockSmtpMailer,
		sync.WaitGroup{},
		limiter.NewUnlimited(),
		"path1",
		"path2",
	}
	r := setupRouter(*app)
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
	var db data.DB = data.MockDB
	k, _ := cipher.GenerateKey(32)
	app := &App{
		db,
		crypto.NewJWTHS256(k),
		&mailer.MockSmtpMailer,
		sync.WaitGroup{},
		limiter.NewUnlimited(),
		"path1",
		"path2",
	}
	r := setupRouter(*app)

	address, email, utype := "0x8ba1f109551bD432803012645Ac136ddd64DBA72", "john.doe@mailservice.com", "talent"
	vt, _ := app.jwt.Create(&data.User{
		Address: address,
		Email:   email,
		Type:    utype}, time.Now())
	vh := app.jwt.Hash(vt)

	tt := []struct {
		name        string
		token, hash string
		status      int
		err         string
	}{
		{
			"valid token+hash",
			vt, vh,
			http.StatusCreated,
			"",
		},
		{
			"valid token invalid hash",
			vt, "hA5h",
			http.StatusUnauthorized,
			`{"error":"Unauthorized"}`,
		},
		{
			"no token no hash",
			"", "",
			http.StatusTemporaryRedirect,
			"",
		},
		{
			"malformated token",
			"eyJhbGciOiJIUzI1N.ZZZZZ.dczrracv.de", "hA5h",
			http.StatusUnauthorized,
			`{"error":"Unauthorized"}`,
		},
		{
			"malformated token no hash",
			"eyJhbGciOiJIUzI1N.ZZZZZ.dczrracv.de", "",
			http.StatusTemporaryRedirect,
			"",
		},
		{
			"malformated token fake hash",
			"eyJhbGciOiJIUzI1N.ZZZZZ.dczrracv.de", "hA5h",
			http.StatusUnauthorized,
			`{"error":"Unauthorized"}`,
		},
		{
			"fake token valid hash",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiIiwiZW1haWwiOiJqb2huLmRvZUBtYWlsc2VydmljZS5jb20iLCJ0eXBlIjoidGFsZW50IiwiaXNzIjoiZmFpcmhpdmUuaW8iLCJleHAiOjE2NDg1NTIsIm5iZiI6MTY0Nzk1MiwiaWF0IjoxNjQ3OTUyfQ.fU_Vi8s1vxU59oRSAnTGj3bN4veMeNuFYBLNWKuOnJE",
			"69E045328DCE6EAD2F52068EF4FB2232EBF12E82FF4712947B2DBC3CCEA8035000CF037B22EE1DD6B0C938AA9B7D329CC6DD111824F01D494FDD75B1150BD628",
			http.StatusUnauthorized,
			`{"error":"Unauthorized"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			token := tc.token
			hash := tc.hash
			req, _ := http.NewRequest("POST", fmt.Sprintf("/activate/%s/%s", token, hash), nil)
			r.ServeHTTP(w, req)

			switch tc.status {
			case http.StatusCreated:
				if w.Code != http.StatusCreated {
					t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusCreated)
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
			case http.StatusTemporaryRedirect:
				if w.Code != http.StatusTemporaryRedirect {
					t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusTemporaryRedirect)
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

	app.db = data.MockErrDB
	r = setupRouter(*app)
	t.Run("faulty DB", func(t *testing.T) {
		w := httptest.NewRecorder()
		token := vt
		hash := vh
		req, _ := http.NewRequest("POST", fmt.Sprintf("/activate/%s/%s", token, hash), nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusInternalServerError)
			t.FailNow()
		}
	})
}

func TestCount(t *testing.T) {
	var db data.DB = data.MockDB
	k, _ := cipher.GenerateKey(32)
	app := &App{
		db,
		crypto.NewJWTHS256(k),
		&mailer.MockSmtpMailer,
		sync.WaitGroup{},
		limiter.NewUnlimited(),
		"path1",
		"path2",
	}
	r := setupRouter(*app)

	t.Run("json normal", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/count?mime=json", app.secpath1, app.secpath2), nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("incorrect status, got %d, want %d", w.Code, http.StatusOK)
			t.FailNow()
		}
		if w.Body == nil {
			t.Errorf("Response body cannot be nil")
			t.FailNow()
		}

		var res struct {
			Users map[string]int
			Total int
		}
		err := json.NewDecoder(w.Body).Decode(&res)
		if err != nil {
			t.Errorf("Cannot decode response body %v, %v", w.Body, err)
			t.FailNow()
		}

		for ut, uc := range data.UsersMapMock {
			if res.Users[ut] != uc {
				t.Errorf("incorrect %q count, got %d, want %d", ut, res.Users[ut], uc)
				t.FailNow()
			}
		}

		if res.Total != data.UsersCountMock {
			t.Errorf("incorrect total count, got %d, want %d", res.Total, data.UsersCountMock)
			t.FailNow()
		}
	})

	t.Run("xml normal", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/count?mime=xml", app.secpath1, app.secpath2), nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("incorrect status, got %d, want %d", w.Code, http.StatusOK)
			t.FailNow()
		}
		if w.Body == nil {
			t.Errorf("Response body cannot be nil")
			t.FailNow()
		}

		type xmlUser struct {
			Type  string
			Value int
		}
		type Count struct {
			Total int
			Users []xmlUser
		}
		var res Count
		err := xml.NewDecoder(w.Body).Decode(&res)
		if err != nil {
			t.Errorf("Cannot decode response body %v, %v", w.Body, err)
			t.FailNow()
		}
		if res.Total != data.UsersCountMock {
			t.Errorf("incorrect total count, got %d, want %d", res.Total, data.UsersCountMock)
			t.FailNow()
		}
		for _, xu := range res.Users {
			if xu.Value != data.UsersMapMock[xu.Type] {
				t.Errorf("incorrect agent count, got %d, want %d", xu.Value, data.UsersMapMock[xu.Type])
				t.FailNow()
			}
		}
	})

	t.Run("html normal", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/count", app.secpath1, app.secpath2), nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("incorrect status, got %d, want %d", w.Code, http.StatusOK)
			t.FailNow()
		}
		if w.Body == nil {
			t.Errorf("Response body cannot be nil")
			t.FailNow()
		}

		users := data.UsersMapMock
		m := map[string]bool{
			fmt.Sprintf(`<td colspan="2">Total: %d</td>`, data.UsersCountMock): false,
		}
		for t, v := range users {
			k := fmt.Sprintf(`%s</td><td class="count">%d</td>`, t, v)
			m[k] = false
		}
		for l, err := w.Body.ReadString('\n'); err == nil; {
			for k := range m {
				if strings.Contains(l, k) {
					m[k] = true
				}
			}
			l, err = w.Body.ReadString('\n')
		}
		for k, v := range m {
			if !v {
				t.Errorf("response body should contain %q", k)
				t.FailNow()
			}
		}
	})

	tt := []struct {
		name         string
		path1, path2 string
		status       int
		body         string
	}{
		{"fakepath1", "fakepath1", app.secpath2, http.StatusNotFound, ""},
		{"fakepath2", app.secpath1, "fakepath2", http.StatusNotFound, ""},
		{"fakepaths", "fakepath1", "fakepath2", http.StatusNotFound, ""},
		{"missing path1", "", "fakepath2", http.StatusNotFound, "404 page not found"},
		{"missing path2", "fakepath1", "", http.StatusNotFound, ""},
		{"no path", "", "", http.StatusNotFound, ""},
	}
	for _, tc := range tt {
		t.Run("json_"+tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/count?mime=json", tc.path1, tc.path2), nil)
			r.ServeHTTP(w, req)
			if w.Code != tc.status {
				t.Errorf("incorrect status, got %d, want %d", w.Code, tc.status)
				t.FailNow()
			}
			if w.Body.String() != tc.body {
				t.Errorf("incorrect body, got %q, want %q", w.Body.String(), tc.body)
				t.FailNow()
			}
		})
	}

	app.db = data.MockErrDB
	r = setupRouter(*app)
	t.Run("json faulty DB", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/count?mime=json", app.secpath1, app.secpath2), nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusInternalServerError)
			t.FailNow()
		}
	})
}

func TestHealth(t *testing.T) {
	var db data.DB = data.MockDB
	k, _ := cipher.GenerateKey(32)
	app := &App{
		db,
		crypto.NewJWTHS256(k),
		&mailer.MockSmtpMailer,
		sync.WaitGroup{},
		limiter.NewUnlimited(),
		"path1",
		"path2",
	}
	r := setupRouter(*app)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("incorrect status code, got %d, want %d\n", w.Code, http.StatusOK)
		t.FailNow()
	}

	headers := w.Result().Header
	if headers.Get("Access-Control-Allow-Methods") == "" {
		t.Errorf("Access-Control-Allow-Methods header cannot be empty, must be %q\n", "POST, OPTIONS")
		t.FailNow()
	}
	if headers.Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("incorrect Content-Type, got %q, want %q\n", headers.Get("Content-Type"), "text/plain; charset=utf-8")
		t.FailNow()
	}

	if w.Body.Len() == 0 {
		t.Error("Body cannot be empty")
		t.FailNow()
	}

	if w.Body.String() != "ok" {
		t.Errorf("incorrect body: must only contain %q", "ok")
		t.FailNow()
	}
}

func TestList(t *testing.T) {
	var db data.DB = data.MockDB
	k, _ := cipher.GenerateKey(32)
	app := &App{
		db,
		crypto.NewJWTHS256(k),
		&mailer.MockSmtpMailer,
		sync.WaitGroup{},
		limiter.NewUnlimited(),
		"path1",
		"path2",
	}
	r := setupRouter(*app)

	t.Run("json normal", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/list", app.secpath1, app.secpath2), nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("incorrect status, got %d, want %d", w.Code, http.StatusOK)
			t.FailNow()
		}
		if w.Body == nil {
			t.Errorf("Response body cannot be nil")
			t.FailNow()
		}

		var res struct {
			Users []*data.User
			Count int
		}
		err := json.NewDecoder(w.Body).Decode(&res)
		if err != nil {
			t.Errorf("Cannot decode response body %v, %v", w.Body, err)
			t.FailNow()
		}

		count := data.UsersCountMock
		if res.Count != count {
			t.Errorf("incorrect count, got %d, want %d", res.Count, count)
			t.FailNow()
		}

		if len(res.Users) != res.Count {
			t.Errorf("incorrect length of Users array, got %d, want %d", len(res.Users), res.Count)
			t.FailNow()
		}

	})

	tt := []struct {
		name         string
		path1, path2 string
		status       int
		body         string
	}{
		{"fakepath1", "fakepath1", app.secpath2, http.StatusNotFound, ""},
		{"fakepath2", app.secpath1, "fakepath2", http.StatusNotFound, ""},
		{"fakepaths", "fakepath1", "fakepath2", http.StatusNotFound, ""},
		{"missing path1", "", "fakepath2", http.StatusNotFound, "404 page not found"},
		{"missing path2", "fakepath1", "", http.StatusNotFound, ""},
		{"no path", "", "", http.StatusNotFound, ""},
	}
	for _, tc := range tt {
		t.Run("json_"+tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/list", tc.path1, tc.path2), nil)
			r.ServeHTTP(w, req)
			if w.Code != tc.status {
				t.Errorf("incorrect status, got %d, want %d", w.Code, tc.status)
				t.FailNow()
			}
			if w.Body.String() != tc.body {
				t.Errorf("incorrect body, got %q, want %q", w.Body.String(), tc.body)
				t.FailNow()
			}
		})
	}

	app.db = data.MockErrDB
	r = setupRouter(*app)
	t.Run("json faulty DB", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/list", app.secpath1, app.secpath2), nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusInternalServerError)
			t.FailNow()
		}
	})
}
