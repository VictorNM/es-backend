package mailer

import (
	"bytes"
	"html/template"
	"net/smtp"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	auth   smtp.Auth
	dialer *gomail.Dialer

	host              string
	port              int
	account, password string
	from              string
}

func New(host string, port int, account, password string, from string) *Mailer {
	return &Mailer{
		auth:   smtp.PlainAuth("", account, password, host),
		dialer: gomail.NewDialer(host, port, account, password),

		host:     host,
		port:     port,
		account:  account,
		password: password,
		from:     from,
	}
}

func (m *Mailer) Send(subject string, tmpl string, data interface{}, to []string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	body, err := parseTemplate(tmpl, data)
	if err != nil {
		return err
	}
	msg.SetBody("text/html", string(body))

	if err := m.dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func parseTemplate(tmpl string, data interface{}) ([]byte, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
