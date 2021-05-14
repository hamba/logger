package logger_test

import (
	"os"

	"github.com/hamba/logger"
	"github.com/hamba/logger/ctx"
)

func ExampleNew() {
	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info).With(ctx.Str("env", "prod"))

	log.Info("redis connection", ctx.Str("redis", "some redis name"), ctx.Int("timeout", 10))
}

func ExampleSyncWriter() {
	log := logger.New(logger.NewSyncWriter(os.Stdout), logger.LogfmtFormat(), logger.Info).With(ctx.Str("env", "prod"))

	log.Info("redis connection", ctx.Str("redis", "some redis name"), ctx.Int("timeout", 10))
}
