package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/hamba/logger/internal/bytes"
)

const (
	// LevelKey is the key used for message levels.
	LevelKey = "lvl"
	// MessageKey is the key used for message descriptions.
	MessageKey = "msg"

	timeFormat = "2006-01-02T15:04:05-0700" // ISO8601 format
)

// Formatter represents a log message formatter.
type Formatter interface {
	// Format formats a log message.
	Format(e *Event) []byte
}

// FormatterFunc is a function formatter.
type FormatterFunc func(e *Event) []byte

// Format formats a log message.
func (f FormatterFunc) Format(e *Event) []byte {
	return f(e)
}

type json struct {
	pool bytes.Pool
}

// JSONFormat formats a log line in json format.
func JSONFormat() Formatter {
	return &json{
		pool: bytes.NewPool(512),
	}
}

func (j *json) Format(e *Event) []byte {
	buf := j.pool.Get()

	_ = buf.WriteByte('{')
	buf.WriteString(`"` + LevelKey + `":"` + e.Lvl.String() + `",`)
	buf.WriteString(`"` + MessageKey + `":`)
	escapeString(buf, e.Msg, true)

	j.writeCtx(buf, e.BaseCtx)
	j.writeCtx(buf, e.Ctx)

	buf.WriteString("}\n")

	j.pool.Put(buf)
	return buf.Bytes()
}

func (j *json) writeCtx(buf *bytes.Buffer, ctx []interface{}) {
	for i := 0; i < len(ctx); i += 2 {
		_ = buf.WriteByte(',')

		k, ok := ctx[i].(string)
		if !ok {
			buf.WriteString(`"` + errorKey + `"`)
			_ = buf.WriteByte(':')
			j.writeValue(buf, ctx[i])
			continue
		}

		buf.WriteString(`"` + k + `"`)
		_ = buf.WriteByte(':')
		j.writeValue(buf, ctx[i+1])
	}
}

func (j *json) writeValue(buf *bytes.Buffer, v interface{}) {
	if v == nil {
		buf.WriteString("null")
		return
	}

	switch val := v.(type) {
	case time.Time:
		_ = buf.WriteByte('"')
		buf.AppendTime(val, timeFormat)
		_ = buf.WriteByte('"')
	case time.Duration:
		escapeString(buf, val.String(), true)
	case bool:
		buf.AppendBool(val)
	case float32:
		buf.AppendFloat(float64(val), 'g', -1, 64)
	case float64:
		buf.AppendFloat(val, 'g', -1, 64)
	case int:
		buf.AppendInt(int64(val))
	case int8:
		buf.AppendInt(int64(val))
	case int16:
		buf.AppendInt(int64(val))
	case int32:
		buf.AppendInt(int64(val))
	case int64:
		buf.AppendInt(val)
	case uint:
		buf.AppendUint(uint64(val))
	case uint8:
		buf.AppendUint(uint64(val))
	case uint16:
		buf.AppendUint(uint64(val))
	case uint32:
		buf.AppendUint(uint64(val))
	case uint64:
		buf.AppendUint(val)
	case string:
		escapeString(buf, val, true)
	default:
		escapeString(buf, fmt.Sprintf("%+v", v), true)
	}
}

type logfmt struct {
	pool bytes.Pool
}

// LogfmtFormat formats a log line in logfmt format.
func LogfmtFormat() Formatter {
	return &logfmt{
		pool: bytes.NewPool(512),
	}
}

func (l *logfmt) Format(e *Event) []byte {
	buf := l.pool.Get()

	buf.WriteString(LevelKey + "=" + e.Lvl.String() + " ")
	buf.WriteString(MessageKey + "=")
	escapeString(buf, e.Msg, needsQuote(e.Msg))

	l.writeCtx(buf, e.BaseCtx)
	l.writeCtx(buf, e.Ctx)

	_ = buf.WriteByte('\n')

	l.pool.Put(buf)
	return buf.Bytes()
}

func (l *logfmt) writeCtx(buf *bytes.Buffer, ctx []interface{}) {
	for i := 0; i < len(ctx); i += 2 {
		_ = buf.WriteByte(' ')

		k, ok := ctx[i].(string)
		if !ok {
			buf.WriteString(errorKey)
			_ = buf.WriteByte('=')
			l.writeValue(buf, ctx[i])
			continue
		}

		buf.WriteString(k)
		_ = buf.WriteByte('=')
		l.writeValue(buf, ctx[i+1])
	}
}

