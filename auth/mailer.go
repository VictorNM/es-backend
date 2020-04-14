package auth

import (
	"fmt"
	"path"
)

// consoleSender simulate sending email by print to stdin
type consoleSender struct {
	repository ReadRepository
	baseURL    string
}

func (sender *consoleSender) SendActivationEmail(userID int) {
	u, err := sender.repository.FindUserByID(userID)
	if err != nil {
		return
	}

	// TODO: review this, may be replace with some Go libs
	link := path.Join(sender.baseURL, "activate", u.ActivationKey)

	fmt.Printf("Click to %q to activate your account", link)

	return
}

func NewConsoleSender(repository ReadRepository, baseURL string) *consoleSender {
	return &consoleSender{repository: repository, baseURL: baseURL}
}
