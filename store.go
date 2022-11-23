package session

import (
	"context"
)

type Store interface {
	// Save store session and return its updated copy
	Save(ctx context.Context, s *Session) (*Session, error)
	// Update session and return its updated copy
	Update(ctx context.Context, s *Session) (*Session, error)
	// Load session by its id
	Load(ctx context.Context, sid string) (*Session, error)
	// Invalidate session by its id
	Invalidate(ctx context.Context, sid string) error
}
