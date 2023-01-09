package session

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

// Session is representation of session in terms of current module.
// Data - should be used to store any data within a session.
//
// IdleTimeout and LastAccessedAt will be used to expire session based
// on user activity within the session.
//
// AbsTimeout and CreatedAt will be used to expire session based
// on full session lifetime.
//
// UID is supposed to store user identity who session belongs to.
//
// Anonym is supposed to use during authentication process.
type Session struct {
	ID   string
	Data map[string]interface{}
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

// CookieConf contains cookie parameters
type CookieConf struct {
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
	MaxAge   int
	SameSite SameSite
}

// Conf contains session parameters
type Conf struct {
	IdleTimeout time.Duration
	AbsTimout   time.Duration
}

// CtxKey type alias for session data attributes keys
type CtxKey string

var (
	base32enc = base32.StdEncoding.WithPadding(base32.NoPadding)
	keyLen    = 32
)

func DefaultCookieConf() CookieConf {
	return CookieConf{
		Secure:   true,
		HTTPOnly: true,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		SameSite: SameSiteStrictMode,
	}
}

func DefaultSessionConf() Conf {
	return Conf{
		AbsTimout:   24 * 7 * time.Hour,
		IdleTimeout: 24 * time.Hour,
	}
}

// NewSession return new session with default configuration
// IdleTimeout = 24h
// AbsTimeout = 7d
// Anonym = true
// Active = true
// Opts: Secure, HTTPOnly, Strict
func NewSession() (Session, error) {
	id, err := generateSessionID()
	if err != nil {
		return Session{}, err
	}
	s := Session{
		ID:          id,
		Data:        make(map[string]interface{}),
		Opts:        DefaultCookieConf(),
		Anonym:      true,
		Active:      true,
		IdleTimeout: DefaultSessionConf().IdleTimeout,
		AbsTimeout:  DefaultSessionConf().AbsTimout,
	}
	return s, nil
}

// WithUserID add user identity to the session
func (s *Session) WithUserID(uid string) {
	s.UID = uid
	s.Anonym = false
}

// WithCookieConf add cookie to the session
func (s *Session) WithCookieConf(cc CookieConf) {
	s.Opts = cc
}

// WithSessionConf configure session timeouts
func (s *Session) WithSessionConf(sc Conf) {
	s.IdleTimeout = sc.IdleTimeout
	s.AbsTimeout = sc.AbsTimout
}

// IsExpired check if session is expired
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

// ValidateSessionID validate session id format
func ValidateSessionID(sid string) error {
	if len(sid) != base32enc.EncodedLen(keyLen) {
		return fmt.Errorf("error validating session: wrong session id length")
	}
	_, err := base32enc.DecodeString(sid)
	if err != nil {
		return fmt.Errorf("error validating session: %w", err)
	}
	return nil
}

func (s *Session) WithAttributes(attrs map[string]interface{}) {
	for k, v := range attrs {
		s.AddAttribute(k, v)
	}
}

// AddAttribute add a new attribute to the session
func (s *Session) AddAttribute(k string, v interface{}) {
	s.Data[k] = v
}

// GetAttribute return a value from the session
// It return nill and false if attribute doesn't exists
func (s *Session) GetAttribute(k string) (interface{}, bool) {
	v, ok := s.Data[k]
	return v, ok
}

// GetString return a value as string from the session by key
// If value isn't a string, it won't be converted
func (s *Session) GetString(k string) (string, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return "", ok
	}

	switch cv := v.(type) {
	case string:
		return cv, true
	default:
		return "", false
	}
}

// GetInt return a value as int from the session by key
// if Value one of (byte/int8, int16, int32, int64, int)
// return (0, false) otherwise
func (s *Session) GetInt(k string) (int, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return 0, ok
	}

	switch cv := v.(type) {
	case byte:
		return int(cv), true
	case int8:
		return int(cv), true
	case int16:
		return int(cv), true
	case int32:
		return int(cv), true
	case int64:
		return int(cv), true
	case int:
		return cv, true
	default:
		return 0, false
	}
}

// GetInt64 return a value as int64 from the session by key
// if Value one of (byte/int8, int16, int32, int64, int)
// return (0, false) otherwise
func (s *Session) GetInt64(k string) (int64, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return 0, ok
	}

	switch cv := v.(type) {
	case byte:
		return int64(cv), true
	case int8:
		return int64(cv), true
	case int16:
		return int64(cv), true
	case int32:
		return int64(cv), true
	case int:
		return int64(cv), true
	case int64:
		return cv, true
	default:
		return 0, false
	}
}

// GetFloat32 return a value as float32 from the session by key
// if Value is float32
// return (0.0, false) otherwise
func (s *Session) GetFloat32(k string) (float32, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return 0, ok
	}

	switch cv := v.(type) {
	case float32:
		return cv, true
	default:
		return 0.0, false
	}
}

