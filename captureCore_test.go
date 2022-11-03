package sallust

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCaptureCore(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		// nolint:staticcheck
		enc = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:  "msg",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		})

		buffer Buffer
		cc     = NewCaptureCore(enc, &buffer, zapcore.InfoLevel)
		l      = zap.New(cc)
	)

	assert.Nil(cc.Check(zapcore.Entry{Level: zapcore.DebugLevel}, nil))
	assert.NotNil(cc.Check(zapcore.Entry{Level: zapcore.InfoLevel}, nil))

	assert.Zero(cc.Len())
	n, err := cc.EachMessage(func(zapcore.Entry, []zapcore.Field) error {
		return errors.New("EachMessage should not have been called")
	})
	assert.Zero(n)
	assert.NoError(err)

	l.Debug("this should not appear")
	n, err = cc.EachMessage(func(zapcore.Entry, []zapcore.Field) error {
		return errors.New("EachMessage should not have been called")
	})
	assert.Zero(n)
	assert.NoError(err)

	cc.ClearMessages()
	n, err = cc.EachMessage(func(zapcore.Entry, []zapcore.Field) error {
		return errors.New("EachMessage should not have been called")
	})
	assert.Zero(n)
	assert.NoError(err)

	l.Info("message", zap.Int("foo", 123))
	l.Sync()
	assert.Equal(1, cc.Len())
	n, err = cc.EachMessage(func(e zapcore.Entry, f []zapcore.Field) error {
		assert.Equal("message", e.Message)
		require.Len(f, 1)
		assert.Equal("foo", f[0].Key)
		assert.Equal(int64(123), f[0].Integer)
		return nil
	})

	assert.Equal(1, n)
	assert.NoError(err)

	l = l.With(zap.String("bar", "asdf"))
	l.Info("another message", zap.Int("foo", 456))
	l.Sync()
	require.NotPanics(func() {
		cc = l.Core().(*CaptureCore)
	})

	assert.Equal(2, cc.Len())
	n, err = cc.EachMessage(func(e zapcore.Entry, f []zapcore.Field) error {
		if e.Message == "message" {
			require.Len(f, 1)
			assert.Equal("foo", f[0].Key)
			assert.Equal(int64(123), f[0].Integer)
		} else if e.Message == "another message" {
			require.Len(f, 2)
			assert.Equal("foo", f[0].Key)
			assert.Equal(int64(456), f[0].Integer)
			assert.Equal("bar", f[1].Key)
			assert.Equal("asdf", f[1].String)
		} else {
			return errors.New("Unrecognized captured message")
		}

		return nil
	})

	assert.Equal(2, n)
	assert.NoError(err)

	n, err = cc.EachMessage(func(e zapcore.Entry, _ []zapcore.Field) error {
		if e.Message != "another message" {
			return nil
		}

		return errors.New("expected")
	})

	assert.Equal(1, n)
	assert.Error(err)
}

func TestCapture(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		enc     = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:  "msg",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		})

		buffer Buffer
		c      = Capture(
			zapcore.NewCore(enc, &buffer, zapcore.InfoLevel),
		)
	)

	require.NotNil(c)

	c.Write(zapcore.Entry{Level: zapcore.InfoLevel}, []zapcore.Field{})
	assert.True(buffer.Len() > 0)

	cc, ok := c.(*CaptureCore)
	require.True(ok)
	assert.Equal(1, cc.Len())
}
