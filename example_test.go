package logger_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/field"
)

func ExampleNew() {
	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info).With(field.Str("env", "prod"))

	log.Info("redis connection", field.Str("redis", "some redis name"), field.Int("timeout", 10))
}

func ExampleSyncWriter() {
	log := logger.New(logger.NewSyncWriter(os.Stdout), logger.LogfmtFormat(), logger.Info).With(field.Str("env", "prod"))

	log.Info("redis connection", field.Str("redis", "some redis name"), field.Int("timeout", 10))
}

func ExampleWithContext() {
	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info).With(field.Str("svc", "api"))

	reqCtx := logger.WithContext(context.Background(), log, field.Str("req_id", "abc-123"), field.Str("method", "GET"))

	reqLog := log.FromContext(reqCtx)
	reqLog.Info("request handled", field.Int("status", 200))
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
