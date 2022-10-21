package sallusthttp

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
)

func TestWith(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	logger := new(zap.Logger)
	original := httptest.NewRequest("PUT", "/test", nil)

	child := With(original, logger)
	require.NotNil(child)
	assert.Equal(original.Method, child.Method)
	assert.Equal(original.URL.String(), child.URL.String())
	ctx := child.Context()
	assert.Equal(logger, sallust.Get(ctx))
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	original := httptest.NewRequest("GET", "/test", nil)

	logger := Get(original)
	assert.Equal(sallust.Default(), logger)

	def := new(zap.Logger)
	logger = GetDefault(original, def)
	assert.Equal(def, logger)

	logger = new(zap.Logger)
	child := With(original, logger)
	assert.Equal(logger, Get(child))
}

func TestRequestInfo(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	request := httptest.NewRequest("GET", "/test/foo/bar", nil)
	request.RemoteAddr = "127.0.0.1:1234"
	kv := RequestInfo(nil, request)

	require.NotNil(kv)
	require.Len(kv, 3)
	assert.True(kv[0].Equals(zap.String(requestMethodKey, "GET")))
	assert.True(kv[1].Equals(zap.String(requestURIKey, "/test/foo/bar")))
	assert.True(kv[2].Equals(zap.String(remoteAddrKey, "127.0.0.1:1234")))
}

func testHeaderMissing(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	request := httptest.NewRequest("GET", "/", nil)
	kv := Header("X-Test", "key")(nil, request)

	require.NotNil(kv)
	require.Len(kv, 1)
	assert.True(kv[0].Equals(zap.String("key", "")))
}

func testHeaderSingleValue(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	request := httptest.NewRequest("GET", "/", nil)

	request.Header.Set("X-Test", "value")
	kv := Header("X-Test", "key")(nil, request)
	require.NotNil(kv)
	require.Len(kv, 1)
	assert.True(kv[0].Equals(zap.String("key", "value")))
}

func testHeaderMultiValue(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	request := httptest.NewRequest("GET", "/", nil)

	request.Header.Add("X-Test", "value1")
	request.Header.Add("X-Test", "value2")
	kv := Header("X-Test", "key")(nil, request)
	require.NotNil(kv)
	require.Len(kv, 1)
	assert.True(kv[0].Equals(zap.Strings("key", []string{"value1", "value2"})))
}

func TestHeader(t *testing.T) {
	t.Run("Missing", testHeaderMissing)
	t.Run("SingleValue", testHeaderSingleValue)
	t.Run("MultiValue", testHeaderMultiValue)
}

func testSetLoggerNilBase(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() {
		SetLogger(nil)
	})
}

func testSetLoggerBaseOnly(t *testing.T) {
	assert := assert.New(t)
	base := sallust.Default()
	request := httptest.NewRequest("GET", "/", nil)
	ctx := SetLogger(base)(context.Background(), request)

	assert.Equal(base, sallust.Get(ctx))
}

func testSetLoggerCustom(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	variables := map[string]string{
		"test": "path variable value",
	}
	request := mux.SetURLVars(
		httptest.NewRequest("GET", "/test/uri", nil),
		variables,
	)
	verify, base := sallust.NewTestLogger()

	request.RemoteAddr = "10.0.0.1:7777"
	request.Header.Set("X-Test", "header value")
	logger := sallust.Get(
		SetLogger(
			base,
			RequestInfo, Header("X-Test", "key1"),
		)(context.Background(), request),
	)
	require.NotEqual(base, logger)
	logger.Info("test message")
	entry := map[string]interface{}{}
	require.NoError(json.Unmarshal(verify.Bytes(), &entry))
	assert.Equal("GET", entry[requestMethodKey])
	assert.Equal("/test/uri", entry[requestURIKey])
	assert.Equal("10.0.0.1:7777", entry[remoteAddrKey])
	assert.Equal("header value", entry["key1"])
	assert.Equal("test message", entry["msg"])
}

func TestSetLogger(t *testing.T) {
	t.Run("NilBase", testSetLoggerNilBase)
	t.Run("BaseOnly", testSetLoggerBaseOnly)
	t.Run("Custom", testSetLoggerCustom)
}
