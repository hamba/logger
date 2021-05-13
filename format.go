package logger

import (
	"fmt"
	"strings"
	"time"

	bytesx "github.com/hamba/logger/internal/bytes"
)

const (
	// LevelKey is the key used for message levels.
	LevelKey = "lvl"
	// MessageKey is the key used for message descriptions.
	MessageKey = "msg"

	// ISO8601 format
	timeFormat = "2006-01-02T15:04:05-0700"
)

// Formatter represents a log message formatter.
type Formatter interface {
	WriteMessage(buf *bytesx.Buffer, time int64, lvl Level, msg string)
	AppendEndMarker(buf *bytesx.Buffer)
	AppendKey(buf *bytesx.Buffer, key string)
	AppendString(buf *bytesx.Buffer, s string)
	AppendBool(buf *bytesx.Buffer, b bool)
	AppendInt(buf *bytesx.Buffer, i int64)
	AppendUint(buf *bytesx.Buffer, i uint64)
	AppendFloat(buf *bytesx.Buffer, f float64)
	AppendTime(buf *bytesx.Buffer, t time.Time)
	AppendDuration(buf *bytesx.Buffer, d time.Duration)
	AppendInterface(buf *bytesx.Buffer, v interface{})
}

type json struct{}

// JSONFormat formats a log line in json format.
func JSONFormat() Formatter {
	return &json{}
}

func (j *json) WriteMessage(buf *bytesx.Buffer, time int64, lvl Level, msg string) {
	buf.WriteString(`{"`)
	buf.WriteString(LevelKey)
	buf.WriteString(`":"`)
	buf.WriteString(lvl.String())
	buf.WriteString(`","`)
	buf.WriteString(MessageKey)
	buf.WriteString(`":`)
	escapeString(buf, msg, true)
}

func (j *json) AppendEndMarker(buf *bytesx.Buffer) {
	buf.WriteString("}\n")
}

func (j *json) AppendKey(buf *bytesx.Buffer, key string) {
	buf.WriteString(`,"`)
	buf.WriteString(key)
	buf.WriteString(`":`)
}

func (j *json) AppendString(buf *bytesx.Buffer, s string) {
	escapeString(buf, s, true)
}

func (j *json) AppendBool(buf *bytesx.Buffer, b bool) {
	buf.AppendBool(b)
}

func (j *json) AppendInt(buf *bytesx.Buffer, i int64) {
	buf.AppendInt(i)
}

func (j *json) AppendUint(buf *bytesx.Buffer, i uint64) {
	buf.AppendUint(i)
}

func (j *json) AppendFloat(buf *bytesx.Buffer, f float64) {
	buf.AppendFloat(f, 'g', -1, 64)
}

func (j *json) AppendTime(buf *bytesx.Buffer, t time.Time) {
	s := t.Format(timeFormat)
	escapeString(buf, s, true)
}

func (j *json) AppendDuration(buf *bytesx.Buffer, d time.Duration) {
	s := d.String()
	escapeString(buf, s, true)
}

func (j *json) AppendInterface(buf *bytesx.Buffer, v interface{}) {
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
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			return true
		}
	}
	return false
}

func (l *logfmt) WriteMessage(buf *bytesx.Buffer, time int64, lvl Level, msg string) {
	buf.WriteString(LevelKey)
	buf.WriteByte('=')
	buf.WriteString(lvl.String())
	buf.WriteByte(' ')
	buf.WriteString(MessageKey)
	buf.WriteByte('=')
	escapeString(buf, msg, l.needsQuote(msg))
}

func (l *logfmt) AppendEndMarker(buf *bytesx.Buffer) {
	buf.WriteByte('\n')
}

func (l *logfmt) AppendKey(buf *bytesx.Buffer, key string) {
	buf.WriteByte(' ')
	buf.WriteString(key)
	buf.WriteByte('=')
}

func (l *logfmt) AppendString(buf *bytesx.Buffer, s string) {
	escapeString(buf, s, l.needsQuote(s))
}

func (l *logfmt) AppendBool(buf *bytesx.Buffer, b bool) {
	buf.AppendBool(b)
}