func (l *logfmt) writeValue(buf *bytes.Buffer, v interface{}) {
	if v == nil {
		return
	}

	switch val := v.(type) {
	case time.Time:
		buf.AppendTime(val, timeFormat)
	case time.Duration:
		escapeString(buf, val.String(), false)
	case bool:
		buf.AppendBool(val)
	case float32:
		buf.AppendFloat(float64(val), 'f', 3, 64)
	case float64:
		buf.AppendFloat(val, 'f', 3, 64)
	case int:
		buf.AppendInt(int64(val))
	case int8:
		buf.AppendInt(int64(val))
	case int16:
		buf.AppendInt(int64(val))
	case int32:
		buf.AppendInt(int64(val))
	case int64:
		buf.AppendInt(val)
	case uint:
		buf.AppendUint(uint64(val))
	case uint8:
		buf.AppendUint(uint64(val))
	case uint16:
		buf.AppendUint(uint64(val))
	case uint32:
		buf.AppendUint(uint64(val))
	case uint64:
		buf.AppendUint(val)
	case string:
		escapeString(buf, val, needsQuote(val))
	default:
		str := fmt.Sprintf("%+v", v)
		escapeString(buf, str, needsQuote(str))
	}
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
	return color(attr)
}

func (c color) Write(buf *bytes.Buffer) {
	if len(c) == 0 {
		return
	}

	_ = buf.WriteByte('\x1b')
	_ = buf.WriteByte('[')
	for i := 0; i < len(c); i++ {
		if i > 0 {
			_ = buf.WriteByte(';')
		}
		buf.AppendInt(int64(c[i]))
	}
	_ = buf.WriteByte('m')
}

func withColor(c color, buf *bytes.Buffer, fn func()) {
	if len(c) == 0 {
		fn()
		return
	}

	c.Write(buf)
	fn()
	noColor.Write(buf)
}

type console struct {
	pool bytes.Pool
}

// ConsoleFormat formats a log line in a console format.
func ConsoleFormat() Formatter {
	return &console{
		pool: bytes.NewPool(512),
	}
}

func (c *console) Format(e *Event) []byte {
	buf := c.pool.Get()

	withColor(c.lvlColor(e.Lvl), buf, func() {
		buf.WriteString(strings.ToUpper(e.Lvl.String()))
	})
	_ = buf.WriteByte(' ')
	escapeString(buf, e.Msg, false)

	c.writeCtx(buf, e.BaseCtx)
	c.writeCtx(buf, e.Ctx)

	_ = buf.WriteByte('\n')

	c.pool.Put(buf)
	return buf.Bytes()
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

func (c *console) writeCtx(buf *bytes.Buffer, ctx []interface{}) {
	for i := 0; i < len(ctx); i += 2 {
		_ = buf.WriteByte(' ')

		k, ok := ctx[i].(string)
		if !ok {
			withColor(newColor(colorRed), buf, func() {
				buf.WriteString(errorKey)
				_ = buf.WriteByte('=')
				c.writeValue(buf, ctx[i], noColor)
			})
			continue
		}

		var nameCol, valCol = newColor(colorCyan), noColor
		if strings.HasPrefix(k, "err") {
			nameCol, valCol = newColor(colorRed), newColor(colorRed)
		}

		withColor(nameCol, buf, func() {
			buf.WriteString(k)
			_ = buf.WriteByte('=')
		})
		c.writeValue(buf, ctx[i+1], valCol)
	}
}

func (c *console) writeValue(buf *bytes.Buffer, v interface{}, color color) {
	if v == nil {
		return
	}

	needsColor := len(color) > 0 && color[0] != colorReset
	if needsColor {
		color.Write(buf)
	}

	switch val := v.(type) {
	case time.Time:
		buf.AppendTime(val, timeFormat)
	case time.Duration:
		escapeString(buf, val.String(), false)
	case bool:
		buf.AppendBool(val)
	case float32:
		buf.AppendFloat(float64(val), 'f', 3, 64)
	case float64:
		buf.AppendFloat(val, 'f', 3, 64)
	case int:
		buf.AppendInt(int64(val))
	case int8:
		buf.AppendInt(int64(val))
	case int16:
		buf.AppendInt(int64(val))
	case int32:
		buf.AppendInt(int64(val))
	case int64:
		buf.AppendInt(val)
	case uint:
		buf.AppendUint(uint64(val))
	case uint8:
		buf.AppendUint(uint64(val))
	case uint16:
		buf.AppendUint(uint64(val))
	case uint32:
		buf.AppendUint(uint64(val))
	case uint64:
		buf.AppendUint(val)
	case string:
		escapeString(buf, val, false)
	default:
		str := fmt.Sprintf("%+v", v)
		escapeString(buf, str, false)
	}

	if needsColor {
		noColor.Write(buf)
	}
}

func needsQuote(s string) bool {
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			return true
		}
	}
	return false
}

func escapeString(buf *bytes.Buffer, s string, quote bool) {
	if quote {
		_ = buf.WriteByte('"')
	}

	for _, r := range s {
		switch r {
		case '\\', '"':
			_ = buf.WriteByte('\\')
			_ = buf.WriteByte(byte(r))
		case '\n':
			buf.WriteString("\\n")
		case '\r':
			buf.WriteString("\\r")
		case '\t':
			buf.WriteString("\\t")
		default:
			_ = buf.WriteByte(byte(r))
		}
	}

	if quote {
		_ = buf.WriteByte('"')
	}
}
