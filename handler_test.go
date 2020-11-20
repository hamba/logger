package logger_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/hamba/logger"
	"github.com/stretchr/testify/assert"
)

func TestBufferedStreamHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	h := logger.BufferedStreamHandler(buf, 2000, time.Second, logger.LogfmtFormat())

	h.Log(&logger.Event{Msg: "some message", Lvl: logger.Error})
	h.(io.Closer).Close()

	assert.Equal(t, "lvl=eror msg=\"some message\"\n", buf.String())
}

func TestBufferedStreamHandler_SendsMessagesAfterFlushInterval(t *testing.T) {
	buf := &bytes.Buffer{}
	h := logger.BufferedStreamHandler(buf, 2000, time.Millisecond, logger.LogfmtFormat())
	defer h.(io.Closer).Close()

	h.Log(&logger.Event{Msg: "some message", Lvl: logger.Error})

	time.Sleep(5 * time.Millisecond)

	assert.Equal(t, "lvl=eror msg=\"some message\"\n", buf.String())
}

func TestBufferedStreamHandler_SendsMessagesAfterFlushBytes(t *testing.T) {
	buf := &bytes.Buffer{}
	h := logger.BufferedStreamHandler(buf, 40, time.Second, logger.LogfmtFormat())
	defer h.(io.Closer).Close()

	h.Log(&logger.Event{Msg: "some message", Lvl: logger.Error})
	h.Log(&logger.Event{Msg: "some message", Lvl: logger.Error})
	h.Log(&logger.Event{Msg: "some message", Lvl: logger.Error})

	time.Sleep(time.Millisecond)

	assert.Equal(t, "lvl=eror msg=\"some message\"\nlvl=eror msg=\"some message\"\n", buf.String())
}

func TestBufferedStreamHandler_DoesntWriteAfterClose(t *testing.T) {
	buf := &bytes.Buffer{}
	h := logger.BufferedStreamHandler(buf, 40, time.Second, logger.LogfmtFormat())
	h.(io.Closer).Close()

	h.Log(&logger.Event{Msg: "some message", Lvl: logger.Error})

	assert.Equal(t, "", buf.String())
}

func TestStreamHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	h := logger.StreamHandler(buf, logger.LogfmtFormat())

	h.Log(&logger.Event{Msg: "some message", Lvl: logger.Error})

	assert.Equal(t, "lvl=eror msg=\"some message\"\n", buf.String())
}

func TestLevelFilterHandler(t *testing.T) {
	count := 0
	testHandler := logger.HandlerFunc(func(e *logger.Event) {
		count++
	})
	h := logger.LevelFilterHandler(logger.Info, testHandler)

	h.Log(&logger.Event{Msg: "test", Lvl: logger.Debug})
	h.Log(&logger.Event{Msg: "test", Lvl: logger.Info})

	assert.Equal(t, 1, count)
}

func TestLevelFilterHandler_TriesToCallUnderlyingClose(t *testing.T) {
	testHandler := logger.HandlerFunc(func(e *logger.Event) {})
	h := logger.LevelFilterHandler(logger.Info, testHandler)
	ch := h.(io.Closer)

	ch.Close()
}

func TestLevelFilterHandler_CallsUnderlyingClose(t *testing.T) {
	testHandler := &CloseableHandler{}
	h := logger.LevelFilterHandler(logger.Info, testHandler)
	ch := h.(io.Closer)

	ch.Close()

	assert.True(t, testHandler.CloseCalled)
}

func TestDiscardHandler(t *testing.T) {
	h := logger.DiscardHandler()

	h.Log(&logger.Event{Msg: "test", Lvl: logger.Crit})
}
