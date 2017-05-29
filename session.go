package spychat

import (
	"io"
	"sync"
	"time"

	"github.com/unixpickle/fbmsgr"
)

const (
	eventBuffSize         = 30
	sessionExpirationTime = time.Hour * 24
	sessionGracePeriod    = time.Hour
)

// A Session manages a session on Facebook Messenger and
// multiplexes the event stream for multiple clients.
type Session struct {
	username string
	password string

	sessLock sync.Mutex
	sess     *fbmsgr.Session

	listenersLock sync.RWMutex
	listeners     []*Listener
	isListening   bool
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
//
// You should not manually call ReadEvent() on the
// session, as that is handled by Listen().
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

// Listen adds a Listener for events.
func (s *Session) Listen() *Listener {
	s.listenersLock.Lock()
	defer s.listenersLock.Unlock()
	stream := make(chan fbmsgr.Event, eventBuffSize)
	listener := &Listener{
		stream: stream,
	}
	s.listeners = append(s.listeners, listener)
	if !s.isListening {
		s.isListening = true
		go s.listenerLoop()
	}
	return listener
}

// Unlisten removes a Listener.
//
// This will not close the Listener's channel.
func (s *Session) Unlisten(l *Listener) {
	s.listenersLock.Lock()
	defer s.listenersLock.Unlock()
	for i, x := range s.listeners {
		if x == l {
			s.listeners[i] = s.listeners[len(s.listeners)-1]
			s.listeners = s.listeners[:len(s.listeners)-1]
			return
		}
	}
}

func (s *Session) listenerLoop() {
	for {
		// Stop the loop if nobody is listening anymore.
		s.listenersLock.Lock()
		if len(s.listeners) == 0 {
			s.isListening = false
			s.listenersLock.Unlock()
			return
		}
		s.listenersLock.Unlock()

		sess, err := s.Session()
		if err != nil {
			s.listenerError(err)
			return
		}
		event, err := sess.ReadEvent()
		if err == io.EOF {
			continue
		} else if err != nil {
			s.listenerError(err)
			return
		}
		s.listenersLock.Lock()
		if len(s.listeners) == 0 {
			s.isListening = false
			s.listenersLock.Unlock()
			return
		}
		for _, listener := range s.listeners {
			// Drop events if the listeners can't keep up.
			select {
			case listener.stream <- event:
			default:
			}
		}
		s.listenersLock.Unlock()
	}
}

func (s *Session) listenerError(e error) {
	s.listenersLock.Lock()
	defer s.listenersLock.Unlock()
	for _, listener := range s.listeners {
		close(listener.stream)
	}
	s.listeners = nil
	s.isListening = false
}

// A Listener provides an event stream from a Session.
type Listener struct {
	stream chan fbmsgr.Event
}

// Chan returns the channel of events.
//
// The channel is closed if the event stream ends
// unexpectedly.
// This might happen in the case where the user changes
// their Facebook password.
func (l *Listener) Chan() <-chan fbmsgr.Event {
	return l.stream
}
