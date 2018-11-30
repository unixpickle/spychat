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
	table map[int64]*Session
}

// NewSessionTable creates an empty SessionTable.
func NewSessionTable() *SessionTable {
	return &SessionTable{table: map[int64]*Session{}}
}

// Add adds a session and returns its ID.
func (s *SessionTable) Add(sess *Session) int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.curID++
	s.table[s.curID] = sess
	return s.curID
}

// Get gets a session by its ID.
func (s *SessionTable) Get(i int64) *Session {
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

// A Session manages a session on Facebook Messenger and
// multiplexes the event stream for multiple clients.
type Session struct {
	username string
	password string

	sessLock sync.Mutex
	sess     *fbmsgr.Session
}

// NewSession authenticates a new session.
func NewSession(user, password string) (*Session, error) {
	res := &Session{
		username: user,
		password: password,
	}
	if _, err := res.Session(); err != nil {
		return nil, err
	}
	return res, nil
}

// Session returns the current underlying *fbmsgr.Session
// for interacting with Messenger.
// This may change between consecutive calls.
func (s *Session) Session() (*fbmsgr.Session, error) {
	s.sessLock.Lock()
	defer s.sessLock.Unlock()
	if s.sess != nil {
		return s.sess, nil
	}
	sess, err := fbmsgr.Auth(s.username, s.password)
	if err != nil {
		return nil, err
	}
	s.sess = sess

	// After a decent amount of time, we destroy the
	// session so that we don't have any permanently
	// dangling resources.
	go func() {
		time.Sleep(sessionExpirationTime)
		s.sessLock.Lock()
		if s.sess == sess {
			s.sess = nil
		}
		s.sessLock.Unlock()
		time.Sleep(sessionGracePeriod)
		sess.Close()
	}()

	return sess, nil
}
