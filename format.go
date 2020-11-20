package logger

import (
	"fmt"
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

var jsonPool = bytes.NewPool(512)

// JSONFormat formats a log line in json format.
func JSONFormat() Formatter {
	writeCtx := func(buf *bytes.Buffer, ctx []interface{}) {
		for i := 0; i < len(ctx); i += 2 {
			_ = buf.WriteByte(',')

			k, ok := ctx[i].(string)
			if !ok {
				buf.WriteString(`"` + errorKey + `"`)
				_ = buf.WriteByte(':')
				formatJSONValue(buf, ctx[i])
				continue
			}

			buf.WriteString(`"` + k + `"`)
			_ = buf.WriteByte(':')
			formatJSONValue(buf, ctx[i+1])
		}
	}

	return FormatterFunc(func(e *Event) []byte {
		buf := jsonPool.Get()

		// Append initial keys to the buffer
		_ = buf.WriteByte('{')
		buf.WriteString(`"` + LevelKey + `":"` + e.Lvl.String() + `",`)
		buf.WriteString(`"` + MessageKey + `":`)
		quoteString(buf, e.Msg)

		writeCtx(buf, e.BaseCtx)
		writeCtx(buf, e.Ctx)

		buf.WriteString("}\n")

		jsonPool.Put(buf)
		return buf.Bytes()
	})
}

// formatJSONValue formats a value, adding it to the buffer.
func formatJSONValue(buf *bytes.Buffer, value interface{}) {
	if value == nil {
		buf.WriteString("null")
		return
	}

	switch v := value.(type) {
	case time.Time:
		_ = buf.WriteByte('"')
		buf.AppendTime(v, timeFormat)
		_ = buf.WriteByte('"')
	case bool:
		buf.AppendBool(v)
	case float32:
		buf.AppendFloat(float64(v), 'g', -1, 64)
	case float64:
		buf.AppendFloat(v, 'g', -1, 64)
	case int:
		buf.AppendInt(int64(v))
	case int8:
		buf.AppendInt(int64(v))
	case int16:
		buf.AppendInt(int64(v))
	case int32:
		buf.AppendInt(int64(v))
	case int64:
		buf.AppendInt(v)
	case uint:
		buf.AppendUint(uint64(v))
	case uint8:
		buf.AppendUint(uint64(v))
	case uint16:
		buf.AppendUint(uint64(v))
	case uint32:
		buf.AppendUint(uint64(v))
	case uint64:
		buf.AppendUint(v)
	case string:
		quoteString(buf, v)
	default:
		quoteString(buf, fmt.Sprintf("%+v", value))
	}
}

var logfmtPool = bytes.NewPool(512)

// LogfmtFormat formats a log line in logfmt format.
func LogfmtFormat() Formatter {
	writeCtx := func(buf *bytes.Buffer, ctx []interface{}) {
		for i := 0; i < len(ctx); i += 2 {
			_ = buf.WriteByte(' ')

			k, ok := ctx[i].(string)
			if !ok {
				buf.WriteString(errorKey)
				_ = buf.WriteByte('=')
				formatLogfmtValue(buf, ctx[i])
				continue
			}

			buf.WriteString(k)
			_ = buf.WriteByte('=')
			formatLogfmtValue(buf, ctx[i+1])
		}
	}

	return FormatterFunc(func(e *Event) []byte {
		buf := logfmtPool.Get()

		// Append initial keys to the buffer
		buf.WriteString(LevelKey + "=" + e.Lvl.String() + " ")
		buf.WriteString(MessageKey + "=")
		logfmtQuoteString(buf, e.Msg)

		writeCtx(buf, e.BaseCtx)
		writeCtx(buf, e.Ctx)

		_ = buf.WriteByte('\n')

		logfmtPool.Put(buf)
		return buf.Bytes()
	})
}

// formatLogfmtValue formats a value, adding it to the buffer.
func formatLogfmtValue(buf *bytes.Buffer, value interface{}) {
	if value == nil {
		return
	}

	switch v := value.(type) {
	case time.Time:
		buf.AppendTime(v, timeFormat)
	case bool:
		buf.AppendBool(v)
	case float32:
		buf.AppendFloat(float64(v), 'f', 3, 64)
	case float64:
		buf.AppendFloat(v, 'f', 3, 64)
	case int:
		buf.AppendInt(int64(v))
	case int8:
		buf.AppendInt(int64(v))
	case int16:
		buf.AppendInt(int64(v))
	case int32:
		buf.AppendInt(int64(v))
	case int64:
		buf.AppendInt(v)
	case uint:
		buf.AppendUint(uint64(v))
	case uint8:
		buf.AppendUint(uint64(v))
	case uint16:
		buf.AppendUint(uint64(v))
	case uint32:
		buf.AppendUint(uint64(v))
	case uint64:
		buf.AppendUint(v)
	case string:
		logfmtQuoteString(buf, v)
	default:
		logfmtQuoteString(buf, fmt.Sprintf("%+v", value))
	}
}

func logfmtQuoteString(buf *bytes.Buffer, s string) {
	needsQuotes := false
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			needsQuotes = true
		}
	}

	if needsQuotes {
		_ = buf.WriteByte('"')
	}

	escapeString(buf, s)

	if needsQuotes {
		_ = buf.WriteByte('"')
	}
}

func quoteString(buf *bytes.Buffer, s string) {
	_ = buf.WriteByte('"')

	escapeString(buf, s)

	_ = buf.WriteByte('"')
}

func escapeString(buf *bytes.Buffer, s string) {
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
}
