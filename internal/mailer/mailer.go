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
	headers = "MIME-Version: 1.0\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
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

func sendEmail(m *SmtpMailer, e, s, n string, data any) (err error) {
	to := []string{e}
	auth := smtp.PlainAuth("", m.from, m.password, m.host)

	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf(`From: %s
To: %s
Subject: %s
%s
`, "no_reply@fairhive-labs.com",
		e,
		s,
		headers)))

	m.t.ExecuteTemplate(&body, n, data)

	fmt.Println("Sending email...")
	r := 3
	for i := 0; i < r; i++ {
		err = smtp.SendMail(m.server, auth, m.from, to, body.Bytes())
		if nil == err {
			break
		}
		fmt.Printf("failed %d/%d, retrying in 1s...\n", i+1, r)
		time.Sleep(1 * time.Second)
	}
	return
}

func (m *SmtpMailer) SendActivationEmail(e, u, h string) (err error) {
	err = sendEmail(m, e, "fairhive - preregistration", "emailActivation",
		struct {
			Hash string
			Url  string
		}{
			Hash: h,
			Url:  u,
		})
	logEmailSent(e, h, err)
	return
}

func (m *MockSmtpMailer) SendActivationEmail(e, u, h string) (err error) {
	// do nothing just log
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
