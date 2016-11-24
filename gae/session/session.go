// Package session provides github.com/speedland/go/web/middleware/session.SessionStore
// implementation for GAE environment
package session

import "time"

// Session is a wrapped struct for github.com/speedland/go/web/middleware/sesison.Session
//go:generate ent -type=Session
type Session struct {
	ID         string    `json:"id" ent:"id"`
	CSRFSecret string    `json:"csrf_secret" datastore:",noindex"`
	Timestamp  time.Time `json:"key" ent:"timestamp" datastore:",noindex"`
	Data       []byte    `json:"data" datastore:",noindex"`
}
