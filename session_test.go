package session_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/asstart/go-session"
	"github.com/stretchr/testify/assert"
)

func TestSessionValidationSuccess(t *testing.T) {
	sid := "A7TF7SGM5WZRW7WMGY7BRJPQOGWGXATZWT35HXPKHRO3DU2J3L4Q"
	err := session.ValidateSessionID(sid)
	assert.Nil(t, err)
}

func TestSessionValidationFail(t *testing.T) {
	tt := []struct {
		name   string
		sid    string
		expErr error
	}{
		{"shorter 1 symbol sid", "A7TF7SGM5WZRW7WMGY7BRJPQOGWGXATZWT35HXPKHRO3DU2J3L4", errors.New("error validating session: wrong session id length")},
		{"1 symbol sid", "A", errors.New("error validating session: wrong session id length")},
		{"empty sid", "", errors.New("error validating session: wrong session id length")},
		{"wrong alphabet sid", "!7TF7SGM5WZRW7WMGY7BRJPQOGWGXATZWT35HXPKHRO3DU2J3L4Q", fmt.Errorf("error validating session: %w", fmt.Errorf("illegal base32 data at input byte 0"))},
		{"longer 1 symbok sid", "A7TF7SGM5WZRW7WMGY7BRJPQOGWGXATZWT35HXPKHRO3DU2J3L4QA", errors.New("error validating session: wrong session id length")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := session.ValidateSessionID(tc.sid)
			assert.Equal(t, tc.expErr.Error(), err.Error())
		})
	}
}

func TestSessionIsExpired(t *testing.T) {
	tt := []struct {
		name       string
		expExpired bool
		s          session.Session
	}{
		{"inactive anon session", true, session.Session{Anonym: true, Active: false, IdleTimeout: 60 * time.Minute, AbsTimeout: 60 * time.Minute, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
		{"inactive user session", true, session.Session{UID: "111", Anonym: false, Active: false, IdleTimeout: 60 * time.Minute, AbsTimeout: 60 * time.Minute, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
		{"active, anon ses, expired idle, not expired abs", true, session.Session{Anonym: true, Active: true, IdleTimeout: 0, AbsTimeout: 60 * time.Minute, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
		{"active, user ses, expired idle, not expired abs", true, session.Session{UID: "1111", Anonym: false, Active: true, IdleTimeout: 0, AbsTimeout: 60 * time.Minute, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
		{"active, anon ses, expired abs, not expired idle", true, session.Session{Anonym: true, Active: true, IdleTimeout: 60 * time.Minute, AbsTimeout: 0, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
		{"active, user ses, expired abs, not expired idle", true, session.Session{UID: "1111", Anonym: false, Active: true, IdleTimeout: 60 * time.Minute, AbsTimeout: 0, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
		{"active, anon ses, expired abs, expired idle", true, session.Session{Anonym: true, Active: true, IdleTimeout: 0, AbsTimeout: 0, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
		{"active, user ses, expired abs, expired idle", true, session.Session{UID: "1111", Anonym: false, Active: true, IdleTimeout: 0, AbsTimeout: 0, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
		{"active, anon ses, not expired abs", false, session.Session{Anonym: true, Active: true, IdleTimeout: 60 * time.Minute, AbsTimeout: 60 * time.Minute, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
		{"active, user ses, not expired", false, session.Session{UID: "1111", Anonym: false, Active: true, IdleTimeout: 60 * time.Minute, AbsTimeout: 60 * time.Minute, CreatedAt: time.Now(), LastAccessedAt: time.Now()}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expExpired, tc.s.IsExpired())
		})
	}
}
