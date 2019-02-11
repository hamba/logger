package logger

import (
	"errors"
	"io"
)

const errorKey = "LOGGER_ERROR"

// List of predefined log Levels
const (
	Crit Level = iota
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

// Logger represents a log writer.
type Logger interface {
	io.Closer

	// Debug logs a debug message.
	Debug(msg string, ctx ...interface{})
	// Info logs an informational message.
	Info(msg string, ctx ...interface{})
	// Warn logs a warning message.
	Warn(msg string, ctx ...interface{})
	// Error logs an error message.
	Error(msg string, ctx ...interface{})
	// Crit logs a critical message.
	Crit(msg string, ctx ...interface{})
}

type logger struct {
	h   Handler
	ctx []interface{}
}

// New creates a new Logger.
func New(h Handler, ctx ...interface{}) Logger {
	return &logger{
		h:   h,
		ctx: ctx,
	}
}

// Debug logs a debug message.
func (l *logger) Debug(msg string, ctx ...interface{}) {
	l.write(msg, Debug, ctx)
}

// Info logs an informational message.
func (l *logger) Info(msg string, ctx ...interface{}) {
	l.write(msg, Info, ctx)
}

// Warn logs a warning message.
func (l *logger) Warn(msg string, ctx ...interface{}) {
	l.write(msg, Warn, ctx)
}

// Error logs an error message.
func (l *logger) Error(msg string, ctx ...interface{}) {
	l.write(msg, Error, ctx)
}

// Crit logs a critical message.
func (l *logger) Crit(msg string, ctx ...interface{}) {
	l.write(msg, Crit, ctx)
}

func (l *logger) write(msg string, lvl Level, ctx []interface{}) {
	ctx = normalize(ctx)

	l.h.Log(msg, lvl, merge(l.ctx, ctx))
}

// Close closes the logger.
func (l *logger) Close() error {
	if h, ok := l.h.(io.Closer); ok {
		return h.Close()
	}

	return nil
}

func normalize(ctx []interface{}) []interface{} {
	// ctx needs to be even as they are key/value pairs
	if len(ctx)%2 != 0 {
		ctx = append(ctx, nil, errorKey, "Normalised odd number of arguments by adding nil")
	}

	return ctx
}

func merge(prefix, suffix []interface{}) []interface{} {
	newCtx := make([]interface{}, len(prefix)+len(suffix))
	n := copy(newCtx, prefix)
	copy(newCtx[n:], suffix)

	return newCtx
}
