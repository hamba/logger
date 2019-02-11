package logger_test

import "github.com/hamba/logger"

type CloseableHandler struct {
	CloseCalled bool
}

func (h *CloseableHandler) Log(msg string, lvl logger.Level, ctx []interface{}) {}

func (h *CloseableHandler) Close() error {
	h.CloseCalled = true
	return nil
}
