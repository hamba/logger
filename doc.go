// Package logger implements leveled, structured logging.
//
// A [Logger] writes one message per line to an [io.Writer], rendered by a
// [Formatter]: [JSONFormat], [LogfmtFormat] or [ConsoleFormat]. Messages carry
// key-value fields, created with the functions in the ctx subpackage, either
// attached to the Logger with [Logger.With] or passed with each message.
// Messages more verbose than the Logger's [Level] are dropped.
//
// A Logger may be used by multiple goroutines if its writer allows it; wrap
// writers that do not with [NewSyncWriter]. Fields can also travel in a
// [context.Context]; see [WithContext] and [Logger.FromContext].
//
// Example usage:
//
//	import (
//		"os"
//
//		"github.com/hamba/logger/v2"
//		"github.com/hamba/logger/v2/ctx"
//	)
//
//	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info)
//
//	// Logger can have scoped context
//	log = log.With(ctx.Str("env", "prod"))
//
//	// All messages can have a context
//	log.Error("connection error", ctx.Str("redis", conn.Name()), ctx.Int("timeout", conn.Timeout()))
package logger
