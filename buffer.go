package sallust

import (
	"bytes"
	"io"
	"sync"

	"go.uber.org/zap"
)

// Buffer is a zap.Sink that captures all logging to an in-memory buffer.
// An optional max size can be placed on the buffer, at which time the buffer
// is cleared before writing more log information.  All methods of a Buffer
// may be called concurrently.  The zero value of a Buffer is valid, and is
// an unlimited size in-memory sink for logs.
//
// Deprecated:  This will be removed in a future release.  Use the zaptest package instead.
type Buffer struct {
	lock   sync.Mutex
	buffer bytes.Buffer
	limit  int
}

var _ zap.Sink = (*Buffer)(nil)

// Len returns the current size of the internal buffer
func (b *Buffer) Len() (n int) {
	b.lock.Lock()
	n = b.buffer.Len()
	b.lock.Unlock()
	return
}

// Limit returns the current size limit of the buffer, in bytes.
// A nonpositive value indicates no limit.
func (b *Buffer) Limit() (l int) {
	b.lock.Lock()
	l = b.limit
	b.lock.Unlock()
	return
}

// SetLimit changes the size limit of this buffer.  If the new limit
// is smaller than the current size of the buffer, the buffer is reset.
// A nonpositive value indicates no limit.
func (b *Buffer) SetLimit(l int) {
	b.lock.Lock()

	if l < 1 {
		b.limit = 0
	} else {
		b.limit = l
		if b.buffer.Len() > b.limit {
			b.buffer.Reset()
		}
	}

	b.lock.Unlock()
}

// Write appends log output to the buffer.  If the output would grow the buffer
// beyond its limit, the buffer is cleared first.  If the limit is smaller than
// the length of p, the write proceeds anyway because a client may still want
// to capture the log output.
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.lock.Lock()
	if b.limit > 0 && b.buffer.Len()+len(p) > b.limit {
		b.buffer.Reset()
	}

	n, err = b.buffer.Write(p)
	b.lock.Unlock()
	return
}

// WriteTo writes the current buffer's contents to the supplied writer.
// The buffer's contents are reset after writing.
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	b.lock.Lock()
	n, err = b.buffer.WriteTo(w)
	b.lock.Unlock()
	return
}

// Sync is a nop, as there is no underlying I/O done by a Buffer
func (b *Buffer) Sync() error {
	return nil
}

// Reset resets the internal buffer, retaining storage for future writes
func (b *Buffer) Reset() {
	b.lock.Lock()
	b.buffer.Reset()
	b.lock.Unlock()
}

// Close is a nop, as no underlying I/O is performed by a Buffer
func (b *Buffer) Close() error {
	return nil
}

// String returns the current buffer's log contents.
func (b *Buffer) String() (s string) {
	b.lock.Lock()
	s = b.buffer.String()
	b.lock.Unlock()
	return
}
