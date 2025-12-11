// Package ctx implements log context convenience functions.
package ctx

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/go-stack/stack"
	"github.com/hamba/logger/v2"
	"go.opentelemetry.io/otel/trace"
)

// Str returns a string context field.
func Str(k, s string) logger.Field {
	return func(e *logger.Event) {
		e.AppendString(k, s)
	}
}

// Strs returns a string slice context field.
func Strs(k string, s []string) logger.Field {
	return func(e *logger.Event) {
		e.AppendStrings(k, s)
	}
}

// Bytes returns a byte slice context field.
func Bytes(k string, p []byte) logger.Field {
	return func(e *logger.Event) {
		e.AppendBytes(k, p)
	}
}

// Bool returns a boolean context field.
func Bool(k string, b bool) logger.Field {
	return func(e *logger.Event) {
		e.AppendBool(k, b)
	}
}

// Int returns an int context field.
func Int(k string, i int) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

// Ints returns an int slice context field.
func Ints(k string, a []int) logger.Field {
	return func(e *logger.Event) {
		e.AppendInts(k, a)
	}
}

// Int8 returns an int8 context field.
func Int8(k string, i int8) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

// Int16 returns an int16 context field.
func Int16(k string, i int16) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

// Int32 returns an int32 context field.
func Int32(k string, i int32) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

// Int64 returns an int64 context field.
func Int64(k string, i int64) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, i)
	}
}

// Uint returns a uint context field.
func Uint(k string, i uint) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

// Uint8 returns a uint8 context field.
func Uint8(k string, i uint8) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

// Uint16 returns a uint16 context field.
func Uint16(k string, i uint16) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

// Uint32 returns a uint32 context field.
func Uint32(k string, i uint32) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

// Uint64 returns a uint64 context field.
func Uint64(k string, i uint64) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, i)
	}
}

// Float32 returns a float32 context field.
func Float32(k string, f float32) logger.Field {
	return func(e *logger.Event) {
		e.AppendFloat(k, float64(f))
	}
}

// Float64 returns a float64 context field.
func Float64(k string, f float64) logger.Field {
	return func(e *logger.Event) {
		e.AppendFloat(k, f)
	}
}

// Error returns an error context field.
func Error(k string, err error) logger.Field {
	return func(e *logger.Event) {
		e.AppendString(k, err.Error())
	}
}

// Err returns an error context field with the key set to "error".
func Err(err error) logger.Field {
	return func(e *logger.Event) {
		e.AppendString("error", err.Error())
	}
}

// Stack return a stack string context field.
func Stack(k string) logger.Field {
	return func(e *logger.Event) {
		cs := stack.Trace().TrimRuntime()

		c := stack.Caller(3)
		cs = cs.TrimBelow(c)

		e.AppendString(k, fmt.Sprintf("%+v", cs))
	}
}

// Caller returns a caller string context field.
func Caller(k string) logger.Field {
	return func(e *logger.Event) {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			e.AppendString(k, file+":"+strconv.Itoa(line))
		}
	}
}

// Time returns a time context field.
func Time(k string, t time.Time) logger.Field {
	return func(e *logger.Event) {
		e.AppendTime(k, t)
	}
}

// Duration returns a duration context field.
func Duration(k string, d time.Duration) logger.Field {
	return func(e *logger.Event) {
		e.AppendDuration(k, d)
	}
}

// Interface returns an interface context field.
func Interface(k string, v any) logger.Field {
	return func(e *logger.Event) {
		e.AppendInterface(k, v)
	}
}

// Span represents an open telemetry span.
type Span interface {
	IsRecording() bool
	SpanContext() trace.SpanContext
}

// TraceID returns an open telemetry trace ID context field.
func TraceID(k string, span Span) logger.Field {
	if !span.IsRecording() || !span.SpanContext().HasTraceID() {
		return func(*logger.Event) {}
	}
	return func(e *logger.Event) {
		e.AppendString(k, span.SpanContext().TraceID().String())
	}
}
