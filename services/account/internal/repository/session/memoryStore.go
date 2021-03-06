package session

import (
	"sync"

	"github.com/team-ravl/civic-qa/services/common/model"
)

const intialCapacity = 100

// MemoryStore implements SessionStore,
// Threadsafe via RWMutex
type MemoryStore struct {
	data map[Token]model.SessionState
	lock sync.RWMutex
}

// NewMemoryStore returns a MemorySessionStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[Token]model.SessionState, intialCapacity)}
}

// Create begins a new session for a given SessionState, returning a new SessionToken
func (s *MemoryStore) Create(state model.SessionState) (Token, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// generate a new token as the key for the state
	token, err := generateToken()
	if err != nil {
		return InvalidSessionToken, err
	}
	s.data[token] = state
	return token, nil
}

// Get returns a SessionState for a given SessionToken if it exists
func (s *MemoryStore) Get(token Token) (*model.SessionState, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	state, ok := s.data[token]
	if !ok {
		return nil, ErrStateNotFound
	}

	return &state, nil
}

// Delete removes a SessionState given a SessionToken.
func (s *MemoryStore) Delete(token Token) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// check if the session exists
	// Note: this is not needed, but added for consistency with other implementations
	_, ok := s.data[token]
	if !ok {
		return ErrStateNotFound
	}

	// delete the session
	delete(s.data, token)
	return nil
}
