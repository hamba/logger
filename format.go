package logger

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/hamba/logger/v2/internal/bytes"
)

const (
	// TimestampKey is the key of the timestamp in the output of [JSONFormat]
	// and [LogfmtFormat].
	TimestampKey = "ts"
	// LevelKey is the key of the level in the output of [JSONFormat]
	// and [LogfmtFormat].
	LevelKey = "lvl"
	// MessageKey is the key of the message in the output of [JSONFormat]
	// and [LogfmtFormat].
	MessageKey = "msg"
)

// Formatter renders log messages into a buffer.
//
// For each message, a [Logger] calls AppendBeginMarker and WriteMessage, then
// AppendKey followed by one value method for each field, and finally
// AppendEndMarker and AppendLineBreak.
//
// Formatter uses a type from an internal package, so it can only be
// implemented within this module. Use [JSONFormat], [LogfmtFormat] or
// [ConsoleFormat] to obtain one.
type Formatter interface {
	// WriteMessage writes the timestamp, level and message that start a log
	// line. A zero ts means the Logger has no timestamp and none is written.
	WriteMessage(buf *bytes.Buffer, ts time.Time, lvl Level, msg string)
	// AppendBeginMarker writes any text that opens a log line.
	AppendBeginMarker(buf *bytes.Buffer)
	// AppendEndMarker writes any text that closes a log line.
	AppendEndMarker(buf *bytes.Buffer)
	// AppendLineBreak writes the text that terminates a log line.
	AppendLineBreak(buf *bytes.Buffer)
	// AppendArrayStart writes any text that opens an array value.
	AppendArrayStart(buf *bytes.Buffer)
	// AppendArraySep writes the text between two elements of an array value.
	AppendArraySep(buf *bytes.Buffer)
	// AppendArrayEnd writes any text that closes an array value.
	AppendArrayEnd(buf *bytes.Buffer)
	// AppendKey writes the separator before a field and the field key.
	// It is followed by a call that writes the value.
	AppendKey(buf *bytes.Buffer, key string)
	// AppendString writes s as a field value, escaping it as required.
	AppendString(buf *bytes.Buffer, s string)
	// AppendBool writes b as a field value.
	AppendBool(buf *bytes.Buffer, b bool)
	// AppendInt writes i as a field value.
	AppendInt(buf *bytes.Buffer, i int64)
	// AppendUint writes i as a field value.
	AppendUint(buf *bytes.Buffer, i uint64)
	// AppendFloat writes f as a field value.
	AppendFloat(buf *bytes.Buffer, f float64)
	// AppendTime writes t as a field value.
	AppendTime(buf *bytes.Buffer, t time.Time)
	// AppendDuration writes d as a field value.
	AppendDuration(buf *bytes.Buffer, d time.Duration)
	// AppendInterface writes the text form of v as a field value.
	AppendInterface(buf *bytes.Buffer, v any)
}

type json struct{}

// JSONFormat returns a [Formatter] that writes each message as a single-line
// JSON object. Times are rendered according to [TimeFormat], and values
// passed to AppendInterface are rendered as strings, or null for nil.
//
// Field keys are written as given and are not escaped.
func JSONFormat() Formatter {
	return &json{}
}

func (j *json) WriteMessage(buf *bytes.Buffer, ts time.Time, lvl Level, msg string) {
	if !ts.IsZero() {
		buf.WriteString("\"" + TimestampKey + "\":")
		j.AppendTime(buf, ts)
		buf.WriteString(",\"" + LevelKey + "\":\"")
	} else {
		buf.WriteString("\"" + LevelKey + "\":\"")
	}
	buf.WriteString(lvl.String())
	buf.WriteString("\",\"" + MessageKey + "\":")
	appendString(buf, msg, true)
}

func (j *json) AppendBeginMarker(buf *bytes.Buffer) {
	buf.WriteString("{")
}

func (j *json) AppendEndMarker(buf *bytes.Buffer) {
	buf.WriteString("}")
}

func (j *json) AppendLineBreak(buf *bytes.Buffer) {
	buf.WriteString("\n")
}

func (j *json) AppendArrayStart(buf *bytes.Buffer) {
	buf.WriteByte('[')
}

func (j *json) AppendArraySep(buf *bytes.Buffer) {
	buf.WriteByte(',')
}

func (j *json) AppendArrayEnd(buf *bytes.Buffer) {
	buf.WriteByte(']')
}

func (j *json) AppendKey(buf *bytes.Buffer, key string) {
	buf.WriteString(`,"`)
	buf.WriteString(key)
	buf.WriteString(`":`)
}

func (j *json) AppendString(buf *bytes.Buffer, s string) {
	appendString(buf, s, true)
}

func (j *json) AppendBool(buf *bytes.Buffer, b bool) {
	buf.AppendBool(b)
}

