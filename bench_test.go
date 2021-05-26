package logger_test

import (
	"io"
	"os"
	"testing"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
)

func BenchmarkLogger_Logfmt(b *testing.B) {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkLogger_Json(b *testing.B) {
	log := logger.New(io.Discard, logger.JSONFormat(), logger.Debug)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkLogger_LogfmtWithTS(b *testing.B) {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug)

	cancel := log.WithTimestamp()
	defer cancel()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkLogger_JsonWithTS(b *testing.B) {
	log := logger.New(io.Discard, logger.JSONFormat(), logger.Debug)

	cancel := log.WithTimestamp()
	defer cancel()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkLogger_LogfmtCtx(b *testing.B) {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", 1))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}

func BenchmarkLogger_JsonCtx(b *testing.B) {
	log := logger.New(io.Discard, logger.JSONFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", 1))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}

func BenchmarkLevelLogger_Logfmt(b *testing.B) {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", os.Getpid()))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Debug("debug", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			log.Info("info", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			log.Warn("warn", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			log.Error("error", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}

func BenchmarkLevelLogger_Json(b *testing.B) {
	log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", os.Getpid()))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Debug("debug", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			log.Info("info", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			log.Warn("warn", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
			log.Error("error", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}
