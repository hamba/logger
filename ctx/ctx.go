// Package ctx provides functions that create the [logger.Field] values passed
// to the log methods and [logger.Logger.With] of a [logger.Logger].
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

// Str returns a [logger.Field] with key k and the string value s.
func Str(k, s string) logger.Field {
	return func(e *logger.Event) {
		e.AppendString(k, s)
	}
}

// Strs returns a [logger.Field] with key k and the strings s as an array.
func Strs(k string, s []string) logger.Field {
	return func(e *logger.Event) {
		e.AppendStrings(k, s)
	}
}

// Bytes returns a [logger.Field] with key k and the bytes p as an array of
// integers, not as a string.
func Bytes(k string, p []byte) logger.Field {
	return func(e *logger.Event) {
		e.AppendBytes(k, p)
	}
}

// Bool returns a [logger.Field] with key k and the bool value b.
func Bool(k string, b bool) logger.Field {
	return func(e *logger.Event) {
		e.AppendBool(k, b)
	}
}

// Int returns a [logger.Field] with key k and the int value i.
func Int(k string, i int) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

// Ints returns a [logger.Field] with key k and the ints a as an array.
func Ints(k string, a []int) logger.Field {
	return func(e *logger.Event) {
		e.AppendInts(k, a)
	}
}

// Int8 returns a [logger.Field] with key k and the int8 value i.
func Int8(k string, i int8) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

// Int16 returns a [logger.Field] with key k and the int16 value i.
func Int16(k string, i int16) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

// Int32 returns a [logger.Field] with key k and the int32 value i.
func Int32(k string, i int32) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

// Int64 returns a [logger.Field] with key k and the int64 value i.
func Int64(k string, i int64) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, i)
	}
}

// Uint returns a [logger.Field] with key k and the uint value i.
func Uint(k string, i uint) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

// Uint8 returns a [logger.Field] with key k and the uint8 value i.
func Uint8(k string, i uint8) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

// Uint16 returns a [logger.Field] with key k and the uint16 value i.
func Uint16(k string, i uint16) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

// Uint32 returns a [logger.Field] with key k and the uint32 value i.
func Uint32(k string, i uint32) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

// Uint64 returns a [logger.Field] with key k and the uint64 value i.
func Uint64(k string, i uint64) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, i)
	}
}

// Float32 returns a [logger.Field] with key k and the float32 value f.
// The value is converted to a float64 before it is rendered, so the output
// may show digits that are not significant in f.
func Float32(k string, f float32) logger.Field {
	return func(e *logger.Event) {
		e.AppendFloat(k, float64(f))
	}
}

// Float64 returns a [logger.Field] with key k and the float64 value f.
func Float64(k string, f float64) logger.Field {
	return func(e *logger.Event) {
		e.AppendFloat(k, f)
	}
}

// Error returns a [logger.Field] with key k and the text of err.
// The field panics when it is applied if err is nil.
func Error(k string, err error) logger.Field {
	return func(e *logger.Event) {
		e.AppendString(k, err.Error())
	}
}

// Err returns a [logger.Field] with key "error" and the text of err.
// The field panics when it is applied if err is nil.
func Err(err error) logger.Field {
	return func(e *logger.Event) {
		e.AppendString("error", err.Error())
	}
}

// Stack returns a [logger.Field] with key k and the call stack of the code
// that logged the message.
//
// Stack must be passed directly to a log method of a [logger.Logger]. Used
// in other ways, such as with [logger.Logger.With] or inside another
// [logger.Field], it does not report the caller's stack.
func Stack(k string) logger.Field {
	return func(e *logger.Event) {
		cs := stack.Trace().TrimRuntime()

		c := stack.Caller(3)
		cs = cs.TrimBelow(c)

		e.AppendString(k, fmt.Sprintf("%+v", cs))
	}
}

// Caller returns a [logger.Field] with key k and the "file:line" of the code
// that logged the message. The field is omitted if the location cannot be
// determined.
//
// Caller must be passed directly to a log method of a [logger.Logger]. Used
// in other ways, such as with [logger.Logger.With] or inside another
// [logger.Field], it does not report the caller's location.
func Caller(k string) logger.Field {
	return func(e *logger.Event) {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			e.AppendString(k, file+":"+strconv.Itoa(line))
		}
	}
}

// Time returns a [logger.Field] with key k and the time value t.
func Time(k string, t time.Time) logger.Field {
	return func(e *logger.Event) {
		e.AppendTime(k, t)
	}
}

// Duration returns a [logger.Field] with key k and the duration value d.
func Duration(k string, d time.Duration) logger.Field {
	return func(e *logger.Event) {
		e.AppendDuration(k, d)
	}
}

// Interface returns a [logger.Field] with key k and the value v, rendered as
// text. A nil v renders as null in JSON and as nothing in other formats.
func Interface(k string, v any) logger.Field {
	return func(e *logger.Event) {
		e.AppendInterface(k, v)
	}
}

// Span is the subset of the OpenTelemetry [trace.Span] interface used by
// [TraceID].
type Span interface {
	IsRecording() bool
	SpanContext() trace.SpanContext
}

// TraceID returns a [logger.Field] with key k and the OpenTelemetry trace ID
// of span. The field is omitted if span is not recording or has no trace ID.
func TraceID(k string, span Span) logger.Field {
	if !span.IsRecording() || !span.SpanContext().HasTraceID() {
		return func(*logger.Event) {}
	}
	return func(e *logger.Event) {
		e.AppendString(k, span.SpanContext().TraceID().String())
	}
}