func (j *json) AppendInt(buf *bytes.Buffer, i int64) {
	buf.AppendInt(i)
}

func (j *json) AppendUint(buf *bytes.Buffer, i uint64) {
	buf.AppendUint(i)
}

func (j *json) AppendFloat(buf *bytes.Buffer, f float64) {
	buf.AppendFloat(f, 'g', -1, 64)
}

func (j *json) AppendTime(buf *bytes.Buffer, t time.Time) {
	switch TimeFormat {
	case TimeFormatUnix:
		buf.AppendInt(t.Unix())
	default:
		buf.WriteByte('"')
		buf.AppendTime(t, TimeFormat)
		buf.WriteByte('"')
	}
}

func (j *json) AppendDuration(buf *bytes.Buffer, d time.Duration) {
	buf.WriteByte('"')
	buf.AppendDuration(d)
	buf.WriteByte('"')
}

func (j *json) AppendInterface(buf *bytes.Buffer, v any) {
	if v == nil {
		buf.WriteString("null")
		return
	}

	j.AppendString(buf, fmt.Sprintf("%+v", v))
}

type logfmt struct{}

// LogfmtFormat returns a [Formatter] that writes each message as a line of
// space-separated key=value pairs. Times are rendered according to
// [TimeFormat]. Strings are quoted when they contain spaces, control
// characters, '=' or '"'.
//
// Field keys are written as given and are not escaped.
func LogfmtFormat() Formatter {
	return &logfmt{}
}

func (l *logfmt) needsQuote(s string) bool {
	for i := range len(s) {
		b := s[i]
		if b <= ' ' || b == '=' || b == '"' {
			return true
		}
	}
	return false
}

func (l *logfmt) WriteMessage(buf *bytes.Buffer, ts time.Time, lvl Level, msg string) {
	if !ts.IsZero() {
		buf.WriteString(TimestampKey + "=")
		l.AppendTime(buf, ts)
		buf.WriteString(" " + LevelKey + "=")
	} else {
		buf.WriteString(LevelKey + "=")
	}
	buf.WriteString(lvl.String())
	buf.WriteString(" " + MessageKey + "=")
	appendString(buf, msg, l.needsQuote(msg))
}

func (l *logfmt) AppendBeginMarker(*bytes.Buffer) {}

func (l *logfmt) AppendEndMarker(*bytes.Buffer) {}

func (l *logfmt) AppendLineBreak(buf *bytes.Buffer) {
	buf.WriteByte('\n')
}

func (l *logfmt) AppendArrayStart(_ *bytes.Buffer) {}

func (l *logfmt) AppendArraySep(buf *bytes.Buffer) {
	buf.WriteByte(',')
}

func (l *logfmt) AppendArrayEnd(_ *bytes.Buffer) {}

func (l *logfmt) AppendKey(buf *bytes.Buffer, key string) {
	buf.WriteByte(' ')
	buf.WriteString(key)
	buf.WriteByte('=')
}

func (l *logfmt) AppendString(buf *bytes.Buffer, s string) {
	appendString(buf, s, l.needsQuote(s))
}

func (l *logfmt) AppendBool(buf *bytes.Buffer, b bool) {
	buf.AppendBool(b)
}

func (l *logfmt) AppendInt(buf *bytes.Buffer, i int64) {
	buf.AppendInt(i)
}

func (l *logfmt) AppendUint(buf *bytes.Buffer, i uint64) {
	buf.AppendUint(i)
}

func (l *logfmt) AppendFloat(buf *bytes.Buffer, f float64) {
	buf.AppendFloat(f, 'f', 3, 64)
}

func (l *logfmt) AppendTime(buf *bytes.Buffer, t time.Time) {
	switch TimeFormat {
	case TimeFormatUnix:
		buf.AppendInt(t.Unix())
	default:
		buf.AppendTime(t, TimeFormat)
	}
}

func (l *logfmt) AppendDuration(buf *bytes.Buffer, d time.Duration) {
	buf.AppendDuration(d)
}

func (l *logfmt) AppendInterface(buf *bytes.Buffer, v any) {
	if v == nil {
		return
	}

	l.AppendString(buf, fmt.Sprintf("%+v", v))
}

const (
	// Foreground text colors.
	_ = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	_
	colorCyan
	colorWhite

	colorReset = 0
	colorBold  = 1
)

var noColor = newColor(colorReset)

type color []int

func newColor(attr ...int) color {
	return attr
}

func (c color) Write(buf *bytes.Buffer) {
	buf.WriteByte('\x1b')
	buf.WriteByte('[')
	for i := range c {
		if i > 0 {
			buf.WriteByte(';')
		}
		buf.AppendInt(int64(c[i]))
	}
	buf.WriteByte('m')
}

