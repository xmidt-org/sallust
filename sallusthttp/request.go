package sallusthttp

import (
	"context"
	"net/http"
	"net/textproto"

	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
)

const (
	requestMethodKey = "requestMethod"
	requestURIKey    = "requestURI"
	remoteAddrKey    = "remoteAddr"
)

// With associates a zap.Logger with the given request
func With(parent *http.Request, logger *zap.Logger, b ...Builder) *http.Request {
	for _, f := range b {
		logger = f(parent, logger)
	}

	return parent.WithContext(
		sallust.With(parent.Context(), logger),
	)
}

// Get returns the zap.Logger associated with the given HTTP request
func Get(r *http.Request) *zap.Logger {
	return sallust.Get(r.Context())
}

// GetDefault returns the zap.Logger associated with the request,
// falling back to the given default if no such logger is found
func GetDefault(r *http.Request, def *zap.Logger) *zap.Logger {
	return sallust.GetDefault(r.Context(), def)
}

// SetLogger produces a zap RequestFunc that inserts a zap Logger into the context.
// Zero or more LoggerFuncs can be provided to added key/values.  Note that nothing is added to
// the base logger by default.  If no LoggerFuncs are supplied, the base Logger is added to the
// context as is.  In particular, RequestInfo must be used to inject the request method, uri, etc.
//
// The base logger must be non-nil.  There is no default applied.
//
// The returned function can be used with xcontext.Populate.
func SetLogger(l *zap.Logger, lf ...sallust.LoggerFunc) func(context.Context, *http.Request) context.Context {
	if l == nil {
		panic("The base Logger cannot be nil")
	}

	if len(lf) > 0 {
		return func(ctx context.Context, request *http.Request) context.Context {
			kv := []zap.Field{}
			for _, f := range lf {
				kv = f(kv, request)
			}

			return sallust.With(
				ctx,
				l.With(kv...),
			)
		}
	}

	return func(ctx context.Context, _ *http.Request) context.Context {
		return sallust.With(ctx, l)
	}
}

// RequestInfo is a LoggerFunc that adds the request information described by logging keys in this package.
func RequestInfo(kv []zap.Field, request *http.Request) []zap.Field {
	return append(kv,
		zap.String(requestMethodKey, request.Method),
		zap.String(requestURIKey, request.RequestURI),
		zap.String(remoteAddrKey, request.RemoteAddr),
	)

}

// Header returns a logger func that extracts the value of a header and inserts it as the
// value of a logging key.  If the header is not present in the request, a blank string
// is set as the logging key's value.
func Header(headerName, keyName string) sallust.LoggerFunc {
	headerName = textproto.CanonicalMIMEHeaderKey(headerName)

	return func(kv []zap.Field, request *http.Request) []zap.Field {
		values := request.Header[headerName]
		switch len(values) {
		case 0:
			return append(kv, zap.String(keyName, ""))
		case 1:
			return append(kv, zap.String(keyName, values[0]))
		default:
			return append(kv, zap.Strings(keyName, values))
		}
	}
}
