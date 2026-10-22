package logger

import (
	"io"
	"sync"
)

// SyncWriter is an [io.Writer] that serialises calls to Write on the wrapped
// writer, so it is safe for concurrent use by multiple goroutines. The zero
// value is not usable; create a SyncWriter with [NewSyncWriter].
type SyncWriter struct {
	mu sync.Mutex
	w  io.Writer
}

// NewSyncWriter returns a [SyncWriter] that writes to w.
func NewSyncWriter(w io.Writer) *SyncWriter {
	return &SyncWriter{w: w}
}

// Write writes p to the wrapped writer, allowing one call at a time, and
// returns the byte count and error of the wrapped writer.
func (w *SyncWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	n, err = w.w.Write(p)
	w.mu.Unlock()

	return n, err
}
