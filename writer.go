package ctxio

import (
	"context"
	"io"
	"sync"
)

// Writer holds the necessary state to count the number
type Writer struct {
	writer     io.Writer
	ctx        context.Context
	progressFn func(int64)
	mu         sync.RWMutex
	n          int64
	err        error
}

func noopProgressFn(int64) {}

// NewWriter creates a new Writer with the given io.Writer.
func NewWriter(ctx context.Context, w io.Writer, progressFn func(int64)) *Writer {
	if ctx == nil {
		ctx = context.Background()
	}

	if progressFn == nil {
		progressFn = noopProgressFn
	}

	return &Writer{
		writer:     w,
		ctx:        ctx,
		progressFn: progressFn,
	}
}

// Write writes the contents of p into the underlying io.Writer.
// Then it increments the number of bytes written.
func (w *Writer) Write(p []byte) (int, error) {
	if w.ctx.Err() != nil {
		return 0, w.ctx.Err()
	}

	n, err := w.writer.Write(p)

	w.mu.Lock()
	defer w.mu.Unlock()

	w.n += int64(n)
	w.err = err
	w.progressFn(w.n)

	return n, err
}

// N returns the number of bytes written.
func (w *Writer) N() int64 {
	w.mu.RLock()
	n := w.n
	w.mu.RUnlock()

	return n
}

// Err gets the last error from the Writer.
func (w *Writer) Err() error {
	w.mu.RLock()
	err := w.err
	w.mu.RUnlock()

	return err
}
