package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
)

const (
	LogKeySID        = "session.sid"
	LogKeyRQID       = "session.rqud"
	LogKeyDebugError = "session.dbg_error"
)

var ErrSessionNotFound = errors.New("sessionservice: session not found")

type Service interface {
	CreateAnonymSession(ctx context.Context, cc CookieConf, sc Conf, keyAndValues ...interface{}) (*Session, error)
	CreateUserSession(ctx context.Context, uid string, cc CookieConf, sc Conf, keyAndValues ...interface{}) (*Session, error)
	LoadSession(ctx context.Context, sid string) (*Session, error)
	InvalidateSession(ctx context.Context, sid string) error
	AddAttributes(ctx context.Context, sid string, keyAndValues ...interface{}) (*Session, error)
	RemoveAttributes(ctx context.Context, sid string, keys ...string) (*Session, error)
}

type sessionService struct {
	Logger      logr.Logger
	SStore      Store
	CtxReqIDKey interface{}
}

/*

NewService Create implementation of Service to work with session

logr.Logger may be useful only for debugging purposes,
Nnone of errors will be logged as logr.Error.
It's up to service that call these methods
to decide what is really error in terms of application and what is not.

reqIDKey is key to extract request id from the context

*/
func NewService(s Store, l logr.Logger, reqIDKey interface{}) Service {
	ss := sessionService{
		Logger:      l,
		SStore:      s,
		CtxReqIDKey: reqIDKey,
	}
	return &ss
}

// CreateAnonymSession create new anonym session, store it based on provided implementation of Store and return
// keyAndValues attributes which should be added to a session during creation
func (ss *sessionService) CreateAnonymSession(ctx context.Context, cc CookieConf, sc Conf, keyAndValues ...interface{}) (*Session, error) {
	ss.Logger.V(0).Info(
		"session.CreateAnonymSession() started",
		LogKeyRQID, ctx.Value(ss.CtxReqIDKey))
	defer ss.Logger.V(0).Info(
		"session.CreateAnonymSession() finished",
		LogKeyRQID, ctx.Value(ss.CtxReqIDKey))

	data, err := parseAttrs(keyAndValues...)
	if err != nil {
		err = fmt.Errorf("session.CreateAnonymSession() error: %v", err)
		ss.Logger.V(0).Info(
			"session.CreateAnonymSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err,
		)
		return nil, err
	}

	s, err := NewSession()
	if err != nil {
		err = fmt.Errorf("session.CreateAnonymSession() error creating anon session: %w", err)
		ss.Logger.V(0).Info(
			"session.CreateAnonymSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err)
		return nil, err
	}

	s.WithCookieConf(cc)
	s.WithSessionConf(sc)
	s.WithAttributes(data)

	svdS, err := ss.SStore.Save(ctx, &s)
	if err != nil {
		err = fmt.Errorf("session.CreateAnonymSession() Save error: %w", err)
		ss.Logger.V(0).Info(
			"session.CreateAnonymSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err)
		return nil, err
	}

	return svdS, nil
}

// CreateUserSession create new session, store it based on provided implementation of Store and return
// keyAndValues attributes which should be added to a session during creation
func (ss *sessionService) CreateUserSession(ctx context.Context, uid string, cc CookieConf, sc Conf, keyAndValues ...interface{}) (*Session, error) {
	ss.Logger.V(0).Info("session.CreateUserSession() started", LogKeyRQID, ctx.Value(ss.CtxReqIDKey))
	defer ss.Logger.V(0).Info("session.CreateUserSession() finished", LogKeyRQID, ctx.Value(ss.CtxReqIDKey))

	data, err := parseAttrs(keyAndValues...)
	if err != nil {
		err = fmt.Errorf("session.CreateUserSession() error: %w", err)
		ss.Logger.V(0).Info(
			"session.CreateUserSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err,
		)
		return nil, err
	}

	s, err := NewSession()
	if err != nil {
		err = fmt.Errorf("session.CreateUserSession() error creating user session: %w", err)
		ss.Logger.V(0).Info(
			"session.CreateUserSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err)
		return nil, err
	}

	s.WithCookieConf(cc)
	s.WithUserID(uid)
	s.WithSessionConf(sc)
	s.WithAttributes(data)

	svdS, err := ss.SStore.Save(ctx, &s)
	if err != nil {
		err = fmt.Errorf("session.CreateUserSession() Save error: %w", err)
		ss.Logger.V(0).Info(
			"session.CreateUserSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err)
		return nil, err
	}

	return svdS, nil
}

// LoadSession return session loaded from storage based on implementation of Store
func (ss *sessionService) LoadSession(ctx context.Context, sid string) (*Session, error) {
	ss.Logger.V(0).Info("session.LoadSession() started", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIDKey))
	defer ss.Logger.V(0).Info("session.LoadSession() finished", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIDKey))

	s, err := ss.SStore.Load(ctx, sid)

	if err == ErrSessionNotFound {
		return nil, ErrSessionNotFound
	}

	if err != nil {
		err = fmt.Errorf("session.LoadSession() Load error: %w", err)
		ss.Logger.V(0).Info(
			"session.LoadSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err)
		return nil, err
	}

	return s, nil
}

// InvalidateSession invalidate session in storage based on implementation of Store
func (ss *sessionService) InvalidateSession(ctx context.Context, sid string) error {
	ss.Logger.V(0).Info("session.InvalidateSession() started", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIDKey))
	defer ss.Logger.V(0).Info("session.InvalidateSession() finished", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIDKey))

	err := ss.SStore.Invalidate(ctx, sid)
	if err != nil {
		err = fmt.Errorf("session.InvalidateSession() Invalidate error: %w", err)
		ss.Logger.V(0).Info(
			"session.InvalidateSession() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err)
		return err
	}
	return nil
}

func (ss *sessionService) AddAttributes(ctx context.Context, sid string, keyAndValues ...interface{}) (*Session, error) {
	ss.Logger.V(0).Info("session.AddAttributes() started", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIDKey))
	defer ss.Logger.V(0).Info("session.AddAttributes() finished", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIDKey))

	if len(keyAndValues) == 0 {
		err := fmt.Errorf("session.AddAttributes() no attributes to add")
		ss.Logger.V(0).Info(
			"session.AddAttributes() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err,
			LogKeySID, sid,
		)
		return nil, err
	}

	data, err := parseAttrs(keyAndValues...)
	if err != nil {
		err = fmt.Errorf("session.AddAttributes() error: %v", err)
		ss.Logger.V(0).Info(
			"session.AddAttributes() error",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeySID, sid,
			LogKeyDebugError, err,
		)
		return nil, err
	}

	s, err := ss.SStore.AddAttributes(ctx, sid, data)

	if err == ErrSessionNotFound {
		return nil, ErrSessionNotFound
	}

	if err != nil {
		err = fmt.Errorf("session.AddAttributes() AddAttributes unexpected error: %w", err)
		ss.Logger.V(0).Info("session.AddAttributes() AddAttributes unexpected error",
			LogKeySID, sid,
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err)
		return nil, err
	}

	return s, nil
}

