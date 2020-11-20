package logger_test

import "github.com/hamba/logger"

type CloseableHandler struct {
	CloseCalled bool
}

func (h *CloseableHandler) Log(e *logger.Event) {}

func (h *CloseableHandler) Close() error {
	h.CloseCalled = true
	return nil
}
