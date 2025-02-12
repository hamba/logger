package logger

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"time"
)

// TimeFormat is the format that times will be added in.
//
// TimeFormat defaults to unix time.
var TimeFormat = TimeFormatUnix

// Time formats.
const (
	TimeFormatUnix    = ""
	TimeFormatISO8601 = "2006-01-02T15:04:05-0700"
)

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

// Level represents the predefined log level.
type Level int

// LevelFromString converts a string to Level.
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

// String returns the string representation of the level.
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

// Field is a context field.
type Field func(*Event)

// Logger is a logger.
type Logger struct {
	w         io.Writer
	isDiscard bool
	fmtr      Formatter
	timeFn    func() time.Time
	ctx       []byte
	lvl       Level
}

// New creates a new Logger.
func New(w io.Writer, fmtr Formatter, lvl Level) *Logger {
	isDiscard := w == io.Discard

	return &Logger{
		w:         w,
		isDiscard: isDiscard,
		fmtr:      fmtr,
		lvl:       lvl,
	}
}

// WithTimestamp adds a timestamp to each log lone. Sub-loggers
// will inherit the timestamp.
//
// WithTimestamp is not thread safe.
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

// With returns a new Logger with the given context.
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

// Trace logs a trace message, intended for fine grained debug messages.
func (l *Logger) Trace(msg string, ctx ...Field) {
	l.write(msg, Trace, ctx)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, ctx ...Field) {
	l.write(msg, Debug, ctx)
}

// Info logs an informational message.
func (l *Logger) Info(msg string, ctx ...Field) {
	l.write(msg, Info, ctx)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, ctx ...Field) {
	l.write(msg, Warn, ctx)
}

// Error logs an error message.
func (l *Logger) Error(msg string, ctx ...Field) {
	l.write(msg, Error, ctx)
}

// Crit logs a critical message.
func (l *Logger) Crit(msg string, ctx ...Field) {
	l.write(msg, Crit, ctx)
}

type writerFunc func([]byte) (int, error)

func (fn writerFunc) Write(p []byte) (n int, err error) {
	return fn(p)
}

// Writer returns an io.Writer that writes at the given level.
// This can be used as a writer with the standard log library.
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
