package mongo

import (
	"context"
	"fmt"
	"reflect"

	"time"

	"github.com/asstart/go-session"
	"github.com/go-logr/logr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStore struct {
	Collecction    *mongo.Collection
	Logger         logr.Logger
	CustomRegistry *bsoncodec.Registry
	CtxReqIDKey    interface{}
}

func NewMongoStore(c *mongo.Collection, l logr.Logger, reqIDKey interface{}) session.Store {
	return &mongoStore{
		Collecction:    c,
		Logger:         l,
		CustomRegistry: getCustomRegisry(),
		CtxReqIDKey:    reqIDKey,
	}
}

type mngSession struct {
	ID             primitive.ObjectID     `bson:"_id"`
	SID            string                 `bson:"sid"`
	Data           map[string]interface{} `bson:"data,omitempty"`
	Opts           mngCookieConf          `bson:"opts"`
	Anonym         bool                   `bson:"anonym"`
	Active         bool                   `bson:"active"`
	UID            string                 `bson:"uid"`
	IdleTimeout    time.Duration          `bson:"idle_timeout"`
	AbsTimeout     time.Duration          `bson:"abs_timeout"`
	LastAccessedAt time.Time              `bson:"last_accessed_at"`
	CreatedAt      time.Time              `bson:"created_at"`
}

type mngCookieConf struct {
	Path     string           `bson:"path"`
	Domain   string           `bson:"domain"`
	Secure   bool             `bson:"secure"`
	HTTPOnly bool             `bson:"http_only"`
	MaxAge   int              `bson:"max_age"`
	SameSite session.SameSite `bson:"same_site"`
}

func fromMngCookieConf(c mngCookieConf) session.CookieConf {
	return session.CookieConf{
		Path:     c.Path,
		Domain:   c.Domain,
		Secure:   c.Secure,
		HTTPOnly: c.HTTPOnly,
		MaxAge:   c.MaxAge,
		SameSite: c.SameSite,
	}
}

func fromMngSession(s *mngSession) session.Session {

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

func (ms *mongoStore) Save(ctx context.Context, s *session.Session) (*session.Session, error) {

	ms.Logger.V(0).Info("session.mongo.Save() started", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
	defer ms.Logger.V(0).Info("session.mongo.Save() finished", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))

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
				{"http_only", s.Opts.HTTPOnly},
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

	var updS mngSession
	sr := ms.Collecction.FindOneAndUpdate(ctx, f, o, opts)
	err := decodeWithRegistry(ms.CustomRegistry, sr, &updS)

	if err == mongo.ErrNoDocuments {
		ms.Logger.V(0).Info("session.mong.Save() session not found", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
		return nil, session.ErrSessionNotFound
	}

	if err != nil {
		err = fmt.Errorf("session.mong.Save() FindOneAndUpdate() unexpected error: %w", err)
		ms.Logger.V(0).Info("session.mong.Save() session not found",
			session.LogKeySID, s.ID,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey),
			session.LogKeyDebugError, err)

		return nil, err
	}

	r := fromMngSession(&updS)

	return &r, nil
}

func (ms *mongoStore) Invalidate(ctx context.Context, sid string) error {
	ms.Logger.V(0).Info("session.mongo.Invalidate() started", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
	defer ms.Logger.V(0).Info("session.mongo.Invalidate() finished", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))

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
		err = fmt.Errorf("session.mongo.Invalidate error: %w", err)
		ms.Logger.V(0).Info(
			"session.mongo.Invalidate() error",
			session.LogKeySID, sid,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey),
			session.LogKeyDebugError, err)
		return err
	}
	return nil
}

