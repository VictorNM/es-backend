package auth

import (
	"path"
)

type sender interface {
	Send(tmpl string, data interface{}, to []string) error
}

type ActivationEmailSender interface {
	SendActivationEmail(userID int)
}

type activationEmailSender struct {
	s          sender
	repository ReadUserRepository
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
		<div>{{ .link }}</div>
	</body>
</html>`

	_ = sender.s.Send(tpl, map[string]string{"link": link}, []string{u.Email})
}

func NewActivationEmailSender(repository ReadUserRepository, path string) *activationEmailSender {
	return &activationEmailSender{repository: repository, path: path}
}
