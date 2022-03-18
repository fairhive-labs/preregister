package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidate(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
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
