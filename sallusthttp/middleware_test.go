package sallusthttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func testMiddlewareDefaults(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		logger  = new(zap.Logger)

		m Middleware
	)

	h := m.Decorate(nil)
	decorator, ok := h.(*handler)
	require.True(ok)
	assert.Equal(zap.NewNop(), decorator.base)
	assert.Equal(http.DefaultServeMux, decorator.next)
	require.NotNil(decorator.builder)
	assert.Equal(logger, decorator.builder(httptest.NewRequest("GET", "/", nil), logger))

	h = m.DecorateFunc(nil)
	decorator, ok = h.(*handler)
	require.True(ok)
	assert.Equal(zap.NewNop(), decorator.base)
	assert.Equal(http.DefaultServeMux, decorator.next)
	require.NotNil(decorator.builder)
	assert.Equal(logger, decorator.builder(httptest.NewRequest("GET", "/", nil), logger))
}

func testMiddlewareDecorate(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)

		next http.Handler = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			logger := Get(request)
			logger.Info("message")
			logger.Sync()

			cc, ok := logger.Core().(*sallust.CaptureCore)
			require.True(ok)

			n, err := cc.EachMessage(func(e zapcore.Entry, f []zapcore.Field) error {
				assert.Equal("message", e.Message)
				assert.Equal(
					[]zapcore.Field{zap.String("foo", "bar")},
					f,
				)

				return nil
			})

			assert.Equal(1, n)
			assert.NoError(err)

			response.WriteHeader(599)
		})

		enc = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
		})

		m = Middleware{
			Base: zap.New(
				sallust.NewCaptureCore(enc, zapcore.AddSync(ioutil.Discard), zapcore.DebugLevel),
			),
		}

		response = httptest.NewRecorder()
		request  = httptest.NewRequest("PUT", "/test", nil)
	)

	m.Builders.Add(
		func(r *http.Request, l *zap.Logger) *zap.Logger {
			return l.With(zap.String("foo", "bar"))
		},
	)

	m.Decorate(next).ServeHTTP(response, request)
	assert.Equal(599, response.Code)
}

func testMiddlewareDecorateFunc(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)

		next = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			logger := Get(request)
			logger.Info("message")
			logger.Sync()

			cc, ok := logger.Core().(*sallust.CaptureCore)
			require.True(ok)

			n, err := cc.EachMessage(func(e zapcore.Entry, f []zapcore.Field) error {
				assert.Equal("message", e.Message)
				assert.Equal(
					[]zapcore.Field{zap.String("foo", "bar")},
					f,
				)

				return nil
			})

			assert.Equal(1, n)
			assert.NoError(err)

			response.WriteHeader(599)
		})

		enc = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey: "msg",
		})

		m = Middleware{
			Base: zap.New(
				sallust.NewCaptureCore(enc, zapcore.AddSync(ioutil.Discard), zapcore.DebugLevel),
			),
		}

		response = httptest.NewRecorder()
		request  = httptest.NewRequest("PUT", "/test", nil)
	)

	m.Builders.Add(
		func(r *http.Request, l *zap.Logger) *zap.Logger {
			return l.With(zap.String("foo", "bar"))
		},
	)

	m.DecorateFunc(next).ServeHTTP(response, request)
	assert.Equal(599, response.Code)
}

func TestMiddleware(t *testing.T) {
	t.Run("Defaults", testMiddlewareDefaults)
	t.Run("Decorate", testMiddlewareDecorate)
	t.Run("DecorateFunc", testMiddlewareDecorateFunc)
}
