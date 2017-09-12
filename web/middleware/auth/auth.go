package auth

import "github.com/speedland/go/web/middleware/session"

type AuthProvider interface {
	Get(*session.Session) (interface{}, error)
	Set(*session.Session, interface{}) error
}
