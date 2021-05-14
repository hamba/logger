package logger

import (
	"sync"
	"time"

	"github.com/hamba/logger/internal/bytes"
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

func (e *Event) AppendString(k, s string) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendString(e.buf, s)
}

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

func (e *Event) AppendBool(k string, b bool) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendBool(e.buf, b)
}

func (e *Event) AppendInt(k string, i int64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendInt(e.buf, i)
}

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

func (e *Event) AppendUint(k string, i uint64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendUint(e.buf, i)
}

func (e *Event) AppendFloat(k string, f float64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendFloat(e.buf, f)
}

func (e *Event) AppendTime(k string, d time.Time) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendTime(e.buf, d)
}

func (e *Event) AppendDuration(k string, d time.Duration) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendDuration(e.buf, d)
}

func (e *Event) AppendInterface(k string, v interface{}) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendInterface(e.buf, v)
}