func withColor(c color, buf *bytes.Buffer, fn func()) {
	c.Write(buf)
	fn()
	noColor.Write(buf)
}

type console struct{}

// ConsoleFormat returns a [Formatter] for reading in a terminal. Each message
// is a line with the time, the upper-case level and the message, followed by
// key=value pairs, coloured with ANSI escape codes. Times are rendered
// according to [TimeFormatConsole]. Strings are never quoted.
//
// Field keys are written as given and are not escaped.
func ConsoleFormat() Formatter {
	return &console{}
}

func (c *console) lvlColor(lvl Level) color {
	switch lvl {
	case Crit:
		return newColor(colorRed, colorBold)
	case Error:
		return newColor(colorRed)
	case Warn:
		return newColor(colorYellow)
	case Info:
		return newColor(colorGreen)
	case Debug:
		return newColor(colorBlue)
	}
	return newColor(colorWhite)
}

func (c *console) WriteMessage(buf *bytes.Buffer, ts time.Time, lvl Level, msg string) {
	if !ts.IsZero() {
		withColor(newColor(colorBlue), buf, func() {
			c.AppendTime(buf, ts)
		})
		buf.WriteByte(' ')
	}
	withColor(c.lvlColor(lvl), buf, func() {
		buf.WriteString(strings.ToUpper(lvl.String()))
	})
	buf.WriteByte(' ')
	appendString(buf, msg, false)
}

func (c *console) AppendBeginMarker(*bytes.Buffer) {}

func (c *console) AppendEndMarker(*bytes.Buffer) {}

func (c *console) AppendLineBreak(buf *bytes.Buffer) {
	buf.WriteByte('\n')
}

func (c *console) AppendArrayStart(_ *bytes.Buffer) {}

func (c *console) AppendArraySep(buf *bytes.Buffer) {
	buf.WriteByte(',')
}

func (c *console) AppendArrayEnd(_ *bytes.Buffer) {}

func (c *console) AppendKey(buf *bytes.Buffer, key string) {
	buf.WriteByte(' ')

	col := newColor(colorCyan)
	if strings.HasPrefix(key, "err") {
		col = newColor(colorRed)
	}

	withColor(col, buf, func() {
		buf.WriteString(key)
		buf.WriteByte('=')
	})
}

func (c *console) AppendString(buf *bytes.Buffer, s string) {
	appendString(buf, s, false)
}

func (c *console) AppendBool(buf *bytes.Buffer, b bool) {
	buf.AppendBool(b)
}

func (c *console) AppendInt(buf *bytes.Buffer, i int64) {
	buf.AppendInt(i)
}

func (c *console) AppendUint(buf *bytes.Buffer, i uint64) {
	buf.AppendUint(i)
}

func (c *console) AppendFloat(buf *bytes.Buffer, f float64) {
	buf.AppendFloat(f, 'f', 3, 64)
}

func (c *console) AppendTime(buf *bytes.Buffer, t time.Time) {
	buf.AppendTime(t, TimeFormatConsole)
}

func (c *console) AppendDuration(buf *bytes.Buffer, d time.Duration) {
	buf.AppendDuration(d)
}

func (c *console) AppendInterface(buf *bytes.Buffer, v any) {
	if v == nil {
		return
	}

	c.AppendString(buf, fmt.Sprintf("%+v", v))
}

const hex = "0123456789abcdef"

//nolint:cyclop // Keeping unsplit for performance.
func appendString(buf *bytes.Buffer, s string, quote bool) {
	if quote {
		buf.WriteByte('"')
	}

	start := 0
	for i := 0; i < len(s); {
		b := s[i]
		if b-0x20 <= 0x5e && b != '"' && b != '\\' {
			i++
			continue
		}

		if start < i {
			buf.WriteString(s[start:i])
		}
		if tryAddASCII(buf, s[i]) {
			i++
			start = i
			continue
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			buf.WriteString(`\ufffd`)
			i++
			start = i
			continue
		}
		buf.WriteString(s[i : i+size])
		i += size
		start = i
	}

	if start < len(s) {
		if start == 0 {
			buf.WriteString(s)
		} else {
			buf.WriteString(s[start:])
		}
	}

	if quote {
		buf.WriteByte('"')
	}
}

func tryAddASCII(buf *bytes.Buffer, b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	switch b {
	case '\\', '"':
		buf.WriteByte('\\')
		buf.WriteByte(b)
	case '\n':
		buf.WriteString("\\n")
	case '\r':
		buf.WriteString("\\r")
	case '\t':
		buf.WriteString("\\t")
	default:
		buf.WriteString(`\u00`)
		buf.WriteByte(hex[b>>4])
		buf.WriteByte(hex[b&0xF])
	}
	return true
}
