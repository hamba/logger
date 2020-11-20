package logger

import "sync"

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{}
	},
}

// Event is a log event.
type Event struct {
	Time    int64
	Msg     string
	Lvl     Level
	BaseCtx []interface{}
	Ctx     []interface{}
}

func newEvent(msg string, lvl Level) *Event {
	e := eventPool.Get().(*Event)
	e.Msg = msg
	e.Lvl = lvl
	return e
}

func putEvent(e *Event) {
	eventPool.Put(e)
}
