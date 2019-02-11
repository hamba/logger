package logger_test

import (
	"testing"

	"github.com/hamba/logger"
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
			assert.Equal(t, tt.want, lvl, )
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
	h := logger.HandlerFunc(func(msg string, lvl logger.Level, ctx []interface{}) {})

	l := logger.New(h)

	assert.Implements(t, (*logger.Logger)(nil), l)
}

func TestLogger(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(l logger.Logger)
		wantMsg string
		wantLvl logger.Level
		wantCtx []interface{}
	}{
		{
			name:    "Debug",
			fn:      func(l logger.Logger) { l.Debug("debug", "level", "debug") },
			wantMsg: "debug",
			wantLvl: logger.Debug,
			wantCtx: []interface{}{"level", "debug"},
		},
		{
			name:    "Info",
			fn:      func(l logger.Logger) { l.Info("info", "level", "info") },
			wantMsg: "info",
			wantLvl: logger.Info,
			wantCtx: []interface{}{"level", "info"},
		},
		{
			name:    "Warn",
			fn:      func(l logger.Logger) { l.Warn("warn", "level", "warn") },
			wantMsg: "warn",
			wantLvl: logger.Warn,
			wantCtx: []interface{}{"level", "warn"},
		},
		{
			name:    "Error",
			fn:      func(l logger.Logger) { l.Error("error", "level", "error") },
			wantMsg: "error",
			wantLvl: logger.Error,
			wantCtx: []interface{}{"level", "error"},
		},
		{
			name:    "Crit",
			fn:      func(l logger.Logger) { l.Crit("critical", "level", "critical") },
			wantMsg: "critical",
			wantLvl: logger.Crit,
			wantCtx: []interface{}{"level", "critical"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var outMsg string
			var outLvl logger.Level
			var outCtx []interface{}

			h := logger.HandlerFunc(func(msg string, lvl logger.Level, ctx []interface{}) {
				outMsg = msg
				outLvl = lvl
				outCtx = ctx
			})
			l := logger.New(h)

			tt.fn(l)

			assert.Equal(t, tt.wantMsg, outMsg)
			assert.Equal(t, tt.wantLvl, outLvl)
			assert.Equal(t, tt.wantCtx, outCtx)
		})
	}
}

func TestLogger_MergesCtx(t *testing.T) {
	var out []interface{}
	h := logger.HandlerFunc(func(msg string, lvl logger.Level, ctx []interface{}) {
		out = ctx
	})
	l := logger.New(h, "a", "b")

	l.Debug("test", "c", "d")

	assert.Equal(t, []interface{}{"a", "b", "c", "d"}, out)
}

func TestLogger_NormalizesCtx(t *testing.T) {
	var out []interface{}
	h := logger.HandlerFunc(func(msg string, lvl logger.Level, ctx []interface{}) {
		out = ctx
	})
	l := logger.New(h)

	l.Debug("test", "a")

	assert.Len(t, out, 4)
	assert.Equal(t, nil, out[1])
}

func TestLogger_TriesToCallUnderlyingClose(t *testing.T) {
	h := logger.HandlerFunc(func(msg string, lvl logger.Level, ctx []interface{}) {})
	l := logger.New(h)

	l.Close()
}

func TestLogger_CallsUnderlyingClose(t *testing.T) {
	h := &CloseableHandler{}
	l := logger.New(h)

	l.Close()

	assert.True(t, h.CloseCalled)
}
