package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()

	address := "0x8ba1f109551bD432803012645Ac136ddd64DBA72"
	email := "john.doe@mailservice.com"
	utype := "talent"

	u := User{
		Address: address,
		Email:   email,
		Type:    utype,
	}

	jsonUser, _ := json.Marshal(u)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonUser))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("code = %d, exp : %d", w.Code, http.StatusNotImplemented)
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
		t.Errorf("code = %d, exp : %d", w.Code, http.StatusOK)
		t.FailNow()
	}

	var res struct {
		Validated bool
		Token     string
	}

	json.NewDecoder(w.Body).Decode(&res)

	if !res.Validated {
		t.Errorf("validated = %v, exp : %v", res.Validated, true)
		t.FailNow()
	}

	if res.Token != token {
		t.Errorf("res.Token = %s, exp : %s", res.Token, token)
		t.FailNow()
	}
}
