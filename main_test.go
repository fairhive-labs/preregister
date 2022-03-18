package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	}

	json.NewDecoder(w.Body).Decode(&res)

	if !res.Validated {
		t.Errorf("validated = %v, exp : %v", res.Validated, true)
	}
}

func TestRegister(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("code = %d, exp : %d", w.Code, http.StatusNotImplemented)
		t.FailNow()
	}
}
