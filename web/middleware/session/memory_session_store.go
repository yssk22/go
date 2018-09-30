package session

import (
	"fmt"

	"context"

	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/uuid"
)

// MemorySessionStore implements SessionStore on Memory
type MemorySessionStore struct {
	store map[uuid.UUID]*Session
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		store: make(map[uuid.UUID]*Session),
	}
}

// Get implements SessionStore#Get
func (s *MemorySessionStore) Get(ctx context.Context, key uuid.UUID) (*Session, error) {
	session, ok := s.store[key]
	if !ok {
		return nil, keyvalue.KeyError(fmt.Sprintf("seesion:%s", key.String()))
	}
	return session, nil
}

// Set implements SessionStore#Set
func (s *MemorySessionStore) Set(ctx context.Context, session *Session) error {
	s.store[session.ID] = session
	return nil
}

// Del implements SessionStore#Del
func (s *MemorySessionStore) Del(ctx context.Context, session *Session) error {
	delete(s.store, session.ID)
	return nil
}

func (s *MemorySessionStore) String() string {
	return "MemorySessionStore"
}
