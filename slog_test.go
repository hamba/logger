package logger_test

import (
	"bytes"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/hamba/logger/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	h := logger.NewHandler(io.Discard, logger.JSONFormat(), slog.LevelInfo)

	assert.IsType(t, &logger.Handler{}, h)
}

func TestHandler_Enabled(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		w      io.Writer
		minLvl slog.Level
		lvl    slog.Level
		want   bool
	}{
		{
			name:   "at min level",
			w:      io.Discard,
			minLvl: slog.LevelInfo,
			lvl:    slog.LevelInfo,
			want:   false, // io.Discard short-circuits
		},
		{
			name:   "above min level",
			w:      &bytes.Buffer{},
			minLvl: slog.LevelInfo,
			lvl:    slog.LevelWarn,
			want:   true,
		},
		{
			name:   "at min level",
			w:      &bytes.Buffer{},
			minLvl: slog.LevelInfo,
			lvl:    slog.LevelInfo,
			want:   true,
		},
		{
			name:   "below min level",
			w:      &bytes.Buffer{},
			minLvl: slog.LevelInfo,
			lvl:    slog.LevelDebug,
			want:   false,
		},
		{
			name:   "discard writer",
			w:      io.Discard,
			minLvl: slog.LevelDebug,
			lvl:    slog.LevelDebug,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			h := logger.NewHandler(test.w, logger.JSONFormat(), test.minLvl)

			got := h.Enabled(t.Context(), test.lvl)

			assert.Equal(t, test.want, got)
		})
	}
}

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		fmtr  logger.Formatter
		level slog.Level
		attrs []slog.Attr
		want  string
	}{
		{
			name:  "json info no attrs",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			want:  `{"lvl":"info","msg":"hello"}` + "\n",
		},
		{
			name:  "json debug",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelDebug,
			want:  `{"lvl":"dbug","msg":"hello"}` + "\n",
		},
		{
			name:  "json warn",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelWarn,
			want:  `{"lvl":"warn","msg":"hello"}` + "\n",
		},
		{
			name:  "json error",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelError,
			want:  `{"lvl":"eror","msg":"hello"}` + "\n",
		},
		{
			name:  "logfmt info no attrs",
			fmtr:  logger.LogfmtFormat(),
			level: slog.LevelInfo,
			want:  "lvl=info msg=hello\n",
		},
		{
			name:  "json with string attr",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.String("env", "prod")},
			want:  `{"lvl":"info","msg":"hello","env":"prod"}` + "\n",
		},
		{
			name:  "json with int attr",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.Int("port", 5432)},
			want:  `{"lvl":"info","msg":"hello","port":5432}` + "\n",
		},
		{
			name:  "json with bool attr",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.Bool("ok", true)},
			want:  `{"lvl":"info","msg":"hello","ok":true}` + "\n",
		},
		{
			name:  "json with float attr",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.Float64("ratio", 1.5)},
			want:  `{"lvl":"info","msg":"hello","ratio":1.5}` + "\n",
		},
		{
			name:  "json with duration attr",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.Duration("latency", 5*time.Millisecond)},
			want:  `{"lvl":"info","msg":"hello","latency":"5ms"}` + "\n",
		},
		{
			name:  "json with any attr",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.Any("val", struct{ X int }{X: 1})},
			want:  `{"lvl":"info","msg":"hello","val":"{X:1}"}` + "\n",
		},
		{
			name:  "json with group attr",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.Group("db", slog.String("host", "localhost"), slog.Int("port", 5432))},
			want:  `{"lvl":"info","msg":"hello","db":{"host":"localhost","port":5432}}` + "\n",
		},
		{
			name:  "logfmt with group attr",
			fmtr:  logger.LogfmtFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.Group("db", slog.String("host", "localhost"), slog.Int("port", 5432))},
			want:  "lvl=info msg=hello db.host=localhost db.port=5432\n",
		},
		{
			name:  "json with empty group attr discarded",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.Group("empty")},
			want:  `{"lvl":"info","msg":"hello"}` + "\n",
		},
		{
			name:  "json with anonymous group attr flattened",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{slog.Group("", slog.String("a", "1"), slog.String("b", "2"))},
			want:  `{"lvl":"info","msg":"hello","a":"1","b":"2"}` + "\n",
		},
		{
			name:  "json with zero attr discarded",
			fmtr:  logger.JSONFormat(),
			level: slog.LevelInfo,
			attrs: []slog.Attr{{}, slog.String("env", "prod")},
			want:  `{"lvl":"info","msg":"hello","env":"prod"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			h := logger.NewHandler(&buf, test.fmtr, slog.LevelDebug)

			r := slog.NewRecord(time.Time{}, test.level, "hello", 0)
			r.AddAttrs(test.attrs...)

			err := h.Handle(t.Context(), r)

			require.NoError(t, err)
			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		fmtr  logger.Formatter
		attrs []slog.Attr
		extra []slog.Attr
		want  string
	}{
		{
			name:  "json pre-serialises attrs",
			fmtr:  logger.JSONFormat(),
			attrs: []slog.Attr{slog.String("env", "prod")},
			want:  `{"lvl":"info","msg":"hello","env":"prod"}` + "\n",
		},
		{
			name:  "logfmt pre-serialises attrs",
			fmtr:  logger.LogfmtFormat(),
			attrs: []slog.Attr{slog.String("env", "prod")},
			want:  "lvl=info msg=hello env=prod\n",
		},
		{
			name:  "json pre-serialised plus inline attrs",
			fmtr:  logger.JSONFormat(),
			attrs: []slog.Attr{slog.String("env", "prod")},
			extra: []slog.Attr{slog.String("req", "GET /")},
			want:  `{"lvl":"info","msg":"hello","env":"prod","req":"GET /"}` + "\n",
		},
		{
			name:  "json chained WithAttrs",
			fmtr:  logger.JSONFormat(),
			attrs: []slog.Attr{slog.String("env", "prod")},
			extra: []slog.Attr{slog.String("svc", "api")},
			want:  `{"lvl":"info","msg":"hello","env":"prod","svc":"api"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			h := logger.NewHandler(&buf, test.fmtr, slog.LevelInfo)
			h2 := h.WithAttrs(test.attrs)

			// Chain a second WithAttrs when extra is provided.
			if len(test.extra) > 0 {
				h2 = h2.WithAttrs(test.extra)
			}

			r := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello", 0)
			err := h2.Handle(t.Context(), r)

			require.NoError(t, err)
			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestHandler_WithAttrs_empty(t *testing.T) {
	t.Parallel()

	h := logger.NewHandler(io.Discard, logger.JSONFormat(), slog.LevelInfo)
	got := h.WithAttrs(nil)

	assert.Same(t, h, got.(*logger.Handler))
}

func TestHandler_WithGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fmtr      logger.Formatter
		group     string
		withAttrs []slog.Attr
		inline    []slog.Attr
		want      string
	}{
		{
			name:   "json nested group with inline attrs",
			fmtr:   logger.JSONFormat(),
			group:  "db",
			inline: []slog.Attr{slog.String("host", "localhost")},
			want:   `{"lvl":"info","msg":"hello","db":{"host":"localhost"}}` + "\n",
		},
		{
			name:   "logfmt prefix group with inline attrs",
			fmtr:   logger.LogfmtFormat(),
			group:  "db",
			inline: []slog.Attr{slog.String("host", "localhost")},
			want:   "lvl=info msg=hello db.host=localhost\n",
		},
		{
			name:      "json nested group with pre-serialised and inline attrs",
			fmtr:      logger.JSONFormat(),
			group:     "db",
			withAttrs: []slog.Attr{slog.String("host", "localhost"), slog.Int("port", 5432)},
			inline:    []slog.Attr{slog.String("driver", "pgx")},
			want:      `{"lvl":"info","msg":"hello","db":{"host":"localhost","port":5432,"driver":"pgx"}}` + "\n",
		},
		{
			name:      "logfmt prefix group with pre-serialised and inline attrs",
			fmtr:      logger.LogfmtFormat(),
			group:     "db",
			withAttrs: []slog.Attr{slog.String("host", "localhost"), slog.Int("port", 5432)},
			inline:    []slog.Attr{slog.String("driver", "pgx")},
			want:      "lvl=info msg=hello db.host=localhost db.port=5432 db.driver=pgx\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			h := logger.NewHandler(&buf, test.fmtr, slog.LevelInfo)
			h2 := h.WithGroup(test.group)
			if len(test.withAttrs) > 0 {
				h2 = h2.WithAttrs(test.withAttrs)
			}

			r := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello", 0)
			r.AddAttrs(test.inline...)
			err := h2.Handle(t.Context(), r)

			require.NoError(t, err)
			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestHandler_WithGroup_WithAttrs_WithGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fmtr logger.Formatter
		want string
	}{
		{
			name: "json chained groups with attrs at each level",
			fmtr: logger.JSONFormat(),
			want: `{"lvl":"info","msg":"hello","a":{"x":"1","b":{"y":"2","z":"3"}}}` + "\n",
		},
		{
			name: "logfmt chained groups with attrs at each level",
			fmtr: logger.LogfmtFormat(),
			want: "lvl=info msg=hello a.x=1 a.b.y=2 a.b.z=3\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			h := logger.NewHandler(&buf, test.fmtr, slog.LevelInfo)

			// WithGroup("a") → WithAttrs([x]) → WithGroup("b") → WithAttrs([y])
			h2 := h.WithGroup("a").
				WithAttrs([]slog.Attr{slog.String("x", "1")}).
				WithGroup("b").
				WithAttrs([]slog.Attr{slog.String("y", "2")})

			r := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello", 0)
			r.AddAttrs(slog.String("z", "3"))
			err := h2.Handle(t.Context(), r)

			require.NoError(t, err)
			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestHandler_WithGroupHandlesEmptyName(t *testing.T) {
	t.Parallel()

	h := logger.NewHandler(io.Discard, logger.JSONFormat(), slog.LevelInfo)
	got := h.WithGroup("")

	assert.Same(t, h, got.(*logger.Handler))
}

func TestHandler_WithGroup_WithAttrsUsingJSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	h := logger.NewHandler(&buf, logger.JSONFormat(), slog.LevelInfo)

	log := slog.New(h).
		With(slog.String("env", "prod")).
		WithGroup("db").
		With(slog.String("host", "localhost"), slog.Int("port", 5432))

	log.Info("connected", slog.String("driver", "pgx"))

	out := buf.String()
	assert.Contains(t, out, `"lvl":"info"`)
	assert.Contains(t, out, `"msg":"connected"`)
	assert.Contains(t, out, `"env":"prod"`)
	assert.Contains(t, out, `"db":{"host":"localhost","port":5432,"driver":"pgx"}`)
}

func TestHandler_WithGroup_WithAttrsUsingLogfmt(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	h := logger.NewHandler(&buf, logger.LogfmtFormat(), slog.LevelInfo)

	log := slog.New(h).
		With(slog.String("env", "prod")).
		WithGroup("db").
		With(slog.String("host", "localhost"), slog.Int("port", 5432))

	log.Info("connected", slog.String("driver", "pgx"))

	out := buf.String()
	assert.Contains(t, out, "lvl=info")
	assert.Contains(t, out, "msg=connected")
	assert.Contains(t, out, "env=prod")
	assert.Contains(t, out, "db.host=localhost")
	assert.Contains(t, out, "db.port=5432")
	assert.Contains(t, out, "db.driver=pgx")
}

func TestHandler_NestedGroups(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fmtr logger.Formatter
		want string
	}{
		{
			name: "json",
			fmtr: logger.JSONFormat(),
			want: `{"lvl":"info","msg":"hello","a":{"b":{"c":"v"}}}` + "\n",
		},
		{
			name: "logfmt",
			fmtr: logger.LogfmtFormat(),
			want: "lvl=info msg=hello a.b.c=v\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			h := logger.NewHandler(&buf, test.fmtr, slog.LevelInfo)

			r := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello", 0)
			r.AddAttrs(slog.Group("a", slog.Group("b", slog.String("c", "v"))))

			err := h.Handle(t.Context(), r)

			require.NoError(t, err)
			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestHandler_MapSlogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		slogLvl slog.Level
		wantLvl string
	}{
		{slog.LevelDebug, "dbug"},
		{slog.LevelDebug - 4, "dbug"}, // custom level below debug
		{slog.LevelInfo, "info"},
		{slog.LevelWarn, "warn"},
		{slog.LevelError, "eror"},
		{slog.LevelError + 4, "eror"}, // custom level above error
	}

	for _, test := range tests {
		t.Run(test.slogLvl.String(), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			h := logger.NewHandler(&buf, logger.LogfmtFormat(), slog.LevelDebug-4)

			r := slog.NewRecord(time.Time{}, test.slogLvl, "msg", 0)
			err := h.Handle(t.Context(), r)

			require.NoError(t, err)
			assert.Contains(t, buf.String(), "lvl="+test.wantLvl)
		})
	}
}
