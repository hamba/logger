package logger

import (
	"errors"
	"io"
)

// List of predefined log Levels.
const (
	Disabled Level = iota
	Crit
	Error
	Warn
	Info
	Debug
)

// Level represents the predefined log level.
type Level int

// LevelFromString converts a string to Level.
func LevelFromString(lvl string) (Level, error) {
	switch lvl {
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

type Logger struct {
	w    io.Writer
	fmtr Formatter
	lvl  Level
	ctx  []byte
}

// New creates a new Logger.
func New(w io.Writer, fmtr Formatter, lvl Level) *Logger {
	return &Logger{
		w:    w,
		fmtr: fmtr,
		lvl:  lvl,
	}
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
		w:    l.w,
		fmtr: l.fmtr,
		lvl:  l.lvl,
		ctx:  b,
	}
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

func (l *Logger) write(msg string, lvl Level, ctx []Field) {
	if lvl > l.lvl {
		return
	}

	e := newEvent(l.fmtr)

	e.fmtr.WriteMessage(e.buf, 0, lvl, msg)
	e.buf.Write(l.ctx)

	for _, field := range ctx {
		field(e)
	}

	e.fmtr.AppendEndMarker(e.buf)

	_, _ = l.w.Write(e.buf.Bytes())

	putEvent(e)
}
