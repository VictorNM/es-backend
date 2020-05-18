package auth

import (
	"html/template"
	"path"
)

type Mailer interface {
	Send(subject string, tmpl string, data interface{}, to []string) error
}

type activationEmailSender struct {
	mailer     Mailer
	repository UserRepository
	path       string
}

func (sender *activationEmailSender) SendActivationEmail(userID int) {
	u, err := sender.repository.FindUserByID(userID)
	if err != nil {
		return
	}

	link := path.Join(sender.path, u.ActivationKey)

	const tpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Activate for email</title>
	</head>
	<body>
		<a href="{{ .link }}">Click here to activate your email</a>
	</body>
</html>`

	_ = sender.mailer.Send("Activate your email!", tpl, map[string]interface{}{"link": template.URL(link)}, []string{u.Email})
}
