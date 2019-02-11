package logger_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/hamba/logger"
)

func BenchmarkLogged_Logfmt(b *testing.B) {
	buf := &bytes.Buffer{}
	l := logger.New(logger.StreamHandler(buf, logger.LogfmtFormat()), "_n", "bench", "_p", 1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Error("some message", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
}

func BenchmarkLogged_Json(b *testing.B) {
	buf := &bytes.Buffer{}
	l := logger.New(logger.StreamHandler(buf, logger.JSONFormat()), "_n", "bench", "_p", 1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Error("some message", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
}

func BenchmarkLevelLogged_Logfmt(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := logger.New(logger.StreamHandler(buf, logger.LogfmtFormat()), "_n", "bench", "_p", os.Getpid())
	for i := 0; i < b.N; i++ {
		l.Debug("debug", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Info("info", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Warn("warn", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Error("error", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
}

func BenchmarkLevelLogged_Json(b *testing.B) {
	buf := &bytes.Buffer{}
	b.ResetTimer()
	l := logger.New(logger.StreamHandler(buf, logger.JSONFormat()), "_n", "bench", "_p", os.Getpid())
	for i := 0; i < b.N; i++ {
		l.Debug("debug", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Info("info", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Warn("warn", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
		l.Error("error", "key", 1, "key2", 3.141592, "key3", "string", "key4", false)
	}
	b.StopTimer()
}
