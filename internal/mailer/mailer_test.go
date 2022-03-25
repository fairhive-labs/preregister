package mailer

import (
	"fmt"
	"os"
	"testing"
)

const (
	host     = "smtp.mailtrap.io"
	port     = 2525
	tmplPath = "templates/**"
	email    = "john.doe@domain.com"
	token    = "T0k3n"
	hash     = "hA5h"
)

var (
	from     = os.Getenv("MAILTRAP_USER")
	password = os.Getenv("MAILTRAP_PASSWORD")
)

func TestNewMailer(t *testing.T) {
	mailer := NewMailer(from, password, host, port)
	if mailer.server == "" {
		t.Errorf("incorrect server, got empty string, want %q", fmt.Sprintf("%s:%d", host, port))
		t.FailNow()
	}
	if mailer.t == nil {
		t.Errorf("template cannot be nil")
		t.FailNow()
	}
}

func TestSendActivationEmail(t *testing.T) {
	m := NewMailer(from, password, host, port)
	if err := m.SendActivationEmail(email, fmt.Sprintf("http://fairhive.io/activate/%s", token), hash); err != nil {
		t.Errorf("error sending activation email : %v", err)
		t.FailNow()
	}
}
