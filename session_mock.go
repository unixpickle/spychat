package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/fbmsgr"
)

// NewMockSession creates a Session that uses a mocked
// version of Facebook Messenger.
func NewMockSession(configPath string) Session {
	data, err := ioutil.ReadFile(configPath)
	essentials.Must(err)
	var res mockSession
	essentials.Must(essentials.AddCtx("parse mock session", json.Unmarshal(data, &res)))
	return &res
}

type mockSession struct {
	ThreadsResult []*fbmsgr.ThreadInfo    `json:"threads"`
	ThreadResult  []*fbmsgr.GenericAction `json:"thread"`
}

func (m *mockSession) Login(username, password string) error {
	if username == "username" && password == "password" {
		return nil
	}
	return errors.New("login incorrect")
}

func (m *mockSession) Threads() ([]*fbmsgr.ThreadInfo, error) {
	time.Sleep(time.Second)
	return m.ThreadsResult, nil
}

func (m *mockSession) Thread(id string) ([]fbmsgr.Action, error) {
	time.Sleep(time.Second)
	res := make([]fbmsgr.Action, len(m.ThreadResult))
	for i, x := range m.ThreadResult {
		res[i] = x
	}
	return res, nil
}
