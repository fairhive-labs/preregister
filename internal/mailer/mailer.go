package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

const (
	headers = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

type Mailer interface {
	SendActivationEmail(e, u, h string) error
}

type smtpConfig struct {
	from     string
	password string
	host     string
	port     int
	server   string
}

type SmtpMailer struct {
	*smtpConfig
	t *template.Template
}

type MockSmtpMailer struct{}

//go:embed templates/*.html
var tfs embed.FS

func NewMailer(from, password, host string, port int) *SmtpMailer {
	t := template.Must(template.ParseFS(tfs, "templates/*"))
	return &SmtpMailer{&smtpConfig{
		from:     from,
		password: password,
		host:     host,
		port:     port,
		server:   fmt.Sprintf("%s:%d", host, port),
	}, t}
}

func (m *SmtpMailer) SendActivationEmail(e, u, h string) (err error) {
	to := []string{e}
	auth := smtp.PlainAuth("", m.from, m.password, m.host)

	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf("Subject: Complete Preregistration\n%s\n\n", headers)))
	m.t.ExecuteTemplate(&body, "emailActivation", struct {
		Hash string
		Url  string
	}{
		Hash: h,
		Url:  u,
	})

	for i := 0; i < 3; i++ {
		err = smtp.SendMail(m.server, auth, m.from, to, body.Bytes())
		if nil == err {
			break
		}
		time.Sleep(1 * time.Second)
	}
	logEmailSent(e, h, err)
	return
}

func (m *MockSmtpMailer) SendActivationEmail(e, u, h string) (err error) {
	logEmailSent(e, h, err)
	return
}

func logEmailSent(e, h string, err error) {
	if err != nil {
		fmt.Printf("error sending email to %q: %v", e, err)
	} else {
		fmt.Printf("💌 Email to %q: [ \033[1;32mSent\033[0m ]\n🧬 Hash: %s\n", e, h)
	}
}
