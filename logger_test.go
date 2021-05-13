package logger_test

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hamba/logger"
	"github.com/hamba/logger/ctx"
	"github.com/stretchr/testify/assert"
)

func TestLevelFromString(t *testing.T) {
	tests := []struct {
		lvl       string
		want      logger.Level
		wantError bool
	}{
		{
			lvl:       "dbug",
			want:      logger.Debug,
			wantError: false,
		},
		{
			lvl:       "debug",
			want:      logger.Debug,
			wantError: false,
		},
		{
			lvl:       "info",
			want:      logger.Info,
			wantError: false,
		},
		{
			lvl:       "warn",
			want:      logger.Warn,
			wantError: false,
		},
		{
			lvl:       "eror",
			want:      logger.Error,
			wantError: false,
		},
		{
			lvl:       "error",
			want:      logger.Error,
			wantError: false,
		},
		{
			lvl:       "crit",
			want:      logger.Crit,
			wantError: false,
		},
		{
			lvl:       "unkn",
			want:      logger.Level(123),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.lvl, func(t *testing.T) {
			lvl, err := logger.LevelFromString(tt.lvl)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, lvl)
		})
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		lvl  logger.Level
		want string
	}{
		{
			lvl:  logger.Debug,
			want: "dbug",
		},
		{
			lvl:  logger.Info,
			want: "info",
		},
		{
			lvl:  logger.Warn,
			want: "warn",
		},
		{
			lvl:  logger.Error,
			want: "eror",
		},
		{
			lvl:  logger.Crit,
			want: "crit",
		},
		{
			lvl:  logger.Level(123),
			want: "unkn",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.lvl.String())
	}
}

func TestNew(t *testing.T) {
	l := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug)

	assert.Implements(t, (*logger.Logger)(nil), l)
}

func TestLogger(t *testing.T) {
	tests := []struct {
		name string
		fn   func(l logger.Logger)
		want string
	}{
		{
			name: "Debug",
			fn:   func(l logger.Logger) { l.Debug("debug", ctx.Str("level", "debug")) },
			want: "lvl=dbug msg=debug level=debug",
		},
		{
			name: "Info",
			fn:   func(l logger.Logger) { l.Info("info", ctx.Str("level", "info")) },
			want: "lvl=info msg=info level=info",
		},
		{
			name: "Warn",
			fn:   func(l logger.Logger) { l.Warn("warn", ctx.Str("level", "warn")) },
			want: "lvl=warn msg=warn level=warn",
		},
		{
			name: "Error",
			fn:   func(l logger.Logger) { l.Error("error", ctx.Str("level", "error")) },
			want: "lvl=eror msg=error level=error",
		},
		{
			name: "Crit",
			fn:   func(l logger.Logger) { l.Crit("critical", ctx.Str("level", "critical")) },
			want: "lvl=crit msg=critical level=critical",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			log := logger.New(&buf, logger.LogfmtFormat(), logger.Debug)

			test.fn(log)

			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestLogger_DiscardsLogs(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Error)

	log.Debug("some message")

	assert.Equal(t, "", buf.String())
}

func TestLogger_Context(t *testing.T) {
	obj := struct {
		Name string
	}{Name: "test"}

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Info).With(ctx.Str("_n", "bench"), ctx.Int("_p", 1))

	log.Info("some message",
		ctx.Str("str", "string"),
		ctx.Bool("bool", true),
		ctx.Int("int", 1),
		ctx.Int8("int8", 2),
		ctx.Int16("int16", 3),
		ctx.Int32("int32", 4),
		ctx.Int64("int64", 5),
		ctx.Uint("uint", 1),
		ctx.Uint8("uint8", 2),
		ctx.Uint16("uint16", 3),
		ctx.Uint32("uint32", 4),
		ctx.Uint64("uint64", 5),
		ctx.Float32("float32", 1.23),
		ctx.Float64("float64", 4.56),
		ctx.Error("str", errors.New("test error")),
		ctx.Time("str", time.Unix(1541573670, 0).UTC()),
		ctx.Duration("str", time.Second),
		ctx.Interface("str", obj),
	)

	want := `lvl=info msg="some message" _n=bench _p=1 str=string bool=true int=1 int8=2 int16=3 int32=4 int64=5 uint=1 uint8=2 uint16=3 uint32=4 uint64=5 float32=1.230 float64=4.560 str="test error" str=2018-11-07T06:54:30+0000 str=1s str={Name:test}`
	assert.Equal(t, want, buf.String())
}

