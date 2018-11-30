package main

import (
	"sync"
	"time"

	"github.com/unixpickle/fbmsgr"
)

// NewRealSession creates a Session that actually uses
// Facebook Messenger.
func NewRealSession() Session {
	return &realSession{}
}

// A realSession manages an actual Messenger session.
type realSession struct {
	username string
	password string

	sessLock sync.Mutex
	sess     *fbmsgr.Session
}

func (r *realSession) Login(username, password string) error {
	r.username = username
	r.password = password
	_, err := r.rawSession()
	return err
}

func (r *realSession) Threads() ([]*fbmsgr.ThreadInfo, error) {
	sess, err := r.rawSession()
	if err != nil {
		return nil, err
	}
	return sess.AllThreads()
}

func (r *realSession) Thread(id string) ([]fbmsgr.Action, error) {
	sess, err := r.rawSession()
	if err != nil {
		return nil, err
	}
	return sess.ActionLog(id, time.Time{}, 100)
}

// rawSession returns the current underlying session for
// interacting with Messenger.
// This may change between consecutive calls.
func (r *realSession) rawSession() (*fbmsgr.Session, error) {
	r.sessLock.Lock()
	defer r.sessLock.Unlock()
	if r.sess != nil {
		return r.sess, nil
	}
	sess, err := fbmsgr.Auth(r.username, r.password)
	if err != nil {
		return nil, err
	}
	r.sess = sess

	// After a decent amount of time, we destroy the
	// session so that we don't have any permanently
	// dangling resources.
	go func() {
		time.Sleep(sessionExpirationTime)
		r.sessLock.Lock()
		if r.sess == sess {
			r.sess = nil
		}
		r.sessLock.Unlock()
		time.Sleep(sessionGracePeriod)
		sess.Close()
	}()

	return sess, nil
}
