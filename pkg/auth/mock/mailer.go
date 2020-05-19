package mock

type Mailer struct {
	SendFunc func(subject string, tmpl string, data interface{}, to []string) error
}

func (m *Mailer) Send(subject string, tmpl string, data interface{}, to []string) error {
	if m.SendFunc == nil {
		return nil
	}

	return m.SendFunc(subject, tmpl, data, to)
}
