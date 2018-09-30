package session

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/uuid"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/x/xcontext"
	"github.com/yssk22/go/x/xtime"
	"context"
)

// SessionStore is an interface for session storage.
type SessionStore interface {
	Get(context.Context, uuid.UUID) (*Session, error)
	Set(context.Context, *Session) error
	Del(context.Context, *Session) error
}

// Session is an object to represents session
type Session struct {
	ID         uuid.UUID
	CSRFSecret uuid.UUID
	Timestamp  time.Time
	Data       map[string][]byte
	fromStore  bool
}

// NewSession returns a new *Session
func NewSession() *Session {
	return &Session{
		ID:         uuid.New(),
		CSRFSecret: uuid.New(),
		Timestamp:  xtime.Now(),
		Data:       make(map[string][]byte),
		fromStore:  false,
	}
}

// Set sets the session data
func (s *Session) Set(key interface{}, value interface{}) error {
	buff, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("session data must be json compatible: %v", err)
	}
	s.Data[key.(string)] = buff
	return nil
}

// Get sets the session data
func (s *Session) Get(key interface{}, dst interface{}) error {
	buff, ok := s.Data[key.(string)]
	if !ok {
		return keyvalue.KeyError(key.(string))
	}
	if err := json.Unmarshal(buff, dst); err != nil {
		return err
	}
	return nil
}

// Del delete the session data
func (s *Session) Del(key interface{}) error {
	delete(s.Data, key.(string))
	return nil
}

// IsExpired returns true if the session is expired
func (s *Session) IsExpired(maxAge time.Duration) bool {
	return xtime.Now().After(s.Timestamp.Add(maxAge))
}

// Encode returns a encoded strings of session data, which is passed to session store.
func (s *Session) Encode() (string, error) {
	// marshal to json
	buff, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	// and compress
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(buff)
	w.Close()
	// then base64 encoding for safety.
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// Decode decodes a given `data` into *Session object
func (s *Session) Decode(data string) error {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	r, err := zlib.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	var buff bytes.Buffer
	io.Copy(&buff, r)
	r.Close()
	news := &Session{}
	err = json.Unmarshal([]byte(buff.String()), news)
	if err != nil {
		return err
	}
	s.ID = news.ID
	s.CSRFSecret = news.CSRFSecret
	s.Timestamp = news.Timestamp
	s.Data = news.Data
	return nil
}

// FromContext returns a *Session from a context.
func FromContext(ctx context.Context) *Session {
	s, _ := ctx.Value(contextKey).(*Session)
	return s
}

// FromRequest returns a *Session from a request.
func FromRequest(req *web.Request) *Session {
	return FromContext(req.Context())
}

// NewContext returns a new context with he sesison.
func NewContext(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, contextKey, s)
}

var contextKey = xcontext.NewKey("sessionkey")
