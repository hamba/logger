package ctx

import (
	"time"

	"github.com/hamba/logger"
)

func Str(k, s string) logger.Field {
	return func(e *logger.Event) {
		e.AppendString(k, s)
	}
}

func Strs(k string, s []string) logger.Field {
	return func(e *logger.Event) {
		e.AppendStrings(k, s)
	}
}

func Bytes(k string, p []byte) logger.Field {
	return func(e *logger.Event) {
		e.AppendBytes(k, p)
	}
}

func Bool(k string, b bool) logger.Field {
	return func(e *logger.Event) {
		e.AppendBool(k, b)
	}
}

func Int(k string, i int) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

func Ints(k string, a []int) logger.Field {
	return func(e *logger.Event) {
		e.AppendInts(k, a)
	}
}

func Int8(k string, i int8) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

func Int16(k string, i int16) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

func Int32(k string, i int32) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, int64(i))
	}
}

func Int64(k string, i int64) logger.Field {
	return func(e *logger.Event) {
		e.AppendInt(k, i)
	}
}

func Uint(k string, i uint) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

func Uint8(k string, i uint8) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

func Uint16(k string, i uint16) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

func Uint32(k string, i uint32) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, uint64(i))
	}
}

func Uint64(k string, i uint64) logger.Field {
	return func(e *logger.Event) {
		e.AppendUint(k, i)
	}
}

func Float32(k string, f float32) logger.Field {
	return func(e *logger.Event) {
		e.AppendFloat(k, float64(f))
	}
}

func Float64(k string, f float64) logger.Field {
	return func(e *logger.Event) {
		e.AppendFloat(k, f)
	}
}

func Error(k string, err error) logger.Field {
	return func(e *logger.Event) {
		e.AppendString(k, err.Error())
	}
}

func Time(k string, t time.Time) logger.Field {
	return func(e *logger.Event) {
		e.AppendTime(k, t)
	}
}

func Duration(k string, d time.Duration) logger.Field {
	return func(e *logger.Event) {
		e.AppendDuration(k, d)
	}
}

func Interface(k string, v interface{}) logger.Field {
	return func(e *logger.Event) {
		e.AppendInterface(k, v)
	}
}
