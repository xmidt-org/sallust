// SPDX-FileCopyrightText: 2024 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package sallusthttp

import (
	"net/http"

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
	h.next.ServeHTTP(response, With(request, logger))
}
