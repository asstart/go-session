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

func TestCreateAnonSessionBadAttribbutes(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

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
		{"odd number of key value paiers, single key", []interface{}{"key"}, "session.CreateAnonymSession() error: expected even count of key and values, got: 1"},
		{"odd number of key value paiers, multiple keys", []interface{}{"key1", "value1", "key2"}, "session.CreateAnonymSession() error: expected even count of key and values, got: 3"},
		{"invalid int key", []interface{}{1, "value"}, "session.CreateAnonymSession() error: can't convert key of type: int to string"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := service.CreateAnonymSession(context.Background(), session.DefaultCookieConf(), session.DefaultSessionConf(), tc.kv...)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expErr, err.Error())
			assert.Nil(t, s)
		})
	}
}

func TestCreateAnonSessionValidAttributes(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	var k1 = "key1"
	var k2 = "key2"

	tt := []struct {
		name         string
		kv           []interface{}
		expAttrCount int
		expData      map[string]interface{}
	}{
		{"empty key value pairs", []interface{}{}, 0, make(map[string]interface{})},
		{"single pair", []interface{}{k1, "string value"}, 1, map[string]interface{}{k1: "string value"}},
		{"multiple pair", []interface{}{k1, "string value", k2, "string value 2"}, 2, map[string]interface{}{k1: "string value", k2: "string value 2"}},
		{"int value", []interface{}{k1, 1}, 1, map[string]interface{}{k1: 1}},
		{"float value", []interface{}{k1, 1.0}, 1, map[string]interface{}{k1: 1.0}},
		{"bool value", []interface{}{k1, true}, 1, map[string]interface{}{k1: true}},
		{"array value", []interface{}{k1, [1]int{1}}, 1, map[string]interface{}{k1: [1]int{1}}},
		{"slice value", []interface{}{k1, []int{1, 2}}, 1, map[string]interface{}{k1: []int{1, 2}}},
		{"struct value", []interface{}{k1, struct{}{}}, 1, map[string]interface{}{k1: struct{}{}}},
		{"map value", []interface{}{k1, make(map[string]string, 1)}, 1, map[string]interface{}{k1: make(map[string]string, 1)}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			retSes := session.Session{Data: tc.expData}
			smock.EXPECT().Save(gomock.Any(), sesDataMatcher{&retSes}).Return(&retSes, nil)

			created, err := service.CreateAnonymSession(context.Background(), session.DefaultCookieConf(), session.DefaultSessionConf(), tc.kv...)
			assert.Nil(t, err)
			assert.Same(t, &retSes, created)
		})
	}
}

func TestCreateUserSessionBadAttribbutes(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	uid := "2222"

	tt := []struct {
		name   string
		kv     []interface{}
		expErr string
	}{
		{"odd number of key value paiers, single key", []interface{}{"key"}, "session.CreateUserSession() error: expected even count of key and values, got: 1"},
		{"odd number of key value paiers, multiple keys", []interface{}{"key1", "value1", "key2"}, "session.CreateUserSession() error: expected even count of key and values, got: 3"},
		{"invalid int key", []interface{}{1, "value"}, "session.CreateUserSession() error: can't convert key of type: int to string"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := service.CreateUserSession(context.Background(), uid, session.DefaultCookieConf(), session.DefaultSessionConf(), tc.kv...)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expErr, err.Error())
			assert.Nil(t, s)
		})
	}
}

