package logger

import (
	"io"
	"sync"
)

// SyncWriter implements a writer that is synchronised with a lock.
type SyncWriter struct {
	mu sync.Mutex
	w  io.Writer
}

// NewSyncWriter returns a synchronised writer.
func NewSyncWriter(w io.Writer) *SyncWriter {
	return &SyncWriter{w: w}
}

// Write writes to the writer.
func (w *SyncWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	n, err = w.w.Write(p)
	w.mu.Unlock()

	return n, err
}
