package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
)

var UIDAttr SessionKey = "uid"

const (
	LogKeySID        = "session.sid"
	LogKeyRQID       = "session.rqud"
	LogKeyDebugError = "session.dbg_error"
)

var ErrSessionNotFound = errors.New("sessionservice: session not found")

type SessionServiceInterface interface {
	CreateAnonymSession(ctx context.Context, cc CookieConf, sc SessionConf, keyAndValues ...interface{}) (Session, error)
	CreateUserSession(ctx context.Context, uid string, cc CookieConf, sc SessionConf, keyAndValues ...interface{}) (Session, error)
	LoadSession(ctx context.Context, sid string) (Session, error)
	InvalidateSession(ctx context.Context, sid string) error
}

type sessionService struct {
	Logger      logr.Logger
	SStore      SessionStore
	CtxReqIdKey interface{}
}

/*

Create implementation of SessionServiceInterface to work with session

logr.Logger may be usefull only for debugging purposes.
All messages will be logged as logr.V = 10
and none of errors will be logged as logr.Error.
It's up to service that call these methods
to decide what is really error in terms of application and what is not.

reqIdKey is key to extract request id from context

*/
func NewService(s SessionStore, l logr.Logger, reqIdKey interface{}) SessionServiceInterface {
	ss := sessionService{
		Logger:      l,
		SStore:      s,
		CtxReqIdKey: reqIdKey,
	}
	return &ss
}

func (ss *sessionService) CreateAnonymSession(ctx context.Context, cc CookieConf, sc SessionConf, keyAndValues ...interface{}) (Session, error) {
	ss.Logger.V(0).Info(
		"session.CreateAnonymSession() started",
		LogKeyRQID, ctx.Value(ss.CtxReqIdKey))
	defer ss.Logger.V(0).Info(
		"session.CreateAnonymSession() finished",
		LogKeyRQID, ctx.Value(ss.CtxReqIdKey))

	if len(keyAndValues)%2 != 0 {
		err := fmt.Errorf("session.CreateAnonymSession() expected even count of key and values, got: %v", len(keyAndValues))
		ss.Logger.V(0).Info(
			"session.CreateAnonymSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIdKey),
			LogKeyDebugError, err)
		return Session{}, err
	}

	s, err := NewSession()
	if err != nil {
		err := fmt.Errorf("session.CreateAnonymSession() error creating anon session: %w", err)
		ss.Logger.V(0).Info(
			"session.CreateAnonymSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIdKey),
			LogKeyDebugError, err)
		return s, err
	}
	s.WithCookieConf(cc)
	s.WithSessionConf(sc)

	for i := 0; i < len(keyAndValues); i += 2 {
		k, ok := keyAndValues[i].(SessionKey)
		if !ok {
			err := fmt.Errorf("session.CreateAnonymSession() error: can't convert key of type: %T to session.SessionKey", keyAndValues[i])
			ss.Logger.V(0).Info(
				"session.CreateAnonymSession() error",
				LogKeyRQID, ctx.Value(ss.CtxReqIdKey),
				LogKeyDebugError, err)
			return Session{}, err
		}
		s.AddAttribute(k, keyAndValues[i+1])
	}

	svdS, err := ss.SStore.Save(ctx, s)
	if err != nil {
		err := fmt.Errorf("session.CreateAnonymSession() Save error: %w", err)
		ss.Logger.V(0).Info(
			"session.CreateAnonymSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIdKey),
			LogKeyDebugError, err)
		return Session{}, err
	}

	return svdS, nil
}

func (ss *sessionService) CreateUserSession(ctx context.Context, uid string, cc CookieConf, sc SessionConf, keyAndValues ...interface{}) (Session, error) {
	ss.Logger.V(0).Info("session.CreateUserSession() started", LogKeyRQID, ctx.Value(ss.CtxReqIdKey))
	defer ss.Logger.V(0).Info("session.CreateUserSession() finished", LogKeyRQID, ctx.Value(ss.CtxReqIdKey))

	s, err := NewSession()
	if err != nil {
		err := fmt.Errorf("session.CreateUserSession() error creating user session: %w", err)
		ss.Logger.V(0).Info(
			"session.CreateUserSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIdKey),
			LogKeyDebugError, err)
		return Session{}, err
	}
	s.WithCookieConf(cc)
	s.WithUserId(uid)
	s.WithSessionConf(sc)

	s.UID = uid

	svdS, err := ss.SStore.Save(ctx, s)
	if err != nil {
		err := fmt.Errorf("session.CreateUserSession() Save error: %w", err)
		ss.Logger.V(0).Info(
			"session.CreateUserSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIdKey),
			LogKeyDebugError, err)
		return Session{}, err
	}

	return svdS, nil
}

func (ss *sessionService) LoadSession(ctx context.Context, sid string) (Session, error) {
	ss.Logger.V(0).Info("session.LoadSession() started", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIdKey))
	defer ss.Logger.V(0).Info("session.LoadSession() finished", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIdKey))

	s, err := ss.SStore.Load(ctx, sid)

	if err == ErrSessionNotFound {
		return Session{}, ErrSessionNotFound
	}

	if err != nil {
		err := fmt.Errorf("session.LoadSession() Load error: %w", err)
		ss.Logger.V(0).Info(
			"session.LoadSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIdKey),
			LogKeyDebugError, err)
		return Session{}, err
	}

	return s, nil
}

func (ss *sessionService) InvalidateSession(ctx context.Context, sid string) error {
	ss.Logger.V(0).Info("session.InvalidateSession() started", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIdKey))
	defer ss.Logger.V(0).Info("session.InvalidateSession() finished", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIdKey))

	err := ss.SStore.Invalidate(ctx, sid)
	if err != nil {
		err := fmt.Errorf("session.InvalidateSession() Invalidate error: %w", err)
		ss.Logger.V(0).Info(
			"session.InvalidateSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIdKey),
			LogKeyDebugError, err)
		return err
	}
	return nil
}
