package logger_test

import (
	"bytes"
	"errors"
	"io"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug)

	assert.IsType(t, &logger.Logger{}, log)
}

func TestLogger(t *testing.T) {
	tests := []struct {
		name string
		fn   func(l *logger.Logger)
		want string
	}{
		{
			name: "Debug",
			fn:   func(l *logger.Logger) { l.Debug("debug", ctx.Str("level", "debug")) },
			want: "lvl=dbug msg=debug level=debug\n",
		},
		{
			name: "Info",
			fn:   func(l *logger.Logger) { l.Info("info", ctx.Str("level", "info")) },
			want: "lvl=info msg=info level=info\n",
		},
		{
			name: "Warn",
			fn:   func(l *logger.Logger) { l.Warn("warn", ctx.Str("level", "warn")) },
			want: "lvl=warn msg=warn level=warn\n",
		},
		{
			name: "Error",
			fn:   func(l *logger.Logger) { l.Error("error", ctx.Str("level", "error")) },
			want: "lvl=eror msg=error level=error\n",
		},
		{
			name: "Crit",
			fn:   func(l *logger.Logger) { l.Crit("critical", ctx.Str("level", "critical")) },
			want: "lvl=crit msg=critical level=critical\n",
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

	span := &fakeSpan{ID: byte(2), Recording: true}

	_, file, line, _ := runtime.Caller(0)
	caller := file + ":" + strconv.Itoa(line+3)

	log.Info("some message",
		ctx.Str("str", "string"),
		ctx.Strs("strs", []string{"string1", "string2"}),
		ctx.Bytes("bytes", []byte("bytes")),
		ctx.Bool("bool", true),
		ctx.Int("int", 1),
		ctx.Ints("ints", []int{1, 2, 3}),
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
		ctx.Error("err", errors.New("test error")),
		ctx.Err(errors.New("test error")),
		ctx.Time("time", time.Unix(1541573670, 0).UTC()),
		ctx.Duration("dur", time.Second),
		ctx.Interface("obj", obj),
		ctx.Caller("caller"),
		ctx.TraceID("tid", span),
	)

	want := `lvl=info msg="some message" _n=bench _p=1 str=string strs=string1,string2 bytes=98,121,116,101,115 bool=true int=1 ints=1,2,3 int8=2 int16=3 int32=4 int64=5 uint=1 uint8=2 uint16=3 uint32=4 uint64=5 float32=1.230 float64=4.560 err="test error" error="test error" time=1541573670 dur=1s obj={Name:test} caller=` + caller + " tid=01000000000000000000000000000000\n"
	assert.Equal(t, want, buf.String())
}

func TestLogger_Stack(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Info)

	log.Info("some message", ctx.Stack("stack"))

	want := `lvl=info msg="some message" stack=[github.com/hamba/logger/logger/logger_test.go:232]` + "\n"
	assert.Equal(t, want, buf.String())
}

func TestLogger_Timestamp(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Info)
	cancel := log.WithTimestamp()
	defer cancel()

	log.Info("some message")

	want := `^ts=\d+ lvl=info msg="some message"` + "\n$"
	assert.Regexp(t, want, buf.String())
}

func TestLogger_Writer(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Info)
	w := log.Writer(logger.Info)

	n, err := w.Write([]byte("some message\n"))
	require.NoError(t, err)

	want := `lvl=info msg="some message"` + "\n"
	assert.Equal(t, 13, n)
	assert.Equal(t, want, buf.String())
}

type fakeSpan struct {
	Recording bool
	ID        byte
}

func (s *fakeSpan) End(...trace.SpanEndOption)              {}
func (s *fakeSpan) AddEvent(string, ...trace.EventOption)   {}
func (s *fakeSpan) IsRecording() bool                       { return s.Recording }
func (s *fakeSpan) RecordError(error, ...trace.EventOption) {}
func (s *fakeSpan) SetStatus(codes.Code, string)            {}
func (s *fakeSpan) SetName(string)                          {}
func (s *fakeSpan) SetAttributes(...attribute.KeyValue)     {}
func (s *fakeSpan) TracerProvider() trace.TracerProvider    { return nil }

func (s *fakeSpan) SpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID([16]byte{1}),
		SpanID:  trace.SpanID([8]byte{s.ID}),
		Remote:  false,
	})
}
