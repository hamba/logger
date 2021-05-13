package logger_test

import (
	"io"
	"os"
	"testing"

	"github.com/hamba/logger"
	"github.com/hamba/logger/ctx"
)

func BenchmarkLogged_Logfmt(b *testing.B) {
	l := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Error("some message")
		}
	})
}

func BenchmarkLogged_Json(b *testing.B) {
	l := logger.New(io.Discard, logger.JSONFormat(), logger.Debug)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Error("some message")
		}
	})
}

func BenchmarkLogged_LogfmtCtx(b *testing.B) {
	l := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", 1))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Error("some message", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}

func BenchmarkLogged_JsonCtx(b *testing.B) {
	l := logger.New(io.Discard, logger.JSONFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", 1))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Error("some message", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}

func BenchmarkLevelLogged_Logfmt(b *testing.B) {
	l := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", os.Getpid()))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debug("debug", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			l.Info("info", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			l.Warn("warn", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			l.Error("error", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}

func BenchmarkLevelLogged_Json(b *testing.B) {
	l := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", os.Getpid()))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debug("debug", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			l.Info("info", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			l.Warn("warn", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			l.Error("error", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}