func (ms *mongoStore) Update(ctx context.Context, s *session.Session) (*session.Session, error) {
	ms.Logger.V(0).Info("session.mongo.Update() started", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
	defer ms.Logger.V(0).Info("session.mongo.Update() finished", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
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
				{"http_only", s.Opts.HTTPOnly},
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

	var updS mngSession
	sr := ms.Collecction.FindOneAndUpdate(ctx, f, obj, opts)
	err := decodeWithRegistry(ms.CustomRegistry, sr, &updS)

	if err == mongo.ErrNoDocuments {
		ms.Logger.V(0).Info("session.mong.Update() session not found", session.LogKeySID, s.ID, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
		return nil, session.ErrSessionNotFound
	}

	if err != nil {
		err = fmt.Errorf("session.mong.Update() FindOneAndUpdate() unexpected error: %w", err)
		ms.Logger.V(0).Info("session.mong.Update() session not found",
			session.LogKeySID, s.ID,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey),
			session.LogKeyDebugError, err)

		return nil, err
	}

	r := fromMngSession(&updS)

	return &r, nil
}

func (ms *mongoStore) Load(ctx context.Context, sid string) (*session.Session, error) {
	ms.Logger.V(0).Info("session.mongo.Load() started", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
	defer ms.Logger.V(0).Info("session.mongo.Load() finished", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))

	f := bson.D{
		{"sid", sid},
	}

	upd := bson.D{
		{"$currentDate", bson.D{{"last_accessed_at", true}}},
	}

	opts := options.FindOneAndUpdate()
	opts = opts.SetReturnDocument(options.After)

	s := mngSession{}
	sr := ms.Collecction.FindOneAndUpdate(ctx, f, upd, opts)
	err := decodeWithRegistry(ms.CustomRegistry, sr, &s)

	if err == mongo.ErrNoDocuments {
		ms.Logger.V(0).Info("session.mong.Load() session not found", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
		return nil, session.ErrSessionNotFound
	}

	if err != nil {
		err = fmt.Errorf("session.mong.Load() FindOneAndUpdate() unexpected error: %w", err)
		ms.Logger.V(0).Info("session.mong.Load() session not found",
			session.LogKeySID, sid,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey),
			session.LogKeyDebugError, err)

		return nil, err
	}

	r := fromMngSession(&s)

	return &r, nil
}

func (ms *mongoStore) AddAttributes(ctx context.Context, sid string, data map[string]interface{}) (*session.Session, error) {
	ms.Logger.V(0).Info("session.mongo.AddAttributes() started", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
	defer ms.Logger.V(0).Info("session.mongo.AddAttributes() finished", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))

	f := bson.M{"sid": sid}
	up := bson.A{
		bson.D{{"$set", bson.D{
			{"data", bson.D{
				{"$mergeObjects", bson.A{
					"$data", data,
				}},
			}},
		}}},
		bson.D{{"$addFields",
			bson.D{
				{"last_accessed_at", "$$NOW"},
			},
		}},
	}
	opt := options.FindOneAndUpdate()
	opt.SetReturnDocument(options.After)

	var s mngSession
	sr := ms.Collecction.FindOneAndUpdate(ctx, f, up, opt)
	err := decodeWithRegistry(ms.CustomRegistry, sr, &s)

	if err == mongo.ErrNoDocuments {
		ms.Logger.V(0).Info("session.mongo.AddAttributes() FindOneAndUpdate() session not found", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
		return nil, session.ErrSessionNotFound
	}

	if err != nil {
		err = fmt.Errorf("session.mongo.AddAttributes() FindOneAndUpdate() unexpected error: %w", err)
		ms.Logger.V(0).Info("session.mongo.AddAttributes() FindOneAndUpdate() unexpected error",
			session.LogKeySID, sid,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey),
			session.LogKeyDebugError, err,
		)
		return nil, err
	}

	r := fromMngSession(&s)
	return &r, nil
}

func (ms *mongoStore) RemoveAttributes(ctx context.Context, sid string, keys ...string) (*session.Session, error) {
	ms.Logger.V(0).Info("session.mongo.RemoveAttributes() started", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
	defer ms.Logger.V(0).Info("session.mongo.RemoveAttributes() finished", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))

	fullkeys := []string{}
	for _, k := range keys {
		fullkeys = append(fullkeys, fmt.Sprintf("data.%v", k))
	}

	f := bson.M{"sid": sid}
	up := bson.A{
		bson.D{{"$unset", fullkeys}},
		bson.D{{"$addFields",
			bson.D{
				{"last_accessed_at", "$$NOW"},
			},
		}},
	}
	opt := options.FindOneAndUpdate()
	opt.SetReturnDocument(options.After)

	var s mngSession
	sr := ms.Collecction.FindOneAndUpdate(ctx, f, up, opt)
	err := decodeWithRegistry(ms.CustomRegistry, sr, &s)

	if err == mongo.ErrNoDocuments {
		ms.Logger.V(0).Info("session.mongo.AddAttributes() FindOneAndUpdate() session not found", session.LogKeySID, sid, session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey))
		return nil, session.ErrSessionNotFound
	}

	if err != nil {
		err = fmt.Errorf("session.mongo.AddAttributes() FindOneAndUpdate() unexpected error: %w", err)
		ms.Logger.V(0).Info("session.mongo.AddAttributes() FindOneAndUpdate() unexpected error",
			session.LogKeySID, sid,
			session.LogKeyRQID, ctx.Value(ms.CtxReqIDKey),
			session.LogKeyDebugError, err,
		)
		return nil, err
	}

	r := fromMngSession(&s)
	return &r, nil
}

func decodeWithRegistry(r *bsoncodec.Registry, sr *mongo.SingleResult, v interface{}) error {
	if sr.Err() != nil {
		return sr.Err()
	}

	raw, err := sr.DecodeBytes()
	if err != nil {
		return err
	}

	return bson.UnmarshalWithRegistry(r, raw, v)
}

func getCustomRegisry() *bsoncodec.Registry {
	rb := bsoncodec.NewRegistryBuilder()

	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)

	rb.RegisterTypeMapEntry(bsontype.DateTime, reflect.TypeOf(time.Time{}))
	rb.RegisterTypeMapEntry(bson.TypeArray, reflect.TypeOf([]interface{}{}))

	return rb.Build()
}
