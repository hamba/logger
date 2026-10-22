package logger

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"time"
)

// TimeFormat is the layout used by [JSONFormat] and [LogfmtFormat] to render
// timestamps and time fields. The default, [TimeFormatUnix], renders integer
// seconds since the Unix epoch. Any other value is a layout as described in
// [time.Time.Format].
//
// TimeFormat is read on every write, so it must be set before any [Logger]
// is used.
var TimeFormat = TimeFormatUnix

// TimeFormatConsole is the layout used by [ConsoleFormat] to render
// timestamps and time fields. It defaults to [time.Kitchen].
//
// TimeFormatConsole is read on every write, so it must be set before any
// [Logger] is used.
var TimeFormatConsole = time.Kitchen

const (
	// TimeFormatUnix selects integer seconds since the Unix epoch.
	TimeFormatUnix = ""
	// TimeFormatISO8601 is the layout "2006-01-02T15:04:05-0700".
	TimeFormatISO8601 = "2006-01-02T15:04:05-0700"
)

// Level is the severity of a log message and the verbosity threshold of a
// [Logger].
//
// Levels are ordered by verbosity: [Crit] is the least verbose and [Trace] the
// most. A Logger writes a message only if its level is no more verbose than
// the Logger's own level.
type Level int

// List of predefined log Levels.
const (
	Disabled Level = iota
	Crit
	Error
	Warn
	Info
	Debug
	Trace
)

// LevelFromString returns the [Level] named by lvl. The names are "trace",
// "debug", "info", "warn", "error" and "crit", plus the short forms "trce",
// "dbug" and "eror" returned by [Level.String]. Names are case-sensitive.
//
// It returns an error for any other name, including "disabled".
func LevelFromString(lvl string) (Level, error) {
	switch lvl {
	case "trace", "trce":
		return Trace, nil
	case "debug", "dbug":
		return Debug, nil
	case "info":
		return Info, nil
	case "warn":
		return Warn, nil
	case "error", "eror":
		return Error, nil
	case "crit":
		return Crit, nil
	default:
		return 0, errors.New("unknown level " + lvl)
	}
}

// String returns the four-letter name of the level, such as "info" or "dbug".
// Disabled and undefined levels return "unkn".
func (l Level) String() string {
	switch l {
	case Trace:
		return "trce"
	case Debug:
		return "dbug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "eror"
	case Crit:
		return "crit"
	default:
		return "unkn"
	}
}

// Field adds a key and value to an [Event]. The ctx package provides
// functions that create Fields, which are passed to [Logger.With] and the
// log methods of [Logger].
type Field func(*Event)

// Logger writes leveled, structured messages to an [io.Writer].
//
// Each message is written with a single call to the writer's Write method.
// Messages more verbose than the Logger's [Level] are dropped, and errors
// returned by the writer are ignored.
//
// Methods may be called concurrently by multiple goroutines, provided the
// writer allows concurrent calls to Write (see [SyncWriter]), with the
// exception of [Logger.WithTimestamp], which must not be called concurrently
// with any other method. The zero value is not usable; create a Logger with
// [New].
type Logger struct {
	w         io.Writer
	isDiscard bool
	fmtr      Formatter
	timeFn    func() time.Time
	ctx       []byte
	lvl       Level
}

// New returns a [Logger] that writes messages formatted by fmtr to w and
// drops messages more verbose than lvl. If w is [io.Discard], the Logger
// does no work.
func New(w io.Writer, fmtr Formatter, lvl Level) *Logger {
	isDiscard := w == io.Discard

	return &Logger{
		w:         w,
		isDiscard: isDiscard,
		fmtr:      fmtr,
		lvl:       lvl,
	}
}

