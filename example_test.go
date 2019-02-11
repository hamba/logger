package logger_test

import (
	"os"
	"time"

	"github.com/hamba/logger"
)

func ExampleNew() {
	h := logger.LevelFilterHandler(
		logger.Info,
		logger.StreamHandler(os.Stdout, logger.LogfmtFormat()),
	)

	l := logger.New(h, "env", "prod") // The logger can have an initial context

	l.Info("redis connection", "redis", "some redis name", "timeout", 10)
}

func ExampleBufferedStreamHandler() {
	h := logger.BufferedStreamHandler(os.Stdout, 2000, 1*time.Second, logger.LogfmtFormat())

	l := logger.New(h, "env", "prod")

	l.Info("redis connection", "redis", "some redis name", "timeout", 10)
}

func ExampleStreamHandler() {
	h := logger.StreamHandler(os.Stdout, logger.LogfmtFormat())

	l := logger.New(h, "env", "prod")

	l.Info("redis connection", "redis", "some redis name", "timeout", 10)
}

func ExampleLevelFilterHandler() {
	h := logger.LevelFilterHandler(
		logger.Info,
		logger.StreamHandler(os.Stdout, logger.LogfmtFormat()),
	)

	l := logger.New(h, "env", "prod")

	l.Info("redis connection", "redis", "some redis name", "timeout", 10)
}

func ExampleFilterHandler() {
	h := logger.FilterHandler(
		func(msg string, lvl logger.Level, ctx []interface{}) bool {
			return msg == "some condition"
		},
		logger.StreamHandler(os.Stdout, logger.LogfmtFormat()),
	)

	l := logger.New(h, "env", "prod")

	l.Info("redis connection", "redis", "some redis name", "timeout", 10)
}
