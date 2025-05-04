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
	"go.opentelemetry.io/otel/trace"
)

func TestLevelFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		lvl     string
		want    logger.Level
		wantErr require.ErrorAssertionFunc
	}{
		{
			lvl:     "trce",
			want:    logger.Trace,
			wantErr: require.NoError,
		},
		{
			lvl:     "trace",
			want:    logger.Trace,
			wantErr: require.NoError,
		},
		{
			lvl:     "dbug",
			want:    logger.Debug,
			wantErr: require.NoError,
		},
		{
			lvl:     "debug",
			want:    logger.Debug,
			wantErr: require.NoError,
		},
		{
			lvl:     "info",
			want:    logger.Info,
			wantErr: require.NoError,
		},
		{
			lvl:     "warn",
			want:    logger.Warn,
			wantErr: require.NoError,
		},
		{
			lvl:     "eror",
			want:    logger.Error,
			wantErr: require.NoError,
		},
		{
			lvl:     "error",
			want:    logger.Error,
			wantErr: require.NoError,
		},
		{
			lvl:     "crit",
			want:    logger.Crit,
			wantErr: require.NoError,
		},
		{
			lvl:     "unkn",
			wantErr: require.Error,
		},
	}

	for _, test := range tests {
		t.Run(test.lvl, func(t *testing.T) {
			t.Parallel()

			lvl, err := logger.LevelFromString(test.lvl)

			test.wantErr(t, err)
			assert.Equal(t, test.want, lvl)
		})
	}
}

func TestLevel_String(t *testing.T) {
	t.Parallel()

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

	for _, test := range tests {
		t.Run(test.lvl.String(), func(t *testing.T) {
			t.Parallel()

			got := test.lvl.String()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug)

	assert.IsType(t, &logger.Logger{}, log)
}

func TestLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func(l *logger.Logger)
		want string
	}{
		{
			name: "Trace",
			fn:   func(l *logger.Logger) { l.Trace("debug", ctx.Str("level", "trace")) },
			want: "lvl=trce msg=debug level=trace\n",
		},
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
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			log := logger.New(&buf, logger.LogfmtFormat(), logger.Trace)

			test.fn(log)

			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestLogger_DiscardsLogs(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Error)

	log.Debug("some message")

	assert.Empty(t, buf.String())
}

func TestLogger_Context(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Info)

	log.Info("some message", ctx.Stack("stack"))

	want := `lvl=info msg="some message" stack=[github.com/hamba/logger/logger/logger_test.go:259]` + "\n"
	assert.Equal(t, want, buf.String())
}

func TestLogger_Timestamp(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Info)
	cancel := log.WithTimestamp()
	defer cancel()

	log.Info("some message")

	want := `^ts=\d+ lvl=info msg="some message"` + "\n$"
	assert.Regexp(t, want, buf.String())
}

func TestLogger_Writer(t *testing.T) {
	t.Parallel()

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

func (s *fakeSpan) IsRecording() bool {
	return s.Recording
}

func (s *fakeSpan) SpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{s.ID},
		Remote:  false,
	})
}
