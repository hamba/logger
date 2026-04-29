// Package ctx implements log context convenience functions.
package ctx

import (
	"time"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/field"
)

// Str returns a string context field.
//
// Deprecated: use [field.Str] instead.
//
//go:fix inline
func Str(k, s string) logger.Field {
	return field.Str(k, s)
}

// Strs returns a string slice context field.
//
// Deprecated: use [field.Strs] instead.
//
//go:fix inline
func Strs(k string, s []string) logger.Field {
	return field.Strs(k, s)
}

// Bytes returns a byte slice context field.
//
// Deprecated: use [field.Bytes] instead.
//
//go:fix inline
func Bytes(k string, p []byte) logger.Field {
	return field.Bytes(k, p)
}

// Bool returns a boolean context field.
//
// Deprecated: use [field.Bool] instead.
//
//go:fix inline
func Bool(k string, b bool) logger.Field {
	return field.Bool(k, b)
}

// Int returns an int context field.
//
// Deprecated: use [field.Int] instead.
//
//go:fix inline
func Int(k string, i int) logger.Field {
	return field.Int(k, i)
}

// Ints returns an int slice context field.
//
// Deprecated: use [field.Ints] instead.
//
//go:fix inline
func Ints(k string, a []int) logger.Field {
	return field.Ints(k, a)
}

// Int8 returns an int8 context field.
//
// Deprecated: use [field.Int8] instead.
//
//go:fix inline
func Int8(k string, i int8) logger.Field {
	return field.Int8(k, i)
}

// Int16 returns an int16 context field.
//
// Deprecated: use [field.Int16] instead.
//
//go:fix inline
func Int16(k string, i int16) logger.Field {
	return field.Int16(k, i)
}

// Int32 returns an int32 context field.
//
// Deprecated: use [field.Int32] instead.
//
//go:fix inline
func Int32(k string, i int32) logger.Field {
	return field.Int32(k, i)
}

// Int64 returns an int64 context field.
//
// Deprecated: use [field.Int64] instead.
//
//go:fix inline
func Int64(k string, i int64) logger.Field {
	return field.Int64(k, i)
}

// Uint returns a uint context field.
//
// Deprecated: use [field.Uint] instead.
//
//go:fix inline
func Uint(k string, i uint) logger.Field {
	return field.Uint(k, i)
}

// Uint8 returns a uint8 context field.
//
// Deprecated: use [field.Uint8] instead.
//
//go:fix inline
func Uint8(k string, i uint8) logger.Field {
	return field.Uint8(k, i)
}

// Uint16 returns a uint16 context field.
//
// Deprecated: use [field.Uint16] instead.
//
//go:fix inline
func Uint16(k string, i uint16) logger.Field {
	return field.Uint16(k, i)
}

// Uint32 returns a uint32 context field.
//
// Deprecated: use [field.Uint32] instead.
//
//go:fix inline
func Uint32(k string, i uint32) logger.Field {
	return field.Uint32(k, i)
}

// Uint64 returns a uint64 context field.
//
// Deprecated: use [field.Uint64] instead.
//
//go:fix inline
func Uint64(k string, i uint64) logger.Field {
	return field.Uint64(k, i)
}

// Float32 returns a float32 context field.
//
// Deprecated: use [field.Float32] instead.
//
//go:fix inline
func Float32(k string, f float32) logger.Field {
	return field.Float32(k, f)
}

// Float64 returns a float64 context field.
//
// Deprecated: use [field.Float64] instead.
//
//go:fix inline
func Float64(k string, f float64) logger.Field {
	return field.Float64(k, f)
}

// Error returns an error context field.
//
// Deprecated: use [field.Error] instead.
//
//go:fix inline
func Error(k string, err error) logger.Field {
	return field.Error(k, err)
}

// Err returns an error context field with the key set to "error".
//
// Deprecated: use [field.Err] instead.
//
//go:fix inline
func Err(err error) logger.Field {
	return field.Err(err)
}

// Stack return a stack string context field.
//
// Deprecated: use [field.Stack] instead.
//
//go:fix inline
func Stack(k string) logger.Field {
	return field.Stack(k)
}

// Caller returns a caller string context field.
//
// Deprecated: use [field.Caller] instead.
//
//go:fix inline
func Caller(k string) logger.Field {
	return field.Caller(k)
}

// Time returns a time context field.
//
// Deprecated: use [field.Time] instead.
//
//go:fix inline
func Time(k string, t time.Time) logger.Field {
	return field.Time(k, t)
}

// Duration returns a duration context field.
//
// Deprecated: use [field.Duration] instead.
//
//go:fix inline
func Duration(k string, d time.Duration) logger.Field {
	return field.Duration(k, d)
}

// Interface returns an interface context field.
//
// Deprecated: use [field.Interface] instead.
//
//go:fix inline
func Interface(k string, v any) logger.Field {
	return field.Interface(k, v)
}

// Group returns a field that writes all the given fields inside a named group.
//
// Deprecated: use [field.Group] instead.
//
//go:fix inline
func Group(name string, fields ...logger.Field) logger.Field {
	return field.Group(name, fields...)
}

// Span represents an open telemetry span.
//
// Deprecated: use [field.Span] instead.
//
//go:fix inline
type Span = field.Span

// TraceID returns an open telemetry trace ID context field.
//
// Deprecated: use [field.TraceID] instead.
//
//go:fix inline
func TraceID(k string, span field.Span) logger.Field {
	return field.TraceID(k, span)
}
