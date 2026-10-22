package logger

import (
	"sync"
	"time"

	"github.com/hamba/logger/v2/internal/bytes"
)

var eventPool = &sync.Pool{
	New: func() any {
		return &Event{
			buf: bytes.NewBuffer(512),
		}
	},
}

// Event is the log message that a [Field] adds its key-value pairs to.
//
// An Event is only valid during the call to the Field that receives it and
// must not be retained.
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

// AppendString appends the string s under key k.
func (e *Event) AppendString(k, s string) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendString(e.buf, s)
}

// AppendStrings appends the strings s as an array under key k.
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

// AppendBytes appends p as an array of integers under key k, not as a string.
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

// AppendBool appends b under key k.
func (e *Event) AppendBool(k string, b bool) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendBool(e.buf, b)
}

// AppendInt appends i under key k.
func (e *Event) AppendInt(k string, i int64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendInt(e.buf, i)
}

// AppendInts appends a as an array of integers under key k.
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

// AppendUint appends i under key k.
func (e *Event) AppendUint(k string, i uint64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendUint(e.buf, i)
}

// AppendFloat appends f under key k.
func (e *Event) AppendFloat(k string, f float64) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendFloat(e.buf, f)
}

// AppendTime appends t under key k, rendered according to [TimeFormat] or
// [TimeFormatConsole].
func (e *Event) AppendTime(k string, t time.Time) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendTime(e.buf, t)
}

// AppendDuration appends d under key k.
func (e *Event) AppendDuration(k string, d time.Duration) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendDuration(e.buf, d)
}

// AppendInterface appends the text form of v under key k.
func (e *Event) AppendInterface(k string, v any) {
	e.fmtr.AppendKey(e.buf, k)
	e.fmtr.AppendInterface(e.buf, v)
}
