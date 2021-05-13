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

func (e *Event) AppendString(k, v string) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendString(e.buf, v)
}

func (e *Event) AppendBool(k string, b bool) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendBool(e.buf, b)
}

func (e *Event) AppendInt(k string, i int64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendInt(e.buf, i)
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
