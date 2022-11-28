package session

import (
	"context"
)

type Store interface {
	// Save store session and return its updated copy
	Save(ctx context.Context, s *Session) (*Session, error)
	// Save session attributes and return updated copy of session
	AddAttributes(ctx context.Context, sid string, data map[CtxKey]interface{}) (*Session, error)
	// Remove session attributes and return updated copy of session
	RemoveAttributes(ctx context.Context, sid string, keys ...CtxKey) (*Session, error)
	// Load session by its id
	Load(ctx context.Context, sid string) (*Session, error)
	// Invalidate session by its id
	Invalidate(ctx context.Context, sid string) error
}
