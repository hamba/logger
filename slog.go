package logger

import (
	"context"
	"io"
	"log/slog"
)

// Handler is a slog.Handler backed by hamba/logger. It does not synchronize writes;
// wrap the writer with NewSyncWriter for concurrent use in non-posix environments.
type Handler struct {
	w         io.Writer
	isDiscard bool
	fmtr      Formatter
	lvl       slog.Level

	ctx    []byte
	prefix []byte

	groups []string
}

// NewHandler returns a new Handler.
func NewHandler(w io.Writer, fmtr Formatter, lvl slog.Level) *Handler {
	isDiscard := w == io.Discard

	return &Handler{
		w:         w,
		isDiscard: isDiscard,
		fmtr:      fmtr,
		lvl:       lvl,
	}
}

// Enabled returns false when lvl is below the configured minimum or the
// underlying writer is io.Discard.
func (h *Handler) Enabled(_ context.Context, lvl slog.Level) bool {
	return !h.isDiscard && lvl >= h.lvl
}

// Handle writes the record to the underlying writer.
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	e := newEvent(h.fmtr)
	e.prefix = append(e.prefix, h.prefix...)

	e.fmtr.AppendBeginMarker(e.buf)
	e.fmtr.WriteMessage(e.buf, r.Time, mapSlogLevel(r.Level), r.Message)
	e.buf.Write(h.ctx)

	r.Attrs(func(a slog.Attr) bool {
		appendAttr(e, a)
		return true
	})

	for range h.groups {
		e.prefix = e.fmtr.AppendGroupEnd(e.buf, e.prefix)
	}

	e.fmtr.AppendEndMarker(e.buf)
	e.fmtr.AppendLineBreak(e.buf)

	_, err := h.w.Write(e.buf.Bytes())
	putEvent(e)
	return err
}

// WithAttrs returns a new Handler with attrs pre-serialised into the context.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	e := newEvent(h.fmtr)
	e.buf.Write(h.ctx)
	e.prefix = append(e.prefix, h.prefix...)

	for _, a := range attrs {
		appendAttr(e, a)
	}

	newCtx := make([]byte, e.buf.Len())
	copy(newCtx, e.buf.Bytes())

	putEvent(e)

	return &Handler{
		fmtr:      h.fmtr,
		w:         h.w,
		isDiscard: h.isDiscard,
		lvl:       h.lvl,
		ctx:       newCtx,
		prefix:    h.prefix,
		groups:    h.groups,
	}
}

// WithGroup returns a new Handler with name appended to the group stack.
// An empty name is a no-op per the slog spec.
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	groups := make([]string, len(h.groups)+1)
	copy(groups, h.groups)
	groups[len(h.groups)] = name

	e := newEvent(h.fmtr)
	e.buf.Write(h.ctx)
	e.prefix = append(e.prefix, h.prefix...)

	e.prefix = h.fmtr.AppendGroupStart(e.buf, e.prefix, name)

	newCtx := make([]byte, e.buf.Len())
	copy(newCtx, e.buf.Bytes())
	newPrefix := make([]byte, len(e.prefix))
	copy(newPrefix, e.prefix)

	putEvent(e)

	return &Handler{
		fmtr:      h.fmtr,
		w:         h.w,
		isDiscard: h.isDiscard,
		lvl:       h.lvl,
		ctx:       newCtx,
		prefix:    newPrefix,
		groups:    groups,
	}
}

func appendAttr(e *Event, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}

	switch a.Value.Kind() {
	case slog.KindGroup:
		appendGroup(e, a.Key, a.Value)
	case slog.KindString:
		e.AppendString(a.Key, a.Value.String())
	case slog.KindInt64:
		e.AppendInt(a.Key, a.Value.Int64())
	case slog.KindUint64:
		e.AppendUint(a.Key, a.Value.Uint64())
	case slog.KindFloat64:
		e.AppendFloat(a.Key, a.Value.Float64())
	case slog.KindBool:
		e.AppendBool(a.Key, a.Value.Bool())
	case slog.KindTime:
		e.AppendTime(a.Key, a.Value.Time())
	case slog.KindDuration:
		e.AppendDuration(a.Key, a.Value.Duration())
	default:
		e.AppendInterface(a.Key, a.Value.Any())
	}
}

func appendGroup(e *Event, name string, val slog.Value) {
	subs := val.Group()
	if len(subs) == 0 {
		return
	}

	if name == "" {
		// Per the slog spec, an anonymous group is flattened into the
		// enclosing scope.
		for _, a := range subs {
			appendAttr(e, a)
		}
		return
	}

	e.OpenGroup(name)
	for _, a := range subs {
		appendAttr(e, a)
	}
	e.CloseGroup()
}

// mapSlogLevel maps to logger.Level. Custom levels below Debug clamp to
// Debug; above Error clamp to Error. Trace and Crit are never produced.
func mapSlogLevel(lvl slog.Level) Level {
	switch {
	case lvl >= slog.LevelError:
		return Error
	case lvl >= slog.LevelWarn:
		return Warn
	case lvl >= slog.LevelInfo:
		return Info
	default:
		return Debug
	}
}
