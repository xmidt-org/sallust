package sallusthttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
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

		core, logs = observer.New(zapcore.DebugLevel)
		base       = zap.New(core)
		request    = httptest.NewRequest("PUT", "/test", nil)
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
	assert.Equal(1, logs.Len())
	for _, le := range logs.TakeAll() {
		assert.ElementsMatch(expected, le.Context)
		assert.Equal("testHandler", le.Entry.LoggerName)
	}
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
