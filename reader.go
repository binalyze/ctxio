package ctxio

import (
	"context"
	"io"
	"sync"
)

type Reader struct {
	reader     io.Reader
	ctx        context.Context
	progressFn func(int642 int64)
	mu         sync.RWMutex
	n          int64
	err        error
}

func noonProgressFn(int64) {}

// NewReader creates a new Reader with the given io.Reader
func NewReader(ctx context.Context, r io.Reader, progressFn func(int642 int64)) *Reader {
	if ctx == nil {
		ctx = context.Background()
	}

	if progressFn == nil {
		progressFn = noonProgressFn
	}

	return &Reader{
		reader:     r,
		ctx:        ctx,
		progressFn: progressFn,
	}
}

// Read reads the contents into the underlying io.Reader.
// Then it increments the number of bytes read.
func (r *Reader) Read(p []byte) (int, error) {
	if r.ctx.Err() != nil {
		return 0, r.ctx.Err()
	}

	n, err := r.reader.Read(p)

	r.mu.RLock()
	defer r.mu.RUnlock()

	r.n += int64(n)
	r.err = err
	r.progressFn(r.n)

	return n, err
}

// N returns the number of bytes read.
func (r *Reader) N() int64 {
	r.mu.RLock()
	n := r.n
	r.mu.RUnlock()

	return n
}

// Err gets the last error from the Reader.
func (r *Reader) Err() error {
	r.mu.RLock()
	err := r.err
	r.mu.RUnlock()

	return err
}
