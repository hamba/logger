package logger_test

import (
	"log/slog"
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

func ExampleNewHandler() {
	h := logger.NewHandler(os.Stdout, logger.JSONFormat(), slog.LevelInfo)

	log := slog.New(h).With(slog.String("env", "prod")).WithGroup("db")

	log.Info("connected", slog.String("driver", "pgx"))
}

func ExampleNewHandler_logfmt() {
	h := logger.NewHandler(os.Stdout, logger.LogfmtFormat(), slog.LevelInfo)

	log := slog.New(h).With(slog.String("env", "prod")).WithGroup("db")

	log.Info("connected", slog.String("driver", "pgx"))
}