// GetFloat64 return a value as float64 from the session by key
// if Value is one of (float32, float64)
// return (0.0, false) otherwise
func (s *Session) GetFloat64(k string) (float64, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return 0, ok
	}

	switch cv := v.(type) {
	case float32:
		return float64(cv), true
	case float64:
		return cv, true
	default:
		return 0.0, false
	}
}

// GetBool return a value as bool from the session by key
// if Value is bool
// return (false, false) otherwise
func (s *Session) GetBool(k string) (v, ok bool) {
	vl, ok := s.GetAttribute(k)
	if !ok {
		return false, ok
	}

	switch cv := vl.(type) {
	case bool:
		return cv, true
	default:
		return false, false
	}
}

// GetTime return a value as time.Time from the session by key
// if Value is time.Time
// return (time.Time{}, false) otherwise
func (s *Session) GetTime(k string) (time.Time, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return time.Time{}, ok
	}

	switch cv := v.(type) {
	case time.Time:
		return cv, true
	default:
		return time.Time{}, false
	}
}

// GetSlice return a value as []interface{} from the session by key
// if Value is slice
// return (nil, false) otherwise
func (s *Session) GetSlice(k string) ([]interface{}, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return nil, ok
	}

	if reflect.TypeOf(v).Kind() != reflect.Slice {
		return nil, false
	}

	cv, _ := v.([]interface{})
	return cv, true
}

// GetInt32Slice return a value as []int32 from the session by key
// if Value is []int32
// (it won't be converted to []int32, if slice has type of another int)
// return (nil, false) otherwise
func (s *Session) GetInt32Slice(k string) ([]int32, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return nil, ok
	}

	switch cv := v.(type) {
	case []int32:
		return cv, true
	default:
		return nil, false
	}
}

// GetInt64Slice return a value as []int64 from the session by key
// if Value is []int64
// (it won't be converted to []int64, if slice has type of another int)
// return (nil, false) otherwise
func (s *Session) GetInt64Slice(k string) ([]int64, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return nil, ok
	}

	switch cv := v.(type) {
	case []int64:
		return cv, true
	default:
		return nil, false
	}
}

// GetFloat32Slice return a value as []float32 from the session by key
// if Value is []float32
// return (nil, false) otherwise
func (s *Session) GetFloat32Slice(k string) ([]float32, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return nil, ok
	}

	switch cv := v.(type) {
	case []float32:
		return cv, true
	default:
		return nil, false
	}
}

// GetFloat64Slice return a value as []float64 from the session by key
// if Value is []float64
// (it won't be converted to []float64, if slice has type float32)
// return (nil, false) otherwise
func (s *Session) GetFloat64Slice(k string) ([]float64, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return nil, ok
	}

	switch cv := v.(type) {
	case []float64:
		return cv, true
	default:
		return nil, false
	}
}

// GetStringSlice return a value as []string from the session by key
// if Value is []string
// return (nil, false) otherwise
func (s *Session) GetStringSlice(k string) ([]string, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return nil, ok
	}

	switch cv := v.(type) {
	case []string:
		return cv, true
	default:
		return nil, false
	}
}

// GetBoolSlice return a value as []bool from the session by key
// if Value is []bool
// return (nil, false) otherwise
func (s *Session) GetBoolSlice(k string) ([]bool, bool) {
	v, ok := s.GetAttribute(k)
	if !ok {
		return nil, ok
	}

	switch cv := v.(type) {
	case []bool:
		return cv, true
	default:
		return nil, false
	}
}

// GetTimeSlice return a value as []time.Time from the session by key
// if Value is []time.Time
// return (nil, false) otherwise
func (s *Session) GetTimeSlice(k string) ([]time.Time, bool) {
	mapV, ok := s.GetAttribute(k)
	if !ok {
		return nil, ok
	}

	switch cv := mapV.(type) {
	case []time.Time:
		return cv, true
	default:
		return nil, false
	}
}

// GetStruct convert map or struct from the session by key
// and write result to out
// return (false otherwise
//
// to convert value from session to a struct is using
// https://pkg.go.dev/github.com/mitchellh/mapstructure underneath
// so, other its features like tags can be used here
func (s *Session) GetStruct(k string, out interface{}) bool {
	v, ok := s.GetAttribute(k)
	if !ok {
		return false
	}
	err := mapstructure.Decode(v, out)
	return err == nil
}

// GetStructWithDecoder convert map or struct from the session by key
// with custom mapstructure.Decoder (https://pkg.go.dev/github.com/mitchellh/mapstructure)
// return false otherwise
func (s *Session) GetStructWithDecoder(k string, decoder *mapstructure.Decoder) bool {
	v, ok := s.GetAttribute(k)
	if !ok {
		return false
	}
	err := decoder.Decode(v)
	return err == nil
}

func generateSessionID() (string, error) {
	id := make([]byte, keyLen)
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		return "", fmt.Errorf("error generating session id: %w", err)
	}
	return base32enc.EncodeToString(id), nil
}
