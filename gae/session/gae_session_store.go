package session

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine"

	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web/middleware/session"
)

const LoggerKey = session.SessionLoggerKey

// GAESessionStore implements SessionStore on GAE memcache and datastore
type GAESessionStore struct {
	namespace string
}

// NewGAESessionStore returns a new *SessionStore
func NewGAESessionStore(namespace string) *GAESessionStore {
	return &GAESessionStore{
		namespace: namespace,
	}
}

// Get implements SessionStore#Get
func (s *GAESessionStore) Get(ctx context.Context, key uuid.UUID) (*session.Session, error) {
	_ctx, err := appengine.Namespace(ctx, s.namespace)
	if err != nil {
		return nil, err
	}
	wrapper := DefaultSessionKind.MustGet(_ctx, key.String())
	if wrapper == nil {
		return nil, keyvalue.KeyError(fmt.Sprintf("seesion:%s", key.String()))
	}
	sess := &session.Session{}
	sess.ID, _ = uuid.FromString(wrapper.ID)
	sess.CSRFSecret, _ = uuid.FromString(wrapper.CSRFSecret)
	sess.Timestamp = wrapper.Timestamp
	if err := json.Unmarshal(wrapper.Data, &sess.Data); err != nil {
		return nil, err
	}
	return sess, nil
}

// Set implements SessionStore#Set
func (s *GAESessionStore) Set(ctx context.Context, session *session.Session) error {
	_ctx, err := appengine.Namespace(ctx, s.namespace)
	if err != nil {
		return err
	}
	data, err := json.Marshal(session.Data)
	if err != nil {
		return fmt.Errorf("could not marshal session.Data: %v", err)
	}
	wrapper := &Session{
		ID:         session.ID.String(),
		CSRFSecret: session.CSRFSecret.String(),
		Timestamp:  session.Timestamp,
		Data:       data,
	}
	DefaultSessionKind.MustPut(_ctx, wrapper)
	return nil
}

// Del implements SessionStore#Del
func (s *GAESessionStore) Del(ctx context.Context, session *session.Session) error {
	_ctx, err := appengine.Namespace(ctx, s.namespace)
	if err != nil {
		return err
	}
	DefaultSessionKind.MustDelete(_ctx, session.ID.String())
	return nil
}

func (s *GAESessionStore) String() string {
	return fmt.Sprintf("GAESessionStore(ns=%s)", s.namespace)
}