func TestCreateUserSessionValidAttributes(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	uid := "2222"
	var k1 = "key1"
	var k2 = "key2"

	tt := []struct {
		name         string
		kv           []interface{}
		expAttrCount int
		expData      map[string]interface{}
	}{
		{"empty key value pairs", []interface{}{}, 0, make(map[string]interface{})},
		{"single pair", []interface{}{k1, "string value"}, 1, map[string]interface{}{k1: "string value"}},
		{"multiple pair", []interface{}{k1, "string value", k2, "string value 2"}, 2, map[string]interface{}{k1: "string value", k2: "string value 2"}},
		{"int value", []interface{}{k1, 1}, 1, map[string]interface{}{k1: 1}},
		{"float value", []interface{}{k1, 1.0}, 1, map[string]interface{}{k1: 1.0}},
		{"bool value", []interface{}{k1, true}, 1, map[string]interface{}{k1: true}},
		{"array value", []interface{}{k1, [1]int{1}}, 1, map[string]interface{}{k1: [1]int{1}}},
		{"slice value", []interface{}{k1, []int{1, 2}}, 1, map[string]interface{}{k1: []int{1, 2}}},
		{"struct value", []interface{}{k1, struct{}{}}, 1, map[string]interface{}{k1: struct{}{}}},
		{"map value", []interface{}{k1, make(map[string]string, 1)}, 1, map[string]interface{}{k1: make(map[string]string, 1)}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			retSes := session.Session{Data: tc.expData}
			smock.EXPECT().Save(gomock.Any(), sesDataMatcher{&retSes}).Return(&retSes, nil)

			created, err := service.CreateUserSession(context.Background(), uid, session.DefaultCookieConf(), session.DefaultSessionConf(), tc.kv...)
			assert.Nil(t, err)
			assert.Same(t, &retSes, created)
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

	v1, ok := sm.x.(*session.Session)
	if !ok {
		return false
	}

	v2, ok := x.(*session.Session)
	if !ok {
		return false
	}

	if len(v1.Data) != len(v2.Data) {
		return false
	}

	return reflect.DeepEqual(v1.Data, v2.Data)
}

func (sm sesDataMatcher) String() string {
	return fmt.Sprintf("%v", sm.x)
}

type sesFullMatcher struct {
	x interface{}
}

func (sm sesFullMatcher) Matches(x interface{}) bool {
	if sm.x == nil || x == nil {
		return reflect.DeepEqual(sm.x, x)
	}

	v1, ok := sm.x.(*session.Session)
	if !ok {
		return false
	}
	v2, ok := x.(*session.Session)
	if !ok {
		return false
	}

	return session.ValidateSessionID(v2.ID) == nil &&
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
	return fmt.Sprintf("%v", sm.x)
}

func TestSuccessfulCreatingAnonymSession(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	ctx := context.Background()
	cookieConf := session.DefaultCookieConf()
	sconf := session.DefaultSessionConf()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	resSes := session.Session{}

	smock.EXPECT().Save(ctx, sesFullMatcher{
		&session.Session{
			Active:      true,
			Anonym:      true,
			UID:         "",
			Opts:        cookieConf,
			IdleTimeout: sconf.IdleTimeout,
			AbsTimeout:  sconf.AbsTimout,
		},
	}).Return(&resSes, nil)

	created, err := service.CreateAnonymSession(ctx, cookieConf, sconf)

	assert.Nil(t, err)
	assert.Same(t, &resSes, created)
}

func TestErrorSavingAnonymSession(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

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

	smock.EXPECT().Save(ctx, gomock.Any()).Return(nil, retErr)

	created, err := service.CreateAnonymSession(ctx, cookieConf, sconf)

	assert.Nil(t, created)
	assert.Equal(t, expErr, err)
}

func TestSucessfullCreateUserSession(t *testing.T) {

	smock := smocks.NewMockStore(gomock.NewController(t))

	ctx := context.Background()
	cookieConf := session.DefaultCookieConf()
	sconf := session.DefaultSessionConf()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	uid := "1234"

	resSes := session.Session{}

	smock.EXPECT().Save(ctx, sesFullMatcher{
		&session.Session{
			Active:      true,
			Anonym:      false,
			UID:         uid,
			Opts:        cookieConf,
			IdleTimeout: sconf.IdleTimeout,
			AbsTimeout:  sconf.AbsTimout,
		},
	}).Return(&resSes, nil)

	created, err := service.CreateUserSession(ctx, uid, cookieConf, sconf)

	assert.Nil(t, err)
	assert.Same(t, &resSes, created)
}

func TestErrorSavingUserSession(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

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
	smock.EXPECT().Save(ctx, gomock.Any()).Return(nil, retErr)

	uid := "1234"

	created, err := service.CreateUserSession(ctx, uid, cookieConf, sconf)
	assert.Nil(t, created)
	assert.Equal(t, expErr, err)
}

func TestLoadSessionErr(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

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
		{"session not found", session.ErrSessionNotFound, session.ErrSessionNotFound},
		{"session not found", errors.New("some error"), fmt.Errorf("session.LoadSession() Load error: %w", errors.New("some error"))},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			smock.EXPECT().Load(ctx, sid).Return(nil, tc.returnErr)

			created, err := service.LoadSession(ctx, sid)
			assert.Nil(t, created)
			assert.Equal(t, tc.expErr, err)
		})
	}
}

func TestLoadSessionSuccess(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	ctx := context.Background()

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	sid := "1234"

	ses, _ := session.NewSession()
	ses.ID = sid

	smock.EXPECT().Load(ctx, sid).Return(&ses, nil)

	loaded, err := service.LoadSession(ctx, sid)
	assert.Nil(t, err)
	assert.Same(t, &ses, loaded)
}