func (ss *sessionService) RemoveAttributes(ctx context.Context, sid string, keys ...string) (*Session, error) {
	ss.Logger.V(0).Info("session.RemoveAttributes() started", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIDKey))
	defer ss.Logger.V(0).Info("session.RemoveAttributes() finished", LogKeySID, sid, LogKeyRQID, ctx.Value(ss.CtxReqIDKey))

	if len(keys) == 0 {
		err := fmt.Errorf("session.RemoveAttributes() no attributes to remove")
		ss.Logger.V(0).Info(
			"session.RemoveAttributes() no attributes to remove",
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err,
			LogKeySID, sid,
		)
		return nil, err
	}

	s, err := ss.SStore.RemoveAttributes(ctx, sid, keys...)
	if err == ErrSessionNotFound {
		return nil, ErrSessionNotFound
	}

	if err != nil {
		err = fmt.Errorf("session.RemoveAttributes() RemoveAttributes unexpected error: %w", err)
		ss.Logger.V(0).Info("session.RemoveAttributes() RemoveAttributes unexpected error",
			LogKeySID, sid,
			LogKeyRQID, ctx.Value(ss.CtxReqIDKey),
			LogKeyDebugError, err)
		return nil, err
	}

	return s, nil
}

func parseAttrs(keyAndValues ...interface{}) (map[string]interface{}, error) {
	if len(keyAndValues)%2 != 0 {
		return nil, fmt.Errorf("expected even count of key and values, got: %v", len(keyAndValues))
	}

	data := map[string]interface{}{}

	for i := 0; i < len(keyAndValues); i += 2 {
		k, ok := keyAndValues[i].(string)
		if !ok {
			return nil, fmt.Errorf("can't convert key of type: %T to string", keyAndValues[i])
		}
		data[k] = keyAndValues[i+1]
	}
	return data, nil
}
