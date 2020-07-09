package sallusthttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestBuilders(t *testing.T) {
	var (
		assert   = assert.New(t)
		require  = require.New(t)
		expected = []zapcore.Field{
			zap.String(DefaultMethodKey, "PUT"),
			zap.String(DefaultRemoteAddrKey, "127.0.0.1"),
			zap.String(DefaultURIKey, "/test"),
			zap.String("testy", "mctest"),
		}

		buffer sallust.Buffer
		enc    = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			// omit normal keys for easier comparisons
			NameKey: "name",
		})

		core    = sallust.NewCaptureCore(enc, &buffer, zapcore.DebugLevel)
		base    = zap.New(core)
		request = httptest.NewRequest("PUT", "/test", nil)
	)

	request.RemoteAddr = "127.0.0.1"

	var b Builders
	b.Add(Named("testHandler"), DefaultFields)
	b.AddFields(
		func(r *http.Request, f []zap.Field) []zap.Field {
			return append(f, zap.String("testy", "mctest"))
		},
	)

	l := b.Build(request, base)
	require.NotNil(l)

	l.Info("test")
	l.Sync()

	// need to pull out the decorated CaptureCore
	decorated, ok := l.Core().(*sallust.CaptureCore)
	require.True(ok)
	n, err := decorated.EachMessage(func(e zapcore.Entry, actual []zapcore.Field) error {
		assert.ElementsMatch(expected, actual)
		assert.Equal("testHandler", e.LoggerName)
		return nil
	})

	assert.Equal(1, n)
	assert.NoError(err)
}

func TestMethod(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		request = httptest.NewRequest("PUT", "/test", nil)
	)

	f := Method(request, nil)
	require.Len(f, 1)
	assert.Equal(
		zap.String(DefaultMethodKey, "PUT"),
		f[0],
	)
}

func TestMethodCustom(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		request = httptest.NewRequest("PUT", "/test", nil)
	)

	f := MethodCustom("test")(request, nil)
	require.Len(f, 1)
	assert.Equal(
		zap.String("test", "PUT"),
		f[0],
	)
}

func TestURI(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		request = httptest.NewRequest("PUT", "/test", nil)
	)

	f := URI(request, nil)
	require.Len(f, 1)
	assert.Equal(
		zap.String(DefaultURIKey, "/test"),
		f[0],
	)
}

func TestURICustom(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		request = httptest.NewRequest("PUT", "/test", nil)
	)

	f := URICustom("test")(request, nil)
	require.Len(f, 1)
	assert.Equal(
		zap.String("test", "/test"),
		f[0],
	)
}

func TestRemoteAddr(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		request = httptest.NewRequest("PUT", "/test", nil)
	)

	request.RemoteAddr = "172.3.4.5"
	f := RemoteAddr(request, nil)
	require.Len(f, 1)
	assert.Equal(
		zap.String(DefaultRemoteAddrKey, "172.3.4.5"),
		f[0],
	)
}

func TestRemoteAddrCustom(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		request = httptest.NewRequest("PUT", "/test", nil)
	)

	request.RemoteAddr = "19.56.71.123"
	f := RemoteAddrCustom("test")(request, nil)
	require.Len(f, 1)
	assert.Equal(
		zap.String("test", "19.56.71.123"),
		f[0],
	)
}
