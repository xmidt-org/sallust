// SPDX-FileCopyrightText: 2024 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package sallusthttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandler(t *testing.T) {
	var (
		assert     = assert.New(t)
		base       = zap.NewNop()
		contextual = base.With(zap.String("foo", "bar"))

		next http.Handler = http.HandlerFunc(
			func(response http.ResponseWriter, request *http.Request) {
				assert.Equal(contextual, Get(request))
				response.WriteHeader(599)
			},
		)

		h = handler{
			next: next,
			base: base,
			builder: func(r *http.Request, l *zap.Logger) *zap.Logger {
				return contextual
			},
		}

		request  = httptest.NewRequest("PUT", "/test", nil)
		response = httptest.NewRecorder()
	)

	h.ServeHTTP(response, request)
	assert.Equal(599, response.Code) // verify that the handler was called
}
