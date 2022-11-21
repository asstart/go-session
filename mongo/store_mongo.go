package mongo

import (
	"context"
	"fmt"

	"time"

	"github.com/asstart/go-session"
	"github.com/go-logr/logr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStore struct {
	Collecction *mongo.Collection
	Logger      logr.Logger
	CtxReqIdKey interface{}
}

func NewMongoStore(c *mongo.Collection, l logr.Logger, reqIdKey interface{}) session.SessionStore {
	return &mongoStore{
		Collecction: c,
		Logger:      l,
		CtxReqIdKey: reqIdKey,
	}
}

type mngSession struct {
	ID             primitive.ObjectID                 `bson:"_id"`
	SID            string                             `bson:"sid"`
	Data           map[session.SessionKey]interface{} `bson:"data,omitempty"`
	Opts           mngCookieConf                      `bson:"opts"`
	Anonym         bool                               `bson:"anonym"`
	Active         bool                               `bson:"active"`
	UID            string                             `bson:"uid"`
	IdleTimeout    time.Duration                      `bson:"idle_timeout"`
	AbsTimeout     time.Duration                      `bson:"abs_timeout"`
	LastAccessedAt time.Time                          `bson:"last_accessed_at"`
	CreatedAt      time.Time                          `bson:"created_at"`
}

type mngCookieConf struct {
	Path     string           `bson:"path"`
	Domain   string           `bson:"domain"`
	Secure   bool             `bson:"secure"`
	HttpOnly bool             `bson:"http_only"`
	MaxAge   int              `bson:"max_age"`
	SameSite session.SameSite `bson:"same_site"`
}

func toMngCookieConf(c session.CookieConf) mngCookieConf {
	return mngCookieConf{
		Path:     c.Path,
		Domain:   c.Domain,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
		MaxAge:   c.MaxAge,
		SameSite: c.SameSite,
	}
}

func fromMngCookieConf(c mngCookieConf) session.CookieConf {
	return session.CookieConf{
		Path:     c.Path,
		Domain:   c.Domain,
		Secure:   c.Secure,
		HttpOnly: c.HttpOnly,
		MaxAge:   c.MaxAge,
		SameSite: c.SameSite,
	}
}

func toMngSession(s session.Session) (mngSession, error) {

	return mngSession{
		ID:             primitive.NewObjectID(),
		SID:            s.ID,
		Data:           s.Data,
		Opts:           toMngCookieConf(s.Opts),
		Anonym:         s.Anonym,
		Active:         s.Active,
		UID:            s.UID,
		IdleTimeout:    s.IdleTimeout,
		AbsTimeout:     s.AbsTimeout,
		LastAccessedAt: s.LastAccessedAt,
		CreatedAt:      s.CreatedAt,
	}, nil
}

func fromMngSession(s mngSession) session.Session {

	return session.Session{
		ID:             s.SID,
		Data:           s.Data,
		Opts:           fromMngCookieConf(s.Opts),
		Anonym:         s.Anonym,
		Active:         s.Active,
		UID:            s.UID,
		IdleTimeout:    s.IdleTimeout,
		AbsTimeout:     s.AbsTimeout,
		LastAccessedAt: s.LastAccessedAt,
		CreatedAt:      s.CreatedAt,
	}
}

func (ms *mongoStore) Save(ctx context.Context, s session.Session) (session.Session, error) {

	ms.Logger.V(0).Info("session.mongo.Save() started", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey))
	defer ms.Logger.V(0).Info("session.mongo.Save() finished", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey))

	f := bson.D{
		{"_id", primitive.NewObjectID()},
	}
	o := bson.D{
		{"$set", bson.D{
			{"sid", s.ID},
			{"data", s.Data},
			{"opts", bson.D{
				{"path", s.Opts.Path},
				{"domain", s.Opts.Domain},
				{"secure", s.Opts.Secure},
				{"http_only", s.Opts.HttpOnly},
				{"max_age", s.Opts.MaxAge},
				{"same_site", s.Opts.SameSite},
			}},
			{"anonym", s.Anonym},
			{"active", s.Active},
			{"uid", s.UID},
			{"idle_timeout", s.IdleTimeout},
			{"abs_timeout", s.AbsTimeout},
		}},
		{"$currentDate", bson.D{
			{"last_accessed_at", true},
			{"created_at", true},
		}},
	}

	opts := options.FindOneAndUpdate()
	opts = opts.SetUpsert(true)
	opts = opts.SetReturnDocument(options.After)

	sr := ms.Collecction.FindOneAndUpdate(
		ctx,
		f,
		o,
		opts,
	)

	if sr.Err() != nil {
		err := fmt.Errorf("session.mongo.Save() FindOneAndUpdate error: %w", sr.Err())
		ms.Logger.V(0).Info(
			"session.mongo.Save() error",
			session.LogKeySID, s.ID,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey),
			session.LogKeyDebugError, err)
		return session.Session{}, err
	}

	var updS mngSession
	if err := sr.Decode(&updS); err != nil {
		err := fmt.Errorf("session.mongo.Save() Decode result error: %w", err)
		ms.Logger.V(0).Info(
			"session.mongo.Save() error",
			session.LogKeySID, s.ID,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey),
			session.LogKeyDebugError, err)
		return session.Session{}, err
	}

	return fromMngSession(updS), nil
}

