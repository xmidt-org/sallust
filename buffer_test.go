package sallust

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testBufferMessage = "this is a lovely test message"

func testBufferInitial(t *testing.T) {
	var (
		assert  = assert.New(t)
		writeTo = new(bytes.Buffer)
		buffer  Buffer
	)

	assert.Zero(buffer.Len())
	assert.Zero(buffer.Limit())
	assert.NoError(buffer.Sync())
	assert.NoError(buffer.Close())
	assert.Empty(buffer.String())

	n64, err := buffer.WriteTo(writeTo)
	assert.Zero(n64)
	assert.NoError(err)
	assert.Zero(writeTo.Len())
}

func testBufferReset(t *testing.T) {
	var (
		assert  = assert.New(t)
		writeTo = new(bytes.Buffer)
		buffer  Buffer
	)

	buffer.Write([]byte(testBufferMessage))
	buffer.Reset()
	assert.Empty(buffer.String())

	n64, err := buffer.WriteTo(writeTo)
	assert.Zero(n64)
	assert.NoError(err)
	assert.Zero(buffer.Len())
}

func testBufferSetLimit(t *testing.T) {
	var (
		assert = assert.New(t)
		buffer Buffer
	)

	buffer.Write([]byte(testBufferMessage))
	buffer.SetLimit(2 * len(testBufferMessage))
	assert.Equal(len(testBufferMessage), buffer.Len())
	assert.Equal(testBufferMessage, buffer.String())

	buffer.SetLimit(len(testBufferMessage) - 1)
	assert.Zero(buffer.Len())
	assert.Empty(buffer.String())

	buffer.SetLimit(0)
	buffer.Write([]byte(testBufferMessage))
	assert.Equal(len(testBufferMessage), buffer.Len())
	assert.Equal(testBufferMessage, buffer.String())
}

func testBufferWriteNoLimit(t *testing.T) {
	var (
		assert  = assert.New(t)
		writeTo = new(bytes.Buffer)
		buffer  Buffer
	)

	n, err := buffer.Write([]byte(testBufferMessage))
	assert.Equal(len(testBufferMessage), n)
	assert.NoError(err)
	assert.Equal(len(testBufferMessage), buffer.Len())
	assert.Zero(buffer.Limit())
	assert.NoError(buffer.Sync())
	assert.NoError(buffer.Close())
	assert.Equal(testBufferMessage, buffer.String())

	n64, err := buffer.WriteTo(writeTo)
	assert.Equal(int64(len(testBufferMessage)), n64)
	assert.NoError(err)
	assert.Equal(len(testBufferMessage), writeTo.Len())
	writeTo.Reset()
}

func testBufferWriteLimit(t *testing.T) {
	var (
		assert  = assert.New(t)
		writeTo = new(bytes.Buffer)
		buffer  Buffer
	)

	buffer.SetLimit(len(testBufferMessage))
	assert.Equal(len(testBufferMessage), buffer.Limit())

	n, err := buffer.Write([]byte(testBufferMessage))
	assert.Equal(len(testBufferMessage), n)
	assert.NoError(err)
	assert.Equal(len(testBufferMessage), buffer.Len())
	assert.Equal(len(testBufferMessage), buffer.Limit())
	assert.NoError(buffer.Sync())
	assert.NoError(buffer.Close())
	assert.Equal(testBufferMessage, buffer.String())

	// write a second time, to violate the limit
	n, err = buffer.Write([]byte(testBufferMessage))
	assert.Equal(len(testBufferMessage), n)
	assert.NoError(err)
	assert.Equal(len(testBufferMessage), buffer.Len())
	assert.Equal(len(testBufferMessage), buffer.Limit())
	assert.NoError(buffer.Sync())
	assert.NoError(buffer.Close())
	assert.Equal(testBufferMessage, buffer.String())

	n64, err := buffer.WriteTo(writeTo)
	assert.Equal(int64(len(testBufferMessage)), n64)
	assert.NoError(err)
	assert.Equal(len(testBufferMessage), writeTo.Len())
	assert.Equal(len(testBufferMessage), buffer.Limit())
	writeTo.Reset()
}

func TestBuffer(t *testing.T) {
	t.Run("Initial", testBufferInitial)
	t.Run("Reset", testBufferReset)
	t.Run("SetLimit", testBufferSetLimit)
	t.Run("Write", func(t *testing.T) {
		t.Run("NoLimit", testBufferWriteNoLimit)
		t.Run("Limit", testBufferWriteLimit)
	})
}
