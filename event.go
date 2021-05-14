package logger

import (
	"sync"
	"time"

	"github.com/hamba/logger/v2/internal/bytes"
)

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{
			buf: bytes.NewBuffer(512),
		}
	},
}

// Event is a log event.
type Event struct {
	fmtr Formatter
	buf  *bytes.Buffer
}

func newEvent(fmtr Formatter) *Event {
	e := eventPool.Get().(*Event)
	e.fmtr = fmtr
	e.buf.Reset()
	return e
}

func putEvent(e *Event) {
	eventPool.Put(e)
}

// AppendString appends a string to the event.
func (e *Event) AppendString(k, s string) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendString(e.buf, s)
}

// AppendStrings appends strings to the event.
func (e *Event) AppendStrings(k string, s []string) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendArrayStart(e.buf)
	for i, ss := range s {
		if i > 0 {
			e.fmtr.AppendArraySep(e.buf)
		}
		e.fmtr.AppendString(e.buf, ss)
	}
	e.fmtr.AppendArrayEnd(e.buf)
}

// AppendBytes appends bytes to the event.
func (e *Event) AppendBytes(k string, p []byte) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendArrayStart(e.buf)
	for i, b := range p {
		if i > 0 {
			e.fmtr.AppendArraySep(e.buf)
		}
		e.fmtr.AppendInt(e.buf, int64(b))
	}
	e.fmtr.AppendArrayEnd(e.buf)
}

// AppendBool appends a bool to the event.
func (e *Event) AppendBool(k string, b bool) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendBool(e.buf, b)
}

// AppendInt appends an int to the event.
func (e *Event) AppendInt(k string, i int64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendInt(e.buf, i)
}

// AppendInts appends ints to the event.
func (e *Event) AppendInts(k string, a []int) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendArrayStart(e.buf)
	for i, ii := range a {
		if i > 0 {
			e.fmtr.AppendArraySep(e.buf)
		}
		e.fmtr.AppendInt(e.buf, int64(ii))
	}
	e.fmtr.AppendArrayEnd(e.buf)
}

// AppendUint appends a uint to the event.
func (e *Event) AppendUint(k string, i uint64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendUint(e.buf, i)
}

// AppendFloat appends a float to the event.
func (e *Event) AppendFloat(k string, f float64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendFloat(e.buf, f)
}

// AppendTime appends a time to the event.
func (e *Event) AppendTime(k string, d time.Time) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendTime(e.buf, d)
}

// AppendDuration appends a duration to the event.
func (e *Event) AppendDuration(k string, d time.Duration) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendDuration(e.buf, d)
}

// AppendInterface appends a interface to the event.
func (e *Event) AppendInterface(k string, v interface{}) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendInterface(e.buf, v)
}
