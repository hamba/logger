package logger_test

import (
	"context"
	"fmt"
	stdlog "log"
	"os"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
)

func ExampleNew() {
	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info).With(ctx.Str("env", "prod"))

	log.Info("redis connection", ctx.Str("redis", "some redis name"), ctx.Int("timeout", 10))
	log.Debug("dropped, more verbose than Info")

	// Output:
	// lvl=info msg="redis connection" env=prod redis="some redis name" timeout=10
}

func ExampleSyncWriter() {
	log := logger.New(logger.NewSyncWriter(os.Stdout), logger.LogfmtFormat(), logger.Info).With(ctx.Str("env", "prod"))

	log.Info("redis connection", ctx.Str("redis", "some redis name"), ctx.Int("timeout", 10))

	// Output:
	// lvl=info msg="redis connection" env=prod redis="some redis name" timeout=10
}

func ExampleWithContext() {
	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info).With(ctx.Str("svc", "api"))

	reqCtx := logger.WithContext(context.Background(), log, ctx.Str("req_id", "abc-123"), ctx.Str("method", "GET"))

	reqLog := log.FromContext(reqCtx)
	reqLog.Info("request handled", ctx.Int("status", 200))

	// Output:
	// lvl=info msg="request handled" svc=api req_id=abc-123 method=GET status=200
}

func ExampleLogger_Writer() {
	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info)

	std := stdlog.New(log.Writer(logger.Warn), "", 0)
	std.Println("disk almost full")

	// Output:
	// lvl=warn msg="disk almost full"
}

func ExampleLogger_WithTimestamp() {
	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info)

	// The cancel function must be called to stop the clock.
	cancel := log.WithTimestamp()
	defer cancel()

	log.Info("started")
}

func ExampleLevelFromString() {
	lvl, err := logger.LevelFromString("dbug")
	fmt.Println(lvl, err)

	_, err = logger.LevelFromString("verbose")
	fmt.Println(err)

	// Output:
	// dbug <nil>
	// unknown level verbose
}
