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
	w := httptest.NewRecorder()

	address := "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	email := "john.doe@mailservice.com"
	utype := "talent"

	jsonUser, _ := json.Marshal(User{
		Address: address,
		Email:   email,
		Type:    utype,
	})

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonUser))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusCreated)
		t.FailNow()
	}

	var u User
	err := json.NewDecoder(w.Body).Decode(&u)
	if err != nil {
		t.Errorf("Cannot decode response body %v", w.Body)
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
		t.Errorf("UUID is incorrect, cannot be parsed")
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

func TestValidate(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	token := "12345"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/validate/%s", token), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code is incorrect, got %d, want %d", w.Code, http.StatusOK)
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
}
