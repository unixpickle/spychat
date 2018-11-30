package main

import (
	"errors"

	"github.com/unixpickle/fbmsgr"
)

// NewMockSession creates a Session that uses a mocked
// version of Facebook Messenger.
func NewMockSession() Session {
	return &mockSession{}
}

type mockSession struct {
	username string
	password string
}

func (m *mockSession) Login(username, password string) error {
	if username == "username" && password == "password" {
		m.username = username
		m.password = password
		return nil
	}
	return errors.New("login incorrect")
}

func (m *mockSession) Chats() ([]*fbmsgr.ThreadInfo, []*fbmsgr.ParticipantInfo, error) {
	// TODO: this.
	panic("nyi")
}

func (m *mockSession) Thread(id string) ([]fbmsgr.Action, error) {
	// TODO: this.
	panic("nyi")
}
