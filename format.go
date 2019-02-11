package logger

import (
	"fmt"
	"time"

	"github.com/hamba/logger/internal/bytes"
)

const (
	//LevelKey is the key used for message levels.
	LevelKey = "lvl"
	// MessageKey is the key used for message descriptions.
	MessageKey = "msg"

	timeFormat = "2006-01-02T15:04:05-0700" // ISO8601 format
)

// Formatter represents a log message formatter.
type Formatter interface {
	// Format formats a log message.
	Format(msg string, lvl Level, ctx []interface{}) []byte
}

// FormatterFunc is a function formatter.
type FormatterFunc func(msg string, lvl Level, ctx []interface{}) []byte

// Format formats a log message.
func (f FormatterFunc) Format(msg string, lvl Level, ctx []interface{}) []byte {
	return f(msg, lvl, ctx)
}

var jsonPool = bytes.NewPool(512)

// JSONFormat formats a log line in json format.
func JSONFormat() Formatter {
	return FormatterFunc(func(msg string, lvl Level, ctx []interface{}) []byte {
		buf := jsonPool.Get()

		// Append initial keys to the buffer
		buf.WriteByte('{')
		buf.WriteString(`"` + LevelKey + `":"` + lvl.String() + `",`)
		buf.WriteString(`"` + MessageKey + `":`)
		quoteString(buf, msg)

		// Append ctx to the buffer
		for i := 0; i < len(ctx); i += 2 {
			buf.WriteByte(',')

			k, ok := ctx[i].(string)
			if !ok {
				buf.WriteString(`"` + errorKey + `"`)
				buf.WriteByte(':')
				formatJSONValue(buf, ctx[i])
				continue
			}

			buf.WriteString(`"` + k + `"`)
			buf.WriteByte(':')
			formatJSONValue(buf, ctx[i+1])
		}

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
		buf.WriteByte('"')
		buf.AppendTime(v, timeFormat)
		buf.WriteByte('"')
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
	return FormatterFunc(func(msg string, lvl Level, ctx []interface{}) []byte {
		buf := logfmtPool.Get()

		// Append initial keys to the buffer
		buf.WriteString(LevelKey + "=" + lvl.String() + " ")
		buf.WriteString(MessageKey + "=")
		logfmtQuoteString(buf, msg)

		// Append ctx to the buffer
		for i := 0; i < len(ctx); i += 2 {
			buf.WriteByte(' ')

			k, ok := ctx[i].(string)
			if !ok {
				buf.WriteString(errorKey)
				buf.WriteByte('=')
				formatLogfmtValue(buf, ctx[i])
				continue
			}

			buf.WriteString(k)
			buf.WriteByte('=')
			formatLogfmtValue(buf, ctx[i+1])
		}

		buf.WriteByte('\n')

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
		buf.WriteByte('"')
	}

	escapeString(buf, s)

	if needsQuotes {
		buf.WriteByte('"')
	}
}

func quoteString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')

	escapeString(buf, s)

	buf.WriteByte('"')
}

func escapeString(buf *bytes.Buffer, s string) {
	for _, r := range s {
		switch r {
		case '\\', '"':
			buf.WriteByte('\\')
			buf.WriteByte(byte(r))
		case '\n':
			buf.WriteString("\\n")
		case '\r':
			buf.WriteString("\\r")
		case '\t':
			buf.WriteString("\\t")
		default:
			buf.WriteByte(byte(r))
		}
	}
}
