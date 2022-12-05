package session_test

import (
	"errors"
	"fmt"
	"reflect"
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

func TestCastIntAttributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as int", "bool", byte(8), true, 0, reflect.TypeOf(int(0))},
		{"int8 as int", "int8", int8(8), true, 8, reflect.TypeOf(int(0))},
		{"int16 as int", "int16", int16(8), true, 8, reflect.TypeOf(int(0))},
		{"int32 as int", "int32", int32(8), true, 8, reflect.TypeOf(int(0))},
		{"int as int", "int", int(8), true, 8, reflect.TypeOf(int(0))},
		{"int64 as int", "int64", int64(8), false, 0, reflect.TypeOf(int(0))},
		{"string as int", "string", "va", false, 0, reflect.TypeOf(int(0))},
		{"bool as int", "bool", true, false, 0, reflect.TypeOf(int(0))},
		{"time.Time as int", "time.Time", time.Now(), false, 0, reflect.TypeOf(int(0))},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetInt(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastInt64Attributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as int64", "bool", byte(8), true, 0, reflect.TypeOf(int64(0))},
		{"int8 as int64", "int8", int8(8), true, 8, reflect.TypeOf(int64(0))},
		{"int16 as int64", "int16", int16(8), true, 8, reflect.TypeOf(int64(0))},
		{"int32 as int64", "int32", int32(8), true, 8, reflect.TypeOf(int64(0))},
		{"int as int64", "int", int(8), true, 8, reflect.TypeOf(int64(0))},
		{"int64 as int64", "int64", int64(8), true, 8, reflect.TypeOf(int64(0))},
		{"string as int64", "string", "va", false, 0, reflect.TypeOf(int64(0))},
		{"bool as int64", "bool", true, false, 0, reflect.TypeOf(int64(0))},
		{"time.Time as int64", "time.Time", time.Now(), false, 0, reflect.TypeOf(int64(0))},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetInt64(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastStringAttributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as string", "bool", byte(8), false, "", reflect.TypeOf("")},
		{"int8 as string", "int8", int8(8), false, "", reflect.TypeOf("")},
		{"int16 as string", "int16", int16(8), false, "", reflect.TypeOf("")},
		{"int32 as string", "int32", int32(8), false, "", reflect.TypeOf("")},
		{"int as string", "int", int(8), false, "", reflect.TypeOf("")},
		{"int64 as string", "int64", int64(8), false, "", reflect.TypeOf("")},
		{"string as string", "string", "va", true, "va", reflect.TypeOf("")},
		{"bool as string", "bool", true, false, "", reflect.TypeOf("")},
		{"time.Time as string", "time.Time", time.Now(), false, "", reflect.TypeOf("")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetString(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastBoolAttributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as bool", "bool", byte(8), false, false, reflect.TypeOf(false)},
		{"int8 as bool", "int8", int8(8), false, false, reflect.TypeOf(false)},
		{"int16 as bool", "int16", int16(8), false, false, reflect.TypeOf(false)},
		{"int32 as bool", "int32", int32(8), false, false, reflect.TypeOf(false)},
		{"int as bool", "int", int(8), false, false, reflect.TypeOf(false)},
		{"int64 as bool", "int64", int64(8), false, false, reflect.TypeOf(false)},
		{"string as bool", "string", "va", false, false, reflect.TypeOf(false)},
		{"bool as bool", "bool", true, true, true, reflect.TypeOf(false)},
		{"time.Time as bool", "time.Time", time.Now(), false, false, reflect.TypeOf(false)},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetBool(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastFloat32Attributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as float32", "bool", byte(8), false, 0.0, reflect.TypeOf(float32(1.1))},
		{"int8 as float32", "int8", int8(8), false, 0.0, reflect.TypeOf(float32(1.1))},
		{"int16 as float32", "int16", int16(8), false, 0.0, reflect.TypeOf(float32(1.1))},
		{"int32 as float32", "int32", int32(8), false, 0.0, reflect.TypeOf(float32(1.1))},
		{"int as float32", "int", int(8), false, 0.0, reflect.TypeOf(float32(1.1))},
		{"int64 as float32", "int64", int64(8), false, 0.0, reflect.TypeOf(float32(1.1))},
		{"string as float32", "string", "va", false, 0.0, reflect.TypeOf(float32(1.1))},
		{"bool as float32", "bool", true, false, 0.0, reflect.TypeOf(float32(1.1))},
		{"time.Time as float32", "time.Time", time.Now(), false, 0.0, reflect.TypeOf(float32(1.1))},
		{"float32 as float32", "float32", float32(1.1), true, 1.1, reflect.TypeOf(float32(1.1))},
		{"float64 as float32", "float64", float64(1.1), false, 0.0, reflect.TypeOf(float32(1.1))},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetFloat32(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastFloat64Attributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as float64", "bool", byte(8), false, 0.0, reflect.TypeOf(float64(1.1))},
		{"int8 as float64", "int8", int8(8), false, 0.0, reflect.TypeOf(float64(1.1))},
		{"int16 as float64", "int16", int16(8), false, 0.0, reflect.TypeOf(float64(1.1))},
		{"int32 as float64", "int32", int32(8), false, 0.0, reflect.TypeOf(float64(1.1))},
		{"int as float64", "int", int(8), false, 0.0, reflect.TypeOf(float64(1.1))},
		{"int64 as float64", "int64", int64(8), false, 0.0, reflect.TypeOf(float64(1.1))},
		{"string as float64", "string", "va", false, 0.0, reflect.TypeOf(float64(1.1))},
		{"bool as float64", "bool", true, false, 0.0, reflect.TypeOf(float64(1.1))},
		{"time.Time as float64", "time.Time", time.Now(), false, 0.0, reflect.TypeOf(float64(1.1))},
		{"float32 as float64", "float32", float32(1.1), true, 1.1, reflect.TypeOf(float64(1.1))},
		{"float64 as float64", "float64", float64(1.1), true, 0.0, reflect.TypeOf(float64(1.1))},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetFloat64(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastTimeAttributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as time.Time", "bool", byte(8), false, time.Time{}, reflect.TypeOf(time.Time{})},
		{"int8 as time.Time", "int8", int8(8), false, time.Time{}, reflect.TypeOf(time.Time{})},
		{"int16 as time.Time", "int16", int16(8), false, time.Time{}, reflect.TypeOf(time.Time{})},
		{"int32 as time.Time", "int32", int32(8), false, time.Time{}, reflect.TypeOf(time.Time{})},
		{"int as time.Time", "int", int(8), false, time.Time{}, reflect.TypeOf(time.Time{})},
		{"int64 as time.Time", "int64", int64(8), false, time.Time{}, reflect.TypeOf(time.Time{})},
		{"string as time.Time", "string", "va", false, time.Time{}, reflect.TypeOf(time.Time{})},
		{"bool as time.Time", "bool", true, false, time.Time{}, reflect.TypeOf(time.Time{})},
		{"time.Time as time.Time", "time.Time", time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC), true, time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC), reflect.TypeOf(time.Time{})},
		{"float32 as time.Time", "float32", float32(1.1), false, time.Time{}, reflect.TypeOf(time.Time{})},
		{"float64 as time.Time", "float64", float64(1.1), false, time.Time{}, reflect.TypeOf(time.Time{})},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetTime(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastSliceInterfaceAttributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as []interface{}", "bool", byte(8), false, nil, reflect.TypeOf([]interface{}{})},
		{"int8 as []interface{}", "int8", int8(8), false, nil, reflect.TypeOf([]interface{}{})},
		{"int16 as []interface{}", "int16", int16(8), false, nil, reflect.TypeOf([]interface{}{})},
		{"int32 as []interface{}", "int32", int32(8), false, nil, reflect.TypeOf([]interface{}{})},
		{"int as []interface{}", "int", int(8), false, nil, reflect.TypeOf([]interface{}{})},
		{"int64 as []interface{}", "int64", int64(8), false, nil, reflect.TypeOf([]interface{}{})},
		{"string as []interface{}", "string", "va", false, nil, reflect.TypeOf([]interface{}{})},
		{"bool as []interface{}", "bool", true, false, nil, reflect.TypeOf([]interface{}{})},
		{"time.Time as []interface{}", "time.Time", time.Now(), false, nil, reflect.TypeOf([]interface{}{})},
		{"float32 as []interface{}", "float32", float32(1.1), false, nil, reflect.TypeOf([]interface{}{})},
		{"float64 as []interface{}", "float64", float64(1.1), false, nil, reflect.TypeOf([]interface{}{})},
		{"[]interface{} as []interface{}", "[]interface{}", []interface{}{1, 2, 3}, true, []interface{}{1, 2, 3}, reflect.TypeOf([]interface{}{})},
		{"[]int as []interface{}", "[]int", []int{1, 2, 3}, true, []interface{}{1, 2, 3}, reflect.TypeOf([]interface{}{})},
		{"[]string as []interface{}", "[]string", []string{"1", "2", "3"}, true, []interface{}{"1", "2", "3"}, reflect.TypeOf([]interface{}{})},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetSlice(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastSliceIntAttributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as []int32{}", "bool", byte(8), false, nil, reflect.TypeOf([]int32{})},
		{"int8 as []int32{}", "int8", int8(8), false, nil, reflect.TypeOf([]int32{})},
		{"int16 as []int32{}", "int16", int16(8), false, nil, reflect.TypeOf([]int32{})},
		{"int32 as []int32{}", "int32", int32(8), false, nil, reflect.TypeOf([]int32{})},
		{"int as []int32{}", "int", int(8), false, nil, reflect.TypeOf([]int32{})},
		{"int64 as []int32{}", "int64", int64(8), false, nil, reflect.TypeOf([]int32{})},
		{"string as []int32{}", "string", "va", false, nil, reflect.TypeOf([]int32{})},
		{"bool as []int32{}", "bool", true, false, nil, reflect.TypeOf([]int32{})},
		{"time.Time as []int32{}", "time.Time", time.Now(), false, nil, reflect.TypeOf([]int32{})},
		{"float32 as []int32{}", "float32", float32(1.1), false, nil, reflect.TypeOf([]int32{})},
		{"float64 as []int32{}", "float64", float64(1.1), false, nil, reflect.TypeOf([]int32{})},
		{"[]interface{} as []int32{}", "[]interface{}", []interface{}{1, 2, 3}, false, nil, reflect.TypeOf([]int32{})},
		{"[]int as []int{}", "[]int32", []int{1, 2, 3}, false, nil, reflect.TypeOf([]int32{})},
		{"[]int(int64 actually) as []int", "[]int32", []int{10000000000, 10000000001, 10000000002}, false, nil, reflect.TypeOf([]int32{})},
		{"[]32int as []int{}", "[]int32", []int32{1, 2, 3}, true, []int32{1, 2, 3}, reflect.TypeOf([]int32{})},
		{"[]64int as []int{}", "[]int32", []int64{1, 2, 3}, false, nil, reflect.TypeOf([]int32{})},
		{"[]string as []int32{}", "[]string", []string{"1", "2", "3"}, false, nil, reflect.TypeOf([]int32{})},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetInt32Slice(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastSliceInt64Attributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as []int64{}", "bool", byte(8), false, nil, reflect.TypeOf([]int64{})},
		{"int8 as []int64{}", "int8", int8(8), false, nil, reflect.TypeOf([]int64{})},
		{"int16 as []int64{}", "int16", int16(8), false, nil, reflect.TypeOf([]int64{})},
		{"int32 as []int64{}", "int32", int32(8), false, nil, reflect.TypeOf([]int64{})},
		{"int as []int64{}", "int", int(8), false, nil, reflect.TypeOf([]int64{})},
		{"int64 as []int64{}", "int64", int64(8), false, nil, reflect.TypeOf([]int64{})},
		{"string as []int64{}", "string", "va", false, nil, reflect.TypeOf([]int64{})},
		{"bool as []int64{}", "bool", true, false, nil, reflect.TypeOf([]int64{})},
		{"time.Time as []int64{}", "time.Time", time.Now(), false, nil, reflect.TypeOf([]int64{})},
		{"float32 as []int64{}", "float32", float32(1.1), false, nil, reflect.TypeOf([]int64{})},
		{"float64 as []int64{}", "float64", float64(1.1), false, nil, reflect.TypeOf([]int64{})},
		{"[]interface{} as []int64{}", "[]interface{}", []interface{}{1, 2, 3}, false, nil, reflect.TypeOf([]int64{})},
		{"[]int as []int{}", "[]int64", []int{1, 2, 3}, false, nil, reflect.TypeOf([]int64{})},
		{"[]int(int64 actually) as []int", "[]int64", []int{10000000000, 10000000001, 10000000002}, false, nil, reflect.TypeOf([]int64{})},
		{"[]32int as []int{}", "[]int64", []int32{1, 2, 3}, false, nil, reflect.TypeOf([]int64{})},
		{"[]64int as []int{}", "[]int64", []int64{1, 2, 3}, true, []int64{1, 2, 3}, reflect.TypeOf([]int64{})},
		{"[]string as []int64{}", "[]string", []string{"1", "2", "3"}, false, nil, reflect.TypeOf([]int64{})},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetInt64Slice(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastSliceFloat32Attributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as []int64{}", "bool", byte(8), false, nil, reflect.TypeOf([]float32{})},
		{"int8 as []float32{}", "int8", int8(8), false, nil, reflect.TypeOf([]float32{})},
		{"int16 as []float32{}", "int16", int16(8), false, nil, reflect.TypeOf([]float32{})},
		{"int32 as []float32{}", "int32", int32(8), false, nil, reflect.TypeOf([]float32{})},
		{"int as []float32{}", "int", int(8), false, nil, reflect.TypeOf([]float32{})},
		{"int64 as []float32{}", "int64", int64(8), false, nil, reflect.TypeOf([]float32{})},
		{"string as []float32{}", "string", "va", false, nil, reflect.TypeOf([]float32{})},
		{"bool as []float32{}", "bool", true, false, nil, reflect.TypeOf([]float32{})},
		{"time.Time as []float32{}", "time.Time", time.Now(), false, nil, reflect.TypeOf([]float32{})},
		{"float32 as []float32{}", "float32", float32(1.1), false, nil, reflect.TypeOf([]float32{})},
		{"float64 as []float32{}", "float64", float64(1.1), false, nil, reflect.TypeOf([]float32{})},
		{"[]float32 as []float32{}", "[]float32", []float32{1, 2, 3}, true, []float32{1, 2, 3}, reflect.TypeOf([]float32{})},
		{"[]float64 as []float32{}", "[]int64", []float64{1, 2, 3}, false, nil, reflect.TypeOf([]float32{})},
		{"[]string as []float32{}", "[]string", []string{"1", "2", "3"}, false, nil, reflect.TypeOf([]float32{})},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetFloat32Slice(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastSliceFloat64Attributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as []float64{}", "bool", byte(8), false, nil, reflect.TypeOf([]float64{})},
		{"int8 as []float64{}", "int8", int8(8), false, nil, reflect.TypeOf([]float64{})},
		{"int16 as []float64{}", "int16", int16(8), false, nil, reflect.TypeOf([]float64{})},
		{"int32 as []float64{}", "int32", int32(8), false, nil, reflect.TypeOf([]float64{})},
		{"int as []float64{}", "int", int(8), false, nil, reflect.TypeOf([]float64{})},
		{"int64 as []float64{}", "int64", int64(8), false, nil, reflect.TypeOf([]float64{})},
		{"string as []float64{}", "string", "va", false, nil, reflect.TypeOf([]float64{})},
		{"bool as []float64{}", "bool", true, false, nil, reflect.TypeOf([]float64{})},
		{"time.Time as []float64{}", "time.Time", time.Now(), false, nil, reflect.TypeOf([]float64{})},
		{"float32 as []float64{}", "float32", float32(1.1), false, nil, reflect.TypeOf([]float64{})},
		{"float64 as []float64{}", "float64", float64(1.1), false, nil, reflect.TypeOf([]float64{})},
		{"[]float32 as []float64{}", "[]float32", []float32{1, 2, 3}, false, nil, reflect.TypeOf([]float64{})},
		{"[]float64 as []float64{}", "[]int64", []float64{1, 2, 3}, true, []float64{1, 2, 3}, reflect.TypeOf([]float64{})},
		{"[]string as []float64{}", "[]string", []string{"1", "2", "3"}, false, nil, reflect.TypeOf([]float64{})},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetFloat64Slice(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastSliceStringAttributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as []string{}", "bool", byte(8), false, nil, reflect.TypeOf([]string{})},
		{"int8 as []string{}", "int8", int8(8), false, nil, reflect.TypeOf([]string{})},
		{"int16 as []string{}", "int16", int16(8), false, nil, reflect.TypeOf([]string{})},
		{"int32 as []string{}", "int32", int32(8), false, nil, reflect.TypeOf([]string{})},
		{"int as []string{}", "int", int(8), false, nil, reflect.TypeOf([]string{})},
		{"int64 as []string{}", "int64", int64(8), false, nil, reflect.TypeOf([]string{})},
		{"string as []string{}", "string", "va", false, nil, reflect.TypeOf([]string{})},
		{"bool as []string{}", "bool", true, false, nil, reflect.TypeOf([]string{})},
		{"time.Time as []string{}", "time.Time", time.Now(), false, nil, reflect.TypeOf([]string{})},
		{"float32 as []string{}", "float32", float32(1.1), false, nil, reflect.TypeOf([]string{})},
		{"float64 as []string{}", "float64", float64(1.1), false, nil, reflect.TypeOf([]string{})},
		{"[]float32 as []string{}", "[]float32", []float32{1, 2, 3}, false, nil, reflect.TypeOf([]string{})},
		{"[]string as []string{}", "[]string", []string{"1", "2", "3"}, true, []string{"1", "2", "3"}, reflect.TypeOf([]string{})},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetStringSlice(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastSliceBoolAttributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as []bool{}", "bool", byte(8), false, nil, reflect.TypeOf([]bool{})},
		{"int8 as []bool{}", "int8", int8(8), false, nil, reflect.TypeOf([]bool{})},
		{"int16 as []bool{}", "int16", int16(8), false, nil, reflect.TypeOf([]bool{})},
		{"int32 as []bool{}", "int32", int32(8), false, nil, reflect.TypeOf([]bool{})},
		{"int as []bool{}", "int", int(8), false, nil, reflect.TypeOf([]bool{})},
		{"int64 as []bool{}", "int64", int64(8), false, nil, reflect.TypeOf([]bool{})},
		{"string as []bool{}", "string", "va", false, nil, reflect.TypeOf([]bool{})},
		{"bool as []bool{}", "bool", true, false, nil, reflect.TypeOf([]bool{})},
		{"time.Time as []bool{}", "time.Time", time.Now(), false, nil, reflect.TypeOf([]bool{})},
		{"float32 as []bool{}", "float32", float32(1.1), false, nil, reflect.TypeOf([]bool{})},
		{"float64 as []bool{}", "float64", float64(1.1), false, nil, reflect.TypeOf([]bool{})},
		{"[]bool as []bool{}", "[]bool", []bool{true, false}, true, []bool{true, false}, reflect.TypeOf([]bool{})},
		{"[]string as []bool{}", "[]string", []string{"1", "2", "3"}, false, nil, reflect.TypeOf([]bool{})},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetBoolSlice(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}

func TestCastSliceTimeAttributes(t *testing.T) {
	tt := []struct {
		name          string
		key           string
		value         interface{}
		expectedOk    bool
		expectedValue interface{}
		expectedType  reflect.Type
	}{
		{"byte as []time.Time{}", "byte", byte(8), false, nil, reflect.TypeOf([]time.Time{})},
		{"int8 as []time.Time{}", "int8", int8(8), false, nil, reflect.TypeOf([]time.Time{})},
		{"int16 as []time.Time{}", "int16", int16(8), false, nil, reflect.TypeOf([]time.Time{})},
		{"int32 as []time.Time{}", "int32", int32(8), false, nil, reflect.TypeOf([]time.Time{})},
		{"int as []time.Time{}", "int", int(8), false, nil, reflect.TypeOf([]time.Time{})},
		{"int64 as []time.Time{}", "int64", int64(8), false, nil, reflect.TypeOf([]time.Time{})},
		{"string as []time.Time{}", "string", "va", false, nil, reflect.TypeOf([]time.Time{})},
		{"bool as []time.Time{}", "bool", true, false, nil, reflect.TypeOf([]time.Time{})},
		{"time.Time as []time.Time{}", "time.Time", time.Now(), false, nil, reflect.TypeOf([]time.Time{})},
		{"float32 as []time.Time{}", "float32", float32(1.1), false, nil, reflect.TypeOf([]time.Time{})},
		{"float64 as []time.Time{}", "float64", float64(1.1), false, nil, reflect.TypeOf([]time.Time{})},
		{"[]time.Time as []time.Time{}", "[]bool", []time.Time{time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC)}, true, []time.Time{time.Date(2022, 1, 2, 3, 4, 5, 6, time.UTC)}, reflect.TypeOf([]time.Time{})},
		{"[]string as []time.Time{}", "[]string", []string{"1", "2", "3"}, false, nil, reflect.TypeOf([]time.Time{})},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := session.NewSession()
			assert.Nil(t, err)
			s.AddAttribute(tc.key, tc.value)
			v, ok := s.GetTimeSlice(tc.key)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedType, reflect.TypeOf(v))
		})
	}
}
