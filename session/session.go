package session

import (
	"sync"
	"time"

	"github.com/dresbach/project/statemachine"
)

// Session represents a user's session.
type Session struct {
	State     statemachine.State
	LastSeen  time.Time
	Variables map[string]interface{}
}

// Store is a thread-safe in-memory session store.
type Store struct {
	mutex    sync.RWMutex
	sessions map[string]*Session
	timeout  time.Duration
}

// NewStore creates a new session store.
func NewStore(timeout time.Duration) *Store {
	return &Store{
		sessions: make(map[string]*Session),
		timeout:  timeout,
	}
}

// Get returns the session for a given user ID.
func (s *Store) Get(userID string) *Session {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, ok := s.sessions[userID]
	if !ok || time.Since(session.LastSeen) > s.timeout {
		return &Session{
			State:     statemachine.StateRoot,
			LastSeen:  time.Now(),
			Variables: make(map[string]interface{}),
		}
	}
	session.LastSeen = time.Now()
	return session
}

// Set updates the session for a given user ID.
func (s *Store) Set(userID string, session *Session) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.sessions[userID] = session
}
