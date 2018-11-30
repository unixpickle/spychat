package main

import (
	"sync"
	"time"

	"github.com/unixpickle/fbmsgr"
)

const (
	eventBuffSize         = 30
	sessionExpirationTime = time.Hour * 24
	sessionGracePeriod    = time.Hour
)

// A SessionTable maintains a mapping of unique IDs to
// sessions.
type SessionTable struct {
	lock  sync.RWMutex
	curID int64
	table map[int64]Session
}

// NewSessionTable creates an empty SessionTable.
func NewSessionTable() *SessionTable {
	return &SessionTable{table: map[int64]Session{}}
}

// Add adds a session and returns its ID.
func (s *SessionTable) Add(sess Session) int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.curID++
	s.table[s.curID] = sess
	return s.curID
}

// Get gets a session by its ID.
func (s *SessionTable) Get(i int64) Session {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.table[i]
}

// Del deletes a session.
func (s *SessionTable) Del(i int64) {
	s.lock.Lock()
	delete(s.table, i)
	s.lock.Unlock()
}

// A Session is an abstract source of conversations.
type Session interface {
	Login(username, password string) error
	Chats() ([]*fbmsgr.ThreadInfo, []*fbmsgr.ParticipantInfo, error)
	Thread(id string) ([]fbmsgr.Action, error)
}