func (l *logfmt) AppendInt(buf *bytesx.Buffer, i int64) {
	buf.AppendInt(i)
}

func (l *logfmt) AppendUint(buf *bytesx.Buffer, i uint64) {
	buf.AppendUint(i)
}

func (l *logfmt) AppendFloat(buf *bytesx.Buffer, f float64) {
	buf.AppendFloat(f, 'f', 3, 64)
}

func (l *logfmt) AppendTime(buf *bytesx.Buffer, t time.Time) {
	s := t.Format(timeFormat)
	escapeString(buf, s, l.needsQuote(s))
}

func (l *logfmt) AppendDuration(buf *bytesx.Buffer, d time.Duration) {
	s := d.String()
	escapeString(buf, s, l.needsQuote(s))
}

func (l *logfmt) AppendInterface(buf *bytesx.Buffer, v interface{}) {
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

func (c color) Write(buf *bytesx.Buffer) {
	buf.WriteByte('\x1b')
	buf.WriteByte('[')
	for i := 0; i < len(c); i++ {
		if i > 0 {
			buf.WriteByte(';')
		}
		buf.AppendInt(int64(c[i]))
	}
	buf.WriteByte('m')
}

func withColor(c color, buf *bytesx.Buffer, fn func()) {
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

func (c *console) WriteMessage(buf *bytesx.Buffer, time int64, lvl Level, msg string) {
	withColor(c.lvlColor(lvl), buf, func() {
		buf.WriteString(strings.ToUpper(lvl.String()))
	})
	buf.WriteByte(' ')
	escapeString(buf, msg, false)
}

func (c *console) AppendEndMarker(buf *bytesx.Buffer) {
	buf.WriteByte('\n')
}

func (c *console) AppendKey(buf *bytesx.Buffer, key string) {
	buf.WriteByte(' ')

	var col = newColor(colorCyan)
	if strings.HasPrefix(key, "err") {
		col = newColor(colorRed)
	}

	withColor(col, buf, func() {
		buf.WriteString(key)
		buf.WriteByte('=')
	})
}

func (c *console) AppendString(buf *bytesx.Buffer, s string) {
	escapeString(buf, s, false)
}

func (c *console) AppendBool(buf *bytesx.Buffer, b bool) {
	buf.AppendBool(b)
}

func (c *console) AppendInt(buf *bytesx.Buffer, i int64) {
	buf.AppendInt(i)
}

func (c *console) AppendUint(buf *bytesx.Buffer, i uint64) {
	buf.AppendUint(i)
}

func (c *console) AppendFloat(buf *bytesx.Buffer, f float64) {
	buf.AppendFloat(f, 'f', 3, 64)
}

func (c *console) AppendTime(buf *bytesx.Buffer, t time.Time) {
	s := t.Format(timeFormat)
	escapeString(buf, s, false)
}

func (c *console) AppendDuration(buf *bytesx.Buffer, d time.Duration) {
	s := d.String()
	escapeString(buf, s, false)
}

func (c *console) AppendInterface(buf *bytesx.Buffer, v interface{}) {
	if v == nil {
		return
	}

	c.AppendString(buf, fmt.Sprintf("%+v", v))
}

var noEscapeTable = [256]bool{}

func init() {
	for i := 0; i <= 0x7e; i++ {
		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
	}
}

// TODO: clean this up
func escapeString(buf *bytesx.Buffer, s string, quote bool) {
	if quote {
		buf.WriteByte('"')
	}

	var needEscape bool
	for i := 0; i < len(s); i++ {
		if noEscapeTable[s[i]] {
			continue
		}
		needEscape = true
	}

	if !needEscape {
		buf.WriteString(s)
	} else {

		for _, r := range s {
			switch r {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteRune(r)
			case '\n':
				buf.WriteString("\\n")
			case '\r':
				buf.WriteString("\\r")
			case '\t':
				buf.WriteString("\\t")
			default:
				buf.WriteRune(r)
			}
		}
	}

	if quote {
		buf.WriteByte('"')
	}
}
