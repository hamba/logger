package logger

import "context"

type contextKey struct{}

type ctxFields struct {
	b []byte
}

// WithContext returns a copy of ctx carrying the given fields.
// Fields already attached to ctx by WithContext are kept and the new fields
// follow them. Fields added to log with [Logger.With] are not attached.
//
// The fields are evaluated once, when WithContext is called, and formatted
// with the [Formatter] of log. They must only be retrieved with
// [Logger.FromContext] on a [Logger] that uses the same kind of Formatter.
//
// If log discards output or no fields are given, ctx is returned unchanged.
func WithContext(ctx context.Context, log *Logger, fields ...Field) context.Context {
	if log.isDiscard || len(fields) == 0 {
		return ctx
	}

	e := newEvent(log.fmtr)
	defer putEvent(e)

	if existing, _ := ctx.Value(contextKey{}).(*ctxFields); existing != nil {
		e.buf.Write(existing.b)
	}

	for _, field := range fields {
		field(e)
	}

	b := make([]byte, e.buf.Len())
	copy(b, e.buf.Bytes())

	return context.WithValue(ctx, contextKey{}, &ctxFields{b: b})
}

// FromContext returns a [Logger] extended with the fields attached to ctx by
// [WithContext]. The fields follow those added with [Logger.With] and precede
// those passed to the log methods.
//
// If l discards output or ctx carries no fields, FromContext returns l.
// The receiver is never modified.
func (l *Logger) FromContext(ctx context.Context) *Logger {
	if l.isDiscard {
		return l
	}

	fields, ok := ctx.Value(contextKey{}).(*ctxFields)
	if !ok || len(fields.b) == 0 {
		return l
	}

	b := make([]byte, len(l.ctx)+len(fields.b))
	copy(b, l.ctx)
	copy(b[len(l.ctx):], fields.b)

	return &Logger{
		w:         l.w,
		isDiscard: l.isDiscard,
		fmtr:      l.fmtr,
		timeFn:    l.timeFn,
		lvl:       l.lvl,
		ctx:       b,
	}
}
