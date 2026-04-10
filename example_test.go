package logger_test

import (
	"context"
	"os"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
)

func ExampleNew() {
	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info).With(ctx.Str("env", "prod"))

	log.Info("redis connection", ctx.Str("redis", "some redis name"), ctx.Int("timeout", 10))
}

func ExampleSyncWriter() {
	log := logger.New(logger.NewSyncWriter(os.Stdout), logger.LogfmtFormat(), logger.Info).With(ctx.Str("env", "prod"))

	log.Info("redis connection", ctx.Str("redis", "some redis name"), ctx.Int("timeout", 10))
}

func ExampleWithContext() {
	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info).With(ctx.Str("svc", "api"))

	reqCtx := logger.WithContext(context.Background(), log, ctx.Str("req_id", "abc-123"), ctx.Str("method", "GET"))

	reqLog := log.FromContext(reqCtx)
	reqLog.Info("request handled", ctx.Int("status", 200))
}
