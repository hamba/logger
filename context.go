package logger

import "context"

type contextKey struct{}

type ctxFields struct {
	b []byte
}

// WithContext returns a new context carrying the given fields.
// Any fields previously attached via WithContext are preserved;
// the new fields are appended after them. Any fields already
// added to the logger are not attached to the context.
//
// If the logger discards output or no fields are given, ctx is returned
// unchanged.
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

// FromContext returns a new Logger extended with any fields attached to ctx
// via WithContext. The context fields appear after the logger's own pre-rendered
// fields (set via With) and before per-call fields.
//
// If the logger discards output or ctx carries no fields, the receiver is
// returned unchanged.
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
