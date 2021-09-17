package ctxio

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewReader(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := []byte("binalyze")
	rBuf := make([]byte, 4)

	var progress int64
	progressFn := func(n int64) {
		t.Logf("read bytes: %d", n)
		progress = n
	}

	b := bytes.NewBuffer(buf)
	r := NewReader(ctx, b, progressFn)

	// Read first half
	n, err := r.Read(rBuf) // rBuf = "bina"

	require.NoError(t, err)
	require.NoError(t, r.Err())
	require.Equal(t, 4, n)
	require.Equal(t, int64(4), r.N())
	require.Equal(t, string(rBuf), string(buf[:4]))
	require.Equal(t, int64(4), progress)

	// Read second half
	n, err = r.Read(rBuf) // rBuf = "lyze"
	require.NoError(t, err)
	require.NoError(t, r.Err())
	require.Equal(t, 4, n)
	require.Equal(t, int64(8), r.N())
	require.Equal(t, string(rBuf), string(buf[4:]))
	require.Equal(t, int64(8), progress)
}

func TestContextCancel_Reader(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := []byte("binalyze")
	rBuf := make([]byte, 4)

	var progress int64
	progressFn := func(n int64) {
		t.Logf("read bytes: %d", n)
		progress = n
	}

	b := bytes.NewBuffer(buf)
	r := NewReader(ctx, b, progressFn)

	// Read first half
	n, err := r.Read(rBuf) // rBuf = "bina"

	require.NoError(t, err)
	require.NoError(t, r.Err())
	require.Equal(t, 4, n)
	require.Equal(t, int64(4), r.N())
	require.Equal(t, string(rBuf), string(buf[:4]))
	require.Equal(t, int64(4), progress)

	cancel()

	n, err = r.Read(rBuf) // rBuf = "lyze"
	require.Error(t, context.Canceled, err)
	require.Error(t, context.Canceled, r.Err())
	require.Equal(t, 0, n)
	require.Equal(t, int64(4), progress)
}

func TestContextTimeoutCancel_Reader(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	buf := []byte("binalyze")
	rBuf := make([]byte, 4)

	var progress int64
	progressFn := func(n int64) {
		t.Logf("read bytes: %d", n)
		progress = n
	}

	b := bytes.NewBuffer(buf)
	r := NewReader(ctx, b, progressFn)

	// Read first half
	n, err := r.Read(rBuf) // rBuf = "bina"

	require.NoError(t, err)
	require.NoError(t, r.Err())
	require.Equal(t, 4, n)
	require.Equal(t, int64(4), r.N())
	require.Equal(t, string(rBuf), string(buf[:4]))
	require.Equal(t, int64(4), progress)

	time.Sleep(2 * time.Second)

	n, err = r.Read(rBuf) // rBuf = "lyze"
	require.Error(t, context.DeadlineExceeded, err)
	require.Error(t, context.DeadlineExceeded, r.Err())
	require.Equal(t, 0, n)
	require.Equal(t, int64(4), progress)
}
