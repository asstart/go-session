package session_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/asstart/go-session"
	smocks "github.com/asstart/go-session/mocks"
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBadAttribbutes(t *testing.T) {
	smock := smocks.NewMockSessionStore(gomock.NewController(t))

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	tt := []struct {
		name   string
		kv     []interface{}
		expErr string
	}{
		{"odd number of key value paiers, single key", []interface{}{"key"}, "session.CreateAnonymSession() expected even count of key and values, got: 1"},
		{"odd number of key value paiers, multiple keys", []interface{}{"key1", "value1", "key2"}, "session.CreateAnonymSession() expected even count of key and values, got: 3"},
		{"invalid string key", []interface{}{"1", "value"}, "session.CreateAnonymSession() error: can't convert key of type: string to session.SessionKey"},
		{"invalid int key", []interface{}{1, "value"}, "session.CreateAnonymSession() error: can't convert key of type: int to session.SessionKey"},
	}

	emptySes := session.Session{}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := service.CreateAnonymSession(context.Background(), session.DefaultCookieConf(), session.DefaultSessionConf(), tc.kv...)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expErr, err.Error())
			assert.Equal(t, emptySes, s)
		})
	}
}

func TestValidAttributes(t *testing.T) {
	smock := smocks.NewMockSessionStore(gomock.NewController(t))

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	var k1 session.SessionKey = "key1"
	var k2 session.SessionKey = "key2"

	tt := []struct {
		name         string
		kv           []interface{}
		expAttrCount int
		expData      map[session.SessionKey]interface{}
	}{
		{"empty key value pairs", []interface{}{}, 0, make(map[session.SessionKey]interface{})},
		{"single pair", []interface{}{k1, "string value"}, 1, map[session.SessionKey]interface{}{k1: "string value"}},
		{"multiple pair", []interface{}{k1, "string value", k2, "string value 2"}, 2, map[session.SessionKey]interface{}{k1: "string value", k2: "string value 2"}},
		{"int value", []interface{}{k1, 1}, 1, map[session.SessionKey]interface{}{k1: 1}},
		{"float value", []interface{}{k1, 1.0}, 1, map[session.SessionKey]interface{}{k1: 1.0}},
		{"bool value", []interface{}{k1, true}, 1, map[session.SessionKey]interface{}{k1: true}},
		{"array value", []interface{}{k1, [1]int{1}}, 1, map[session.SessionKey]interface{}{k1: [1]int{1}}},
		{"slice value", []interface{}{k1, []int{1, 2}}, 1, map[session.SessionKey]interface{}{k1: []int{1, 2}}},
		{"struct value", []interface{}{k1, struct{}{}}, 1, map[session.SessionKey]interface{}{k1: struct{}{}}},
		{"map value", []interface{}{k1, make(map[string]string, 1)}, 1, map[session.SessionKey]interface{}{k1: make(map[string]string, 1)}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			smock.EXPECT().Save(gomock.Any(), sesDataMatcher{session.Session{Data: tc.expData}}).Return(session.Session{}, nil)

			_, err := service.CreateAnonymSession(context.Background(), session.DefaultCookieConf(), session.DefaultSessionConf(), tc.kv...)
			assert.Nil(t, err)
		})
	}
}

type sesDataMatcher struct {
	x interface{}
}

func (sm sesDataMatcher) Matches(x interface{}) bool {
	if sm.x == nil || x == nil {
		return reflect.DeepEqual(sm.x, x)
	}

	v1, ok := sm.x.(session.Session)
	if !ok {
		return false
	}
	v2, ok := x.(session.Session)
	if !ok {
		return false
	}

	if len(v1.Data) != len(v2.Data) {
		return false
	}

	return reflect.DeepEqual(v1.Data, v2.Data)
}

func (sm sesDataMatcher) String() string {
	return "is session.Session.Data"
}

type sesFullMatcher struct {
	x interface{}
}

func (sm sesFullMatcher) Matches(x interface{}) bool {
	if sm.x == nil || x == nil {
		return reflect.DeepEqual(sm.x, x)
	}

	v1, ok := sm.x.(session.Session)
	if !ok {
		return false
	}
	v2, ok := x.(session.Session)
	if !ok {
		return false
	}

	return session.ValidateSessionId(v2.ID) == nil &&
		reflect.DeepEqual(v1.Active, v2.Active) &&
		reflect.DeepEqual(v1.Anonym, v2.Anonym) &&
		reflect.DeepEqual(v1.UID, v2.UID) &&
		reflect.DeepEqual(v1.Opts, v2.Opts) &&
		reflect.DeepEqual(v1.IdleTimeout, v2.IdleTimeout) &&
		reflect.DeepEqual(v1.AbsTimeout, v2.AbsTimeout) &&
		v2.LastAccessedAt == time.Time{} &&
		v2.CreatedAt == time.Time{}
}

