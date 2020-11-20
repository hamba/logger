package logger

import (
	"io"
	"sync"
	"time"

	"github.com/hamba/logger/internal/bytes"
)

// Handler represents a log handler.
type Handler interface {
	// Log write the log message.
	Log(e *Event)
}

// HandlerFunc is a function handler.
type HandlerFunc func(e *Event)

// Log write the log message.
func (h HandlerFunc) Log(e *Event) {
	h(e)
}

type bufStreamHandler struct {
	flushBytes    int
	flushInterval time.Duration
	w             io.Writer
	fmtr          Formatter

	mx   sync.Mutex
	pool bytes.Pool
	buf  *bytes.Buffer
	ch   chan *bytes.Buffer

	shutdown chan bool
}

// BufferedStreamHandler writes buffered log messages to an io.Writer with the given format.
func BufferedStreamHandler(w io.Writer, flushBytes int, flushInterval time.Duration, fmtr Formatter) Handler {
	pool := bytes.NewPool(flushBytes)

	h := &bufStreamHandler{
		flushBytes:    flushBytes,
		flushInterval: flushInterval,
		fmtr:          fmtr,
		w:             w,
		pool:          pool,
		buf:           pool.Get(),
		ch:            make(chan *bytes.Buffer, 32),
		shutdown:      make(chan bool, 1),
	}

	go h.run()

	return h
}

func (h *bufStreamHandler) run() {
	doneChan := make(chan bool)

	go func() {
		for buf := range h.ch {
			_, _ = h.w.Write(buf.Bytes())
			h.pool.Put(buf)
		}
		doneChan <- true
	}()

	ticker := time.NewTicker(h.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.withBufferLock(func() {
				h.swap()
			})

		case <-doneChan:
			h.shutdown <- true
			return
		}
	}
}

// Log write the log message.
func (h *bufStreamHandler) Log(e *Event) {
	h.withBufferLock(func() {
		// Dont write to a closed buffer
		if h.buf == nil {
			return
		}

		_, _ = h.buf.Write(h.fmtr.Format(e))

		if h.buf.Len() >= h.flushBytes {
			h.swap()
		}
	})
}

// Close closes the handler, waiting for all buffers to be flushed.
func (h *bufStreamHandler) Close() error {
	h.withBufferLock(func() {
		h.swap()
		h.buf = nil
	})

	close(h.ch)
	<-h.shutdown

	return nil
}

func (h *bufStreamHandler) withBufferLock(fn func()) {
	h.mx.Lock()
	fn()
	h.mx.Unlock()
}

func (h *bufStreamHandler) swap() {
	if h.buf == nil || h.buf.Len() == 0 {
		return
	}

	old := h.buf
	h.buf = h.pool.Get()
	h.ch <- old
}

// StreamHandler writes log messages to an io.Writer with the given format.
func StreamHandler(w io.Writer, fmtr Formatter) Handler {
	var mu sync.Mutex

	h := func(e *Event) {
		mu.Lock()
		_, _ = w.Write(fmtr.Format(e))
		mu.Unlock()
	}

	return HandlerFunc(h)
}

// FilterFunc represents a function that can filter messages.
type FilterFunc func(e *Event) bool

// FilterHandler returns a handler that only writes messages to the wrapped
// handler if the given function evaluates true.
func FilterHandler(fn FilterFunc, h Handler) Handler {
	c := &closeHandler{
		Handler: HandlerFunc(func(e *Event) {
			if fn(e) {
				h.Log(e)
			}
		}),
	}

	if ch, ok := h.(io.Closer); ok {
		c.Closer = ch
	}

	return c
}

// LevelFilterHandler returns a handler that filters by log level.
func LevelFilterHandler(maxLvl Level, h Handler) Handler {
	return FilterHandler(func(e *Event) bool {
		return e.Lvl <= maxLvl
	}, h)
}

// DiscardHandler does nothing, discarding all log messages.
func DiscardHandler() Handler {
	return HandlerFunc(func(e *Event) {})
}

// closeHandler wraps a handler allowing it to close if it has a Close method.
type closeHandler struct {
	io.Closer
	Handler
}

func (h *closeHandler) Close() error {
	if h.Closer != nil {
		return h.Closer.Close()
	}

	return nil
}
