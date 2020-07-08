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