func TestInvalidateSession(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

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

func TestAddAttributesValidCases(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	var k1 = "key1"
	var k2 = "key2"
	sid := "1111"

	tt := []struct {
		name         string
		kv           []interface{}
		expAttrCount int
		expData      map[string]interface{}
	}{
		{"single pair", []interface{}{k1, "string value"}, 1, map[string]interface{}{k1: "string value"}},
		{"multiple pair", []interface{}{k1, "string value", k2, "string value 2"}, 2, map[string]interface{}{k1: "string value", k2: "string value 2"}},
		{"int value", []interface{}{k1, 1}, 1, map[string]interface{}{k1: 1}},
		{"float value", []interface{}{k1, 1.0}, 1, map[string]interface{}{k1: 1.0}},
		{"bool value", []interface{}{k1, true}, 1, map[string]interface{}{k1: true}},
		{"array value", []interface{}{k1, [1]int{1}}, 1, map[string]interface{}{k1: [1]int{1}}},
		{"slice value", []interface{}{k1, []int{1, 2}}, 1, map[string]interface{}{k1: []int{1, 2}}},
		{"struct value", []interface{}{k1, struct{}{}}, 1, map[string]interface{}{k1: struct{}{}}},
		{"map value", []interface{}{k1, make(map[string]string, 1)}, 1, map[string]interface{}{k1: make(map[string]string, 1)}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			retSes := session.Session{Data: tc.expData}
			smock.EXPECT().AddAttributes(gomock.Any(), gomock.Eq(sid), gomock.Eq(tc.expData)).Return(&retSes, nil)

			created, err := service.AddAttributes(context.Background(), sid, tc.kv...)
			assert.Nil(t, err)
			assert.Same(t, &retSes, created)
		})
	}
}

func TestAddAttributesInvalidCases(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	service := session.NewService(
		smock,
		logr.Discard(),
		"key",
	)

	sid := "1111"

	tt := []struct {
		name   string
		kv     []interface{}
		expErr string
	}{
		{"empty key value pairs", []interface{}{}, "session.AddAttributes() no attributes to add"},
		{"odd number of key value paiers, single key", []interface{}{"key"}, "session.AddAttributes() error: expected even count of key and values, got: 1"},
		{"odd number of key value paiers, multiple keys", []interface{}{"key1", "value1", "key2"}, "session.AddAttributes() error: expected even count of key and values, got: 3"},
		{"invalid int key", []interface{}{1, "value"}, "session.AddAttributes() error: can't convert key of type: int to string"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := service.AddAttributes(context.Background(), sid, tc.kv...)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expErr, err.Error())
			assert.Nil(t, s)
		})
	}
}

func TestAddAttributesSessionNotFound(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	service := session.NewService(smock, logr.Discard(), "key")

	sid := "1111"

	var k = "attr1"
	v := "value"
	expData := map[string]interface{}{
		k: v,
	}
	smock.EXPECT().AddAttributes(gomock.Any(), gomock.Eq(sid), gomock.Eq(expData)).Return(nil, session.ErrSessionNotFound)
	s, err := service.AddAttributes(context.Background(), sid, k, v)
	assert.Nil(t, s)
	assert.Equal(t, session.ErrSessionNotFound, err)
}

func TestAddAttributesUnexpectedErr(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	service := session.NewService(smock, logr.Discard(), "key")

	sid := "1111"

	var k = "attr1"
	v := "value"
	expData := map[string]interface{}{
		k: v,
	}
	retErr := errors.New("unexpected error")
	smock.
		EXPECT().
		AddAttributes(gomock.Any(), gomock.Eq(sid), gomock.Eq(expData)).
		Return(nil, retErr)
	s, err := service.AddAttributes(context.Background(), sid, k, v)
	assert.Nil(t, s)
	assert.Equal(t, fmt.Errorf("session.AddAttributes() AddAttributes unexpected error: %w", retErr), err)
}

func TestAttributesRemoveNoAttrs(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))
	service := session.NewService(smock, logr.Discard(), "key")

	sid := "1111"
	s, err := service.RemoveAttributes(context.Background(), sid)

	assert.Nil(t, s)
	assert.Equal(t, fmt.Errorf("session.RemoveAttributes() no attributes to remove"), err)
}

func TestRemoveAttributesSessionNotFound(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	service := session.NewService(smock, logr.Discard(), "key")

	sid := "1111"

	var k = "attr1"
	smock.EXPECT().RemoveAttributes(gomock.Any(), gomock.Eq(sid), gomock.Eq(k)).Return(nil, session.ErrSessionNotFound)
	s, err := service.RemoveAttributes(context.Background(), sid, k)
	assert.Nil(t, s)
	assert.Equal(t, session.ErrSessionNotFound, err)
}

func TestRemoveAttributesUnexpectedErr(t *testing.T) {
	smock := smocks.NewMockStore(gomock.NewController(t))

	service := session.NewService(smock, logr.Discard(), "key")

	sid := "1111"

	var k = "attr1"

	retErr := errors.New("unexpected error")
	smock.
		EXPECT().
		RemoveAttributes(gomock.Any(), gomock.Eq(sid), gomock.Eq(k)).
		Return(nil, retErr)
	s, err := service.RemoveAttributes(context.Background(), sid, k)
	assert.Nil(t, s)
	assert.Equal(t, fmt.Errorf("session.RemoveAttributes() RemoveAttributes unexpected error: %w", retErr), err)
}
