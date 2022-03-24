package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

const (
	headers = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

type smtpConfig struct {
	from     string
	password string
	host     string
	port     int
	server   string
}

type Mailer struct {
	*smtpConfig
	t *template.Template
}

func NewMailer(from, password, host string, port int, tmplPath string) *Mailer {
	t := template.Must(template.ParseGlob(tmplPath)) //@TODO : use go-embed
	return &Mailer{&smtpConfig{
		from:     from,
		password: password,
		host:     host,
		port:     port,
		server:   fmt.Sprintf("%s:%d", host, port),
	}, t}
}

func (m *Mailer) SendActivationEmail(e, u, h string) {
	to := []string{e}
	fmt.Println("PlainAuth:", m.from, m.password, m.host)
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

	err := smtp.SendMail(m.server, auth, m.from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("💌 Email to %q: [ \033[1;32mSent\033[0m ]\n🧬 Hash: %s\n", e, h)
	}
}
