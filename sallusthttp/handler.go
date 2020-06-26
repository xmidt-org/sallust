package sallusthttp

import (
	"net/http"

	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
)

// handler is an http.Handler that augments the request Context with
// a logger derived from a given Base logger.
type handler struct {
	next    http.Handler
	base    *zap.Logger
	builder Builder
}

// ServeHTTP creates a logger from the Base and invokes the next handler using
// a request that has that logger in the context.  Downstream HTTP handling code
// may use sallust.Get(request.Context()) to access that logger.
func (h *handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	logger := h.builder(request, h.base)
	h.ServeHTTP(
		response,
		request.WithContext(
			sallust.With(request.Context(), logger),
		),
	)
}
