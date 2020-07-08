package sallusthttp

import (
	"net/http"

	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
)

// With associates a zap.Logger with the given request
func With(parent *http.Request, logger *zap.Logger) *http.Request {
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
