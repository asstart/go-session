package session

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"time"
)

type Session struct {
	ID   string
	Data map[SessionKey]interface{}
	Opts CookieConf

	Anonym bool
	Active bool

	UID string

	IdleTimeout time.Duration
	AbsTimeout  time.Duration

	LastAccessedAt time.Time
	CreatedAt      time.Time
}

type SameSite int

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

type CookieConf struct {
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	MaxAge   int
	SameSite SameSite
}

type SessionConf struct {
	IdleTimeout time.Duration
	AbsTimout   time.Duration
}

type SessionKey string

var (
	base32enc = base32.StdEncoding.WithPadding(base32.NoPadding)
	keyLen    = 32
)

func DefaultCookieConf() CookieConf {
	return CookieConf{
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		SameSite: SameSiteStrictMode,
	}
}

func DefaultSessionConf() SessionConf {
	return SessionConf{
		AbsTimout:   24 * 7 * time.Hour,
		IdleTimeout: 24 * time.Hour,
	}
}

func NewSession() (Session, error) {
	id, err := generateSessionId()
	if err != nil {
		return Session{}, err
	}
	s := Session{
		ID:          id,
		Data:        make(map[SessionKey]interface{}),
		Opts:        DefaultCookieConf(),
		Anonym:      true,
		Active:      true,
		IdleTimeout: DefaultSessionConf().IdleTimeout,
		AbsTimeout:  DefaultSessionConf().AbsTimout,
	}
	return s, nil
}

func (s *Session) WithUserId(uid string) {
	s.UID = uid
	s.Anonym = false
}

func (s *Session) WithCookieConf(cc CookieConf) {
	s.Opts = cc
}

func (s *Session) WithSessionConf(sc SessionConf) {
	s.IdleTimeout = sc.IdleTimeout
	s.AbsTimeout = sc.AbsTimout
}

func (s *Session) AddAttribute(k SessionKey, v interface{}) {
	s.Data[k] = v
}

func (s *Session) GetAttribute(k SessionKey) (interface{}, bool) {
	v, ok := s.Data[k]
	return v, ok
}

func (s *Session) IsExpired() bool {
	if !s.Active {
		return true
	}

	now := time.Now()

	if s.LastAccessedAt.Add(s.IdleTimeout).Before(now) {
		return true
	}

	if s.CreatedAt.Add(s.AbsTimeout).Before(now) {
		return true
	}

	return false
}

func ValidateSessionId(sid string) error {
	if len(sid) != base32enc.EncodedLen(keyLen) {
		return fmt.Errorf("error validating session: wrong session id length")
	}
	_, err := base32enc.DecodeString(sid)
	if err != nil {
		return fmt.Errorf("error validating session: %w", err)
	}
	return nil
}

func generateSessionId() (string, error) {
	id := make([]byte, keyLen)
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		return "", fmt.Errorf("error generating session id: %w", err)
	}
	return base32enc.EncodeToString(id), nil
}