func (ms *mongoStore) Invalidate(ctx context.Context, sid string) error {
	ms.Logger.V(0).Info("session.mongo.Invalidate() started", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey))
	defer ms.Logger.V(0).Info("session.mongo.Invalidate() finished", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey))

	f := bson.D{
		{"sid", sid},
	}

	op := bson.D{
		{"$set", bson.D{
			{"active", false},
		}},
		{"$currentDate", bson.D{
			{"last_accessed_at", true},
		}},
	}

	_, err := ms.Collecction.UpdateOne(ctx, f, op)
	if err != nil {
		err := fmt.Errorf("session.mongo.Invalidate error: %w", err)
		ms.Logger.V(0).Info(
			"session.mongo.Invalidate() error",
			session.LogKeySID, sid,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey),
			session.LogKeyDebugError, err)
		return err
	}
	return nil
}

func (ms *mongoStore) Update(ctx context.Context, s session.Session) (session.Session, error) {
	ms.Logger.V(0).Info("session.mongo.Update() started", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey))
	defer ms.Logger.V(0).Info("session.mongo.Update() finished", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey))
	f := bson.D{
		{"sid", s.ID},
	}
	obj := bson.D{
		{"$set", bson.D{
			{"data", s.Data},
			{"opts", bson.D{
				{"path", s.Opts.Path},
				{"domain", s.Opts.Domain},
				{"secure", s.Opts.Secure},
				{"http_only", s.Opts.HttpOnly},
				{"max_age", s.Opts.MaxAge},
				{"same_site", s.Opts.SameSite},
			}},
			{"anonym", s.Anonym},
			{"active", s.Active},
			{"uid", s.UID},
			{"idle_timeout", s.IdleTimeout},
			{"abs_timeout", s.AbsTimeout},
		}},
		{"$currentDate", bson.D{
			{"last_accessed_at", true},
		}},
	}

	opts := options.FindOneAndUpdate()
	opts = opts.SetReturnDocument(options.After)

	sr := ms.Collecction.FindOneAndUpdate(
		ctx,
		f,
		obj,
		opts,
	)

	if sr.Err() == nil {
		err := fmt.Errorf("session.mongo.Update() FindIneAndUpdate error: %w", sr.Err())
		ms.Logger.V(0).Info(
			"session.mongo.Update() error",
			session.LogKeySID, s.ID,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey),
			session.LogKeyDebugError, err)
		return session.Session{}, err
	}

	var updS session.Session
	if err := sr.Decode(updS); err != nil {
		err := fmt.Errorf("session.mongo.Update() Decode error: %w", sr.Err())
		ms.Logger.V(0).Info(
			"session.mongo.Update() error",
			session.LogKeySID, s.ID,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey),
			session.LogKeyDebugError, err)
		return session.Session{}, err
	}

	return updS, nil
}

func (ms *mongoStore) Load(ctx context.Context, sid string) (session.Session, error) {
	ms.Logger.V(0).Info("session.mongo.Load() started", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey))
	defer ms.Logger.V(0).Info("session.mongo.Load() finished", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey))

	f := bson.D{
		{"sid", sid},
	}

	upd := bson.D{
		{"$currentDate", bson.D{{"last_accessed_at", true}}},
	}

	opts := options.FindOneAndUpdate()
	opts = opts.SetReturnDocument(options.After)

	s := mngSession{}

	err := ms.Collecction.FindOneAndUpdate(ctx, f, upd, opts).Decode(&s)

	if err == mongo.ErrNoDocuments {
		ms.Logger.V(0).Info("session.mong.Load() session not found", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey))
		return session.Session{}, session.ErrSessionNotFound
	}

	if err != nil {
		err := fmt.Errorf("session.mong.Load() FindOneAndUpdate() unexpected error: %w", err)
		ms.Logger.V(0).Info("session.mong.Load() session not found",
			session.LogKeySID, sid,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIdKey),
			session.LogKeyDebugError, err)

		return session.Session{}, err
	}

	return fromMngSession(s), nil
}
