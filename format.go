package logger

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/hamba/logger/v2/internal/bytes"
)

const (
	// TimestampKey is the key used for timestamps.
	TimestampKey = "ts"
	// LevelKey is the key used for message levels.
	LevelKey = "lvl"
	// MessageKey is the key used for message descriptions.
	MessageKey = "msg"
)

// Formatter represents a log message formatter.
type Formatter interface {
	WriteMessage(buf *bytes.Buffer, ts time.Time, lvl Level, msg string)
	AppendBeginMarker(buf *bytes.Buffer)
	AppendEndMarker(buf *bytes.Buffer)
	AppendLineBreak(buf *bytes.Buffer)
	AppendArrayStart(buf *bytes.Buffer)
	AppendArraySep(buf *bytes.Buffer)
	AppendArrayEnd(buf *bytes.Buffer)
	AppendKey(buf *bytes.Buffer, prefix []byte, key string)
	AppendGroupStart(buf *bytes.Buffer, prefix []byte, name string) []byte
	AppendGroupEnd(buf *bytes.Buffer, prefix []byte) []byte
	AppendString(buf *bytes.Buffer, s string)
	AppendBool(buf *bytes.Buffer, b bool)
	AppendInt(buf *bytes.Buffer, i int64)
	AppendUint(buf *bytes.Buffer, i uint64)
	AppendFloat(buf *bytes.Buffer, f float64)
	AppendTime(buf *bytes.Buffer, t time.Time)
	AppendDuration(buf *bytes.Buffer, d time.Duration)
	AppendInterface(buf *bytes.Buffer, v any)
}

type json struct{}

// JSONFormat formats a log line in json format.
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

func (j *json) AppendKey(buf *bytes.Buffer, _ []byte, key string) {
	if buf.Peek() != '{' {
		buf.WriteString(`,"`)
	} else {
		buf.WriteByte('"')
	}
	buf.WriteString(key)
	buf.WriteString(`":`)
}

func (j *json) AppendGroupStart(buf *bytes.Buffer, prefix []byte, name string) []byte {
	if buf.Peek() != '{' {
		buf.WriteString(`,"`)
	} else {
		buf.WriteByte('"')
	}
	buf.WriteString(name)
	buf.WriteString(`":{`)
	return prefix
}

func (j *json) AppendGroupEnd(buf *bytes.Buffer, prefix []byte) []byte {
	buf.WriteByte('}')
	return prefix
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

// LogfmtFormat formats a log line in logfmt format.
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

func (l *logfmt) AppendKey(buf *bytes.Buffer, prefix []byte, key string) {
	buf.WriteByte(' ')
	if len(prefix) > 0 {
		buf.Write(prefix)
	}
	buf.WriteString(key)
	buf.WriteByte('=')
}

func (l *logfmt) AppendGroupStart(_ *bytes.Buffer, prefix []byte, name string) []byte {
	prefix = append(prefix, name...)
	return append(prefix, '.')
}

func (l *logfmt) AppendGroupEnd(_ *bytes.Buffer, prefix []byte) []byte {
	return trimLastGroup(prefix)
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

// ConsoleFormat formats a log line in a console format.
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

func (c *console) AppendKey(buf *bytes.Buffer, prefix []byte, key string) {
	buf.WriteByte(' ')

	col := newColor(colorCyan)
	if strings.HasPrefix(key, "err") {
		col = newColor(colorRed)
	}

	withColor(col, buf, func() {
		if len(prefix) > 0 {
			buf.Write(prefix)
		}
		buf.WriteString(key)
		buf.WriteByte('=')
	})
}

func (c *console) AppendGroupStart(_ *bytes.Buffer, prefix []byte, name string) []byte {
	prefix = append(prefix, name...)
	return append(prefix, '.')
}

func (c *console) AppendGroupEnd(_ *bytes.Buffer, prefix []byte) []byte {
	return trimLastGroup(prefix)
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
	buf.AppendTime(t, time.Kitchen)
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

// trimLastGroup removes the last "name." segment from prefix by scanning
// backwards for the dot preceding the final segment.
func trimLastGroup(prefix []byte) []byte {
	for i := len(prefix) - 2; i >= 0; i-- {
		if prefix[i] == '.' {
			return prefix[:i+1]
		}
	}
	return prefix[:0]
}
