package main

import (
	"testing"
)

func TestSetup(t *testing.T) {
	tn, k, p1, p2 := "Waitlist_UnitTest", "Sup3rSecr3tKAY", "p4th1", "p4th2"
	t.Setenv("FAIRHIVE_PREREGISTER_TABLE_NAME", tn)
	t.Setenv("FAIRHIVE_ENCRYPTION_KEY", k)
	t.Setenv("FAIRHIVE_API_SECURE_PATH1", p1)
	t.Setenv("FAIRHIVE_API_SECURE_PATH2", p2)

	setup()
	if tableName != tn {
		t.Errorf("wrong table name, got %s, want %s", tableName, tn)
		t.FailNow()
	}
	if ek != k {
		t.Errorf("wrong table name, got %s, want %s", ek, k)
		t.FailNow()
	}
	if secpath1 != p1 {
		t.Errorf("wrong secure path #1, got %s, want %s", secpath1, p1)
		t.FailNow()
	}
	if secpath2 != p2 {
		t.Errorf("wrong secure path #2, got %s, want %s", secpath2, p2)
		t.FailNow()
	}

}

func TestNewApp(t *testing.T) {
	tn, k, p1, p2 := "Waitlist_UnitTest", "Sup3rSecr3tKAY", "p4th1", "p4th2"
	t.Setenv("FAIRHIVE_PREREGISTER_TABLE_NAME", tn)
	t.Setenv("FAIRHIVE_ENCRYPTION_KEY", k)
	t.Setenv("FAIRHIVE_API_SECURE_PATH1", p1)
	t.Setenv("FAIRHIVE_API_SECURE_PATH2", p2)
	setup()
	app := newApp()
	if app == nil {
		t.Errorf("app cannot be nil")
		t.FailNow()
	}
}
