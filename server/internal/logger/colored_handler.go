package logger

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"sync"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

type ColoredHandler struct {
	inner  slog.Handler
	writer io.Writer
	opts   *slog.HandlerOptions
	attrs  []slog.Attr
	groups []string
	mu     sync.Mutex
}

func NewColoredHandler(w io.Writer, opts *slog.HandlerOptions) *ColoredHandler {
	return &ColoredHandler{
		inner:  slog.NewTextHandler(w, opts),
		writer: w,
		opts:   opts,
	}
}

func (h *ColoredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *ColoredHandler) Handle(ctx context.Context, r slog.Record) error {
	var buf bytes.Buffer
	var temp slog.Handler = slog.NewTextHandler(&buf, h.opts)
	for _, attr := range h.attrs {
		temp = temp.WithAttrs([]slog.Attr{attr})
	}
	for _, group := range h.groups {
		temp = temp.WithGroup(group)
	}
	if err := temp.Handle(ctx, r); err != nil {
		return err
	}
	line := buf.String()
	switch r.Level {
	case slog.LevelWarn:
		line = colorYellow + line + colorReset
	case slog.LevelError:
		line = colorRed + line + colorReset
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.writer.Write([]byte(line))
	return err
}

func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColoredHandler{
		inner:  h.inner.WithAttrs(attrs),
		writer: h.writer,
		opts:   h.opts,
		attrs:  append(append([]slog.Attr(nil), h.attrs...), attrs...),
		groups: h.groups,
	}
}

func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	return &ColoredHandler{
		inner:  h.inner.WithGroup(name),
		writer: h.writer,
		opts:   h.opts,
		attrs:  h.attrs,
		groups: append(append([]string(nil), h.groups...), name),
	}
}