func (sm sesFullMatcher) String() string {
	return "is session.Session"
}

func TestSuccessfulCreatingAnonymSession(t *testing.T) {
	smock := smocks.NewMockSessionStore(gomock.NewController(t))

	ctx := context.Background()
	cookieConf := session.DefaultCookieConf()
	sconf := session.DefaultSessionConf()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	smock.EXPECT().Save(ctx, sesFullMatcher{
		session.Session{
			Active:      true,
			Anonym:      true,
			UID:         "",
			Opts:        cookieConf,
			IdleTimeout: sconf.IdleTimeout,
			AbsTimeout:  sconf.AbsTimout,
		},
	}).Return(session.Session{}, nil)

	_, err := service.CreateAnonymSession(ctx, cookieConf, sconf)

	assert.Nil(t, err)
}

func TestErrorSavingAnonymSession(t *testing.T) {
	smock := smocks.NewMockSessionStore(gomock.NewController(t))

	ctx := context.Background()
	cookieConf := session.DefaultCookieConf()
	sconf := session.DefaultSessionConf()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	retErr := errors.New("some err")
	expErr := fmt.Errorf("session.CreateAnonymSession() Save error: %w", retErr)

	smock.EXPECT().Save(ctx, gomock.Any()).Return(session.Session{}, retErr)

	s, err := service.CreateAnonymSession(ctx, cookieConf, sconf)

	assert.Equal(t, session.Session{}, s)
	assert.Equal(t, expErr, err)
}

func TestSucessfullCreateUserSession(t *testing.T) {

	smock := smocks.NewMockSessionStore(gomock.NewController(t))

	ctx := context.Background()
	cookieConf := session.DefaultCookieConf()
	sconf := session.DefaultSessionConf()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	uid := "1234"

	smock.EXPECT().Save(ctx, sesFullMatcher{
		session.Session{
			Active:      true,
			Anonym:      false,
			UID:         uid,
			Opts:        cookieConf,
			IdleTimeout: sconf.IdleTimeout,
			AbsTimeout:  sconf.AbsTimout,
		},
	}).Return(session.Session{}, nil)

	_, err := service.CreateUserSession(ctx, uid, cookieConf, sconf)

	assert.Nil(t, err)
}

func TestErrorSavingUserSession(t *testing.T) {
	smock := smocks.NewMockSessionStore(gomock.NewController(t))

	ctx := context.Background()
	cookieConf := session.DefaultCookieConf()
	sconf := session.DefaultSessionConf()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	retErr := errors.New("some error")
	expErr := fmt.Errorf("session.CreateUserSession() Save error: %w", retErr)
	smock.EXPECT().Save(ctx, gomock.Any()).Return(session.Session{}, retErr)

	uid := "1234"

	s, err := service.CreateUserSession(ctx, uid, cookieConf, sconf)
	assert.Equal(t, session.Session{}, s)
	assert.Equal(t, expErr, err)
}

func TestLoadSessionErr(t *testing.T) {
	smock := smocks.NewMockSessionStore(gomock.NewController(t))

	ctx := context.Background()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	sid := "1234"

	tt := []struct {
		name      string
		returnErr error
		returnSes session.Session
		expErr    error
		expSes    session.Session
	}{
		{"session not found", session.ErrSessionNotFound, session.Session{}, session.ErrSessionNotFound, session.Session{}},
		{"session not found", errors.New("some error"), session.Session{}, fmt.Errorf("session.LoadSession() Load error: %w", errors.New("some error")), session.Session{}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			smock.EXPECT().Load(ctx, sid).Return(tc.returnSes, tc.returnErr)

			s, err := service.LoadSession(ctx, sid)
			assert.Equal(t, tc.expSes, s)
			assert.Equal(t, tc.expErr, err)
		})
	}
}

func TestLoadSessionSuccess(t *testing.T) {
	smock := smocks.NewMockSessionStore(gomock.NewController(t))

	ctx := context.Background()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	sid := "1234"

	ses, _ := session.NewSession()
	ses.ID = sid

	smock.EXPECT().Load(ctx, sid).Return(ses, nil)

	s, err := service.LoadSession(ctx, sid)
	assert.Nil(t, err)
	assert.Equal(t, ses, s)
}

func TestInvalidateSession(t *testing.T) {
	smock := smocks.NewMockSessionStore(gomock.NewController(t))

	ctx := context.Background()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	sid := "1234"

	tt := []struct {
		name      string
		returnErr error
		expErr    error
	}{
		{"invalidate failed", errors.New("some error"), fmt.Errorf("session.InvalidateSession() Invalidate error: %w", errors.New("some error"))},
		{"invalidate success", nil, nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			smock.EXPECT().Invalidate(ctx, sid).Return(tc.returnErr)
			err := service.InvalidateSession(ctx, sid)
			assert.Equal(t, tc.expErr, err)
		})
	}
}
