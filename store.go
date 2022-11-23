package session

import (
	"context"
)

type Store interface {
	Save(ctx context.Context, s *Session) (*Session, error)
	Update(ctx context.Context, s *Session) (*Session, error)
	Load(ctx context.Context, sid string) (*Session, error)
	Invalidate(ctx context.Context, sid string) error
}
