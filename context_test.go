package logger_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
	"github.com/stretchr/testify/assert"
)

func TestLogger_WithContextAttachesFields_Logfmt(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Info).With(ctx.Str("svc", "api"))

	goCtx := logger.WithContext(context.Background(), log, ctx.Str("req_id", "abc123"))
	log.FromContext(goCtx).Info("handled")

	assert.Equal(t, "lvl=info msg=handled svc=api req_id=abc123\n", buf.String())
}

func TestLogger_WithContextAttachesFields_JSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := logger.New(&buf, logger.JSONFormat(), logger.Info).With(ctx.Str("svc", "api"))

	goCtx := logger.WithContext(context.Background(), log, ctx.Str("req_id", "abc123"))
	log.FromContext(goCtx).Info("handled")

	assert.Equal(t, `{"lvl":"info","msg":"handled","svc":"api","req_id":"abc123"}`+"\n", buf.String())
}

func TestLogger_WithContextCanLayersFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Info)

	goCtx := logger.WithContext(context.Background(), log, ctx.Str("req_id", "abc123"))
	goCtx = logger.WithContext(goCtx, log, ctx.Str("user", "u456"))
	log.FromContext(goCtx).Info("handled")

	assert.Equal(t, "lvl=info msg=handled req_id=abc123 user=u456\n", buf.String())
}

func TestLogger_WithContextWithNoFieldsReturnsCtxUnchanged(t *testing.T) {
	t.Parallel()

	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Info)
	goCtx := context.Background()

	got := logger.WithContext(goCtx, log)

	assert.Equal(t, goCtx, got)
}

func TestLogger_WithContextWithDiscardReturnsCtxUnchanged(t *testing.T) {
	t.Parallel()

	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Info)
	goCtx := context.Background()

	got := logger.WithContext(goCtx, log, ctx.Str("req_id", "abc123"))

	assert.Equal(t, goCtx, got)
}

func TestLogger_FromContextWithNoContextFieldsReturnsSameLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LogfmtFormat(), logger.Info)

	got := log.FromContext(context.Background())

	assert.Same(t, log, got)
}

func TestLogger_FromContextWithDiscardReturnsSameLogger(t *testing.T) {
	t.Parallel()

	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Info)
	goCtx := logger.WithContext(context.Background(), log, ctx.Str("req_id", "abc123"))

	// Re-attach to a non-discard logger to prove the discard check fires first.
	discardLog := logger.New(io.Discard, logger.LogfmtFormat(), logger.Info)
	got := discardLog.FromContext(goCtx)

	assert.Same(t, discardLog, got)
}