// WithTimestamp makes the Logger add a timestamp to each message.
// Loggers derived from it afterwards inherit the timestamp; those derived
// earlier do not.
//
// Timestamps are refreshed periodically, so they may be slightly behind the
// time of the call. The returned cancel function stops the refreshing and
// must be called when the Logger is no longer needed; timestamps stop
// advancing once it has been called. If the Logger already has a timestamp,
// WithTimestamp changes nothing and returns a cancel function that does
// nothing.
func (l *Logger) WithTimestamp() (cancel func()) {
	if l.timeFn != nil {
		return func() {}
	}

	ctx, cancel := context.WithCancel(context.Background())

	var ts atomic.Value
	ts.Store(time.Now())

	go func() {
		tick := time.NewTicker(100 * time.Millisecond)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				ts.Store(time.Now())
			}
		}
	}()

	l.timeFn = func() time.Time {
		t := ts.Load().(time.Time)
		return t
	}

	return cancel
}

// With returns a new [Logger] that adds ctx to every message, after any fields
// already added to l. The receiver is not modified.
//
// The fields are evaluated once, when With is called, not on each message.
func (l *Logger) With(ctx ...Field) *Logger {
	e := newEvent(l.fmtr)
	defer putEvent(e)

	e.buf.Write(l.ctx)

	for _, field := range ctx {
		field(e)
	}

	b := make([]byte, e.buf.Len())
	copy(b, e.buf.Bytes())

	return &Logger{
		w:         l.w,
		isDiscard: l.isDiscard,
		fmtr:      l.fmtr,
		timeFn:    l.timeFn,
		lvl:       l.lvl,
		ctx:       b,
	}
}

// Trace logs msg and the fields in ctx at the [Trace] level.
// The fields follow those added with [Logger.With].
func (l *Logger) Trace(msg string, ctx ...Field) {
	l.write(msg, Trace, ctx)
}

// Debug logs msg and the fields in ctx at the [Debug] level.
func (l *Logger) Debug(msg string, ctx ...Field) {
	l.write(msg, Debug, ctx)
}

// Info logs msg and the fields in ctx at the [Info] level.
func (l *Logger) Info(msg string, ctx ...Field) {
	l.write(msg, Info, ctx)
}

// Warn logs msg and the fields in ctx at the [Warn] level.
func (l *Logger) Warn(msg string, ctx ...Field) {
	l.write(msg, Warn, ctx)
}

// Error logs msg and the fields in ctx at the [Error] level.
func (l *Logger) Error(msg string, ctx ...Field) {
	l.write(msg, Error, ctx)
}

// Crit logs msg and the fields in ctx at the [Crit] level.
func (l *Logger) Crit(msg string, ctx ...Field) {
	l.write(msg, Crit, ctx)
}

type writerFunc func([]byte) (int, error)

func (fn writerFunc) Write(p []byte) (n int, err error) {
	return fn(p)
}

// Writer returns an [io.Writer] that logs each call to Write as one message at
// level lvl, for use with packages such as [log]. A single trailing newline
// is removed from the message. Write always reports len(p) bytes written and
// a nil error, including when the message is dropped.
func (l *Logger) Writer(lvl Level) io.Writer {
	return writerFunc(func(p []byte) (n int, err error) {
		n = len(p)
		if l.isDiscard {
			return n, nil
		}

		if n > 0 && p[n-1] == '\n' {
			p = p[:n-1]
		}
		l.write(string(p), lvl, nil)

		return n, nil
	})
}

func (l *Logger) write(msg string, lvl Level, ctx []Field) {
	if l.isDiscard || lvl > l.lvl {
		return
	}

	e := newEvent(l.fmtr)

	var ts time.Time
	if l.timeFn != nil {
		ts = l.timeFn()
	}

	e.fmtr.AppendBeginMarker(e.buf)
	e.fmtr.WriteMessage(e.buf, ts, lvl, msg)
	e.buf.Write(l.ctx)

	for _, field := range ctx {
		field(e)
	}

	e.fmtr.AppendEndMarker(e.buf)
	e.fmtr.AppendLineBreak(e.buf)

	_, _ = l.w.Write(e.buf.Bytes())

	putEvent(e)
}
