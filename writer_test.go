package ctxio

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	content := "The ctxio package gives golang io.copy operations the ability to terminate with context and retrieve progress data."

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tempDir := t.TempDir()

	srcFile, err := ioutil.TempFile(tempDir, "source-file-*")
	t.Logf("srcFile: %s", srcFile.Name())
	require.NoError(t, err)
	srcFile.WriteString("The ctxio package gives golang io.copy operations the ability to terminate with context and retrieve progress data.")
	srcFile.Sync()
	defer srcFile.Close()

	dstFile, err := ioutil.TempFile(tempDir, "destination-file-*")
	t.Logf("dstFile: %s", dstFile.Name())
	require.NoError(t, err)
	defer dstFile.Close()

	file, _ := os.Open(srcFile.Name())

	progressFn := func(n int64) {
		t.Logf("Progress bytes: %d", n)
	}

	w := NewWriter(ctx, dstFile, progressFn)

	written, err := io.Copy(w, file)
	t.Logf("written bytes: %d", written)
	require.NoError(t, err)

	expectedContent, err := os.ReadFile(dstFile.Name())
	require.NoError(t, err)

	require.Equal(t, content, string(expectedContent))
}

func TestNewWriter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := []byte("binalyze")

	var progress int64
	progressFn := func(n int64) {
		t.Logf("written bytes: %d", n)
		progress = n
	}

	b := bytes.NewBuffer(nil)
	w := NewWriter(ctx, b, progressFn)

	// Write first half
	n, err := w.Write(buf[:4]) // b = "bina"
	require.NoError(t, err)
	require.NoError(t, w.Err())
	require.Equal(t, 4, n)
	require.Equal(t, int64(4), w.N())
	require.Equal(t, b.String(), string(buf[:4]))
	require.Equal(t, int64(4), progress)

	// Write second half
	n, err = w.Write(buf[4:8]) // b = "lyze"
	require.NoError(t, err)
	require.NoError(t, w.Err())
	require.Equal(t, 4, n)
	require.Equal(t, int64(8), w.N())
	require.Equal(t, b.String(), string(buf))
	require.Equal(t, int64(8), progress)
}

func TestContextCancel_Writer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := []byte("binalyze")

	b := bytes.NewBuffer(nil)

	var progress int64
	progressFn := func(n int64) {
		t.Logf("written bytes: %d", n)
		progress = n
	}

	w := NewWriter(ctx, b, progressFn)

	// Write first half
	n, err := w.Write(buf[:4]) // b = "bina"
	require.NoError(t, err)
	require.NoError(t, w.Err())
	require.Equal(t, 4, n)
	require.Equal(t, b.String(), string(buf[:4]))
	require.Equal(t, int64(4), progress)

	cancel()

	n, err = w.Write(buf[4:8]) // b = "lyze"
	require.Error(t, context.Canceled, err)
	require.Error(t, context.Canceled, w.Err())
	require.Equal(t, 0, n)
	require.Equal(t, int64(4), progress)
}
