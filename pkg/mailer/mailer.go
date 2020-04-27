package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type Mailer struct {
	host              string
	port              int
	account, password string
	from              string
}

func NewMailer(host string, port int, account, password string, from string) *Mailer {
	return &Mailer{
		host:     host,
		port:     port,
		account:  account,
		password: password,
		from:     from,
	}
}

func (m *Mailer) Send(tmpl string, data interface{}, to []string) error {
	a := smtp.PlainAuth("", m.account, m.password, m.host)

	msg, err := parseTemplate(tmpl, data)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	err = smtp.SendMail(addr, a, m.from, to, msg)
	if err != nil {
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
