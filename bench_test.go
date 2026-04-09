package logger_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/ctx"
)

func BenchmarkLogger_Logfmt(b *testing.B) {
	log := logger.New(discard{}, logger.LogfmtFormat(), logger.Debug)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkLogger_Json(b *testing.B) {
	log := logger.New(discard{}, logger.JSONFormat(), logger.Debug)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkLogger_LogfmtWriter(b *testing.B) {
	log := logger.New(discard{}, logger.LogfmtFormat(), logger.Debug)
	w := log.Writer(logger.Info)

	p := []byte("some message")

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = w.Write(p)
		}
	})
}

func BenchmarkLogger_LogfmtWithTS(b *testing.B) {
	log := logger.New(discard{}, logger.LogfmtFormat(), logger.Debug)

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
	log := logger.New(discard{}, logger.JSONFormat(), logger.Debug)

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
	log := logger.New(discard{}, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", 1))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}

func BenchmarkLogger_JsonCtx(b *testing.B) {
	log := logger.New(discard{}, logger.JSONFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", 1))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message", ctx.Int("key", 1), ctx.Float64("key2", 3.141592), ctx.Str("key3", "string"), ctx.Bool("key4", false))
		}
	})
}

func BenchmarkLogger_WithContext(b *testing.B) {
	log := logger.New(discard{}, logger.LogfmtFormat(), logger.Debug)
	goCtx := context.Background()

	_ = logger.WithContext(goCtx, log, ctx.Str("req_id", "abc123"), ctx.Int("attempt", 1))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = logger.WithContext(goCtx, log, ctx.Str("req_id", "abc123"), ctx.Int("attempt", 1))
		}
	})
}

func BenchmarkHandler_Logfmt(b *testing.B) {
	h := logger.NewHandler(discard{}, logger.LogfmtFormat(), slog.LevelDebug)
	log := slog.New(h)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkLogger_FromContext(b *testing.B) {
	log := logger.New(discard{}, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"))
	goCtx := logger.WithContext(context.Background(), log, ctx.Str("req_id", "abc123"), ctx.Int("attempt", 1))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.FromContext(goCtx).Info("some message")
		}
	})
}

func BenchmarkHandler_Json(b *testing.B) {
	h := logger.NewHandler(discard{}, logger.JSONFormat(), slog.LevelDebug)
	log := slog.New(h)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkHandler_LogfmtWithGroup(b *testing.B) {
	h := logger.NewHandler(discard{}, logger.LogfmtFormat(), slog.LevelDebug)
	log := slog.New(h.WithGroup("service"))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkHandler_JsonWithGroup(b *testing.B) {
	h := logger.NewHandler(discard{}, logger.JSONFormat(), slog.LevelDebug)
	log := slog.New(h.WithGroup("service"))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message")
		}
	})
}

func BenchmarkHandler_LogfmtWithGroupAndAttrs(b *testing.B) {
	h := logger.NewHandler(discard{}, logger.LogfmtFormat(), slog.LevelDebug)
	log := slog.New(h.WithGroup("db").WithAttrs([]slog.Attr{
		slog.String("host", "localhost"),
		slog.Int("port", 5432),
	}))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message", slog.String("driver", "pgx"))
		}
	})
}

func BenchmarkHandler_JsonWithGroupAndAttrs(b *testing.B) {
	h := logger.NewHandler(discard{}, logger.JSONFormat(), slog.LevelDebug)
	log := slog.New(h.WithGroup("db").WithAttrs([]slog.Attr{
		slog.String("host", "localhost"),
		slog.Int("port", 5432),
	}))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Error("some message", slog.String("driver", "pgx"))
		}
	})
}

func BenchmarkLevelLogger_Logfmt(b *testing.B) {
	log := logger.New(discard{}, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", os.Getpid()))

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
	log := logger.New(discard{}, logger.LogfmtFormat(), logger.Debug).With(ctx.Str("_n", "bench"), ctx.Int("_p", os.Getpid()))

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

type discard struct{}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}
