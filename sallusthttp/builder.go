// SPDX-FileCopyrightText: 2024 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package sallusthttp

import (
	"net/http"

	"go.uber.org/zap"
)

const (
	// DefaultMethodKey is the default logging key for an HTTP request method
	DefaultMethodKey = "method"

	// DefaultURIKey is the default logging key for an HTTP request URI
	DefaultURIKey = "uri"

	// DefaultRemoteAddrKey is the default logging key for a request's remote address
	DefaultRemoteAddrKey = "remoteAddr"
)

// Builder is a strategy for augmenting a zap.Logger for an HTTP request.
// A Builder may copy information from the request, or it my inject information
// from other sources.
type Builder func(*http.Request, *zap.Logger) *zap.Logger

// FieldBuilder is a strategy for the special case of building sequences
// of zap logging fields for an HTTP request.  Using this strategy to build
// up multiple fields at once can be more efficient than multiple calls to zap.Logger.With
// required by the Builder strategy.
type FieldBuilder func(*http.Request, []zap.Field) []zap.Field

// Fields adapts a sequence of FieldBuilder strategies into a single Builder.
func Fields(fb ...FieldBuilder) Builder {
	return func(r *http.Request, l *zap.Logger) *zap.Logger {
		var fields []zap.Field
		for _, f := range fb {
			fields = f(r, fields)
		}

		// NOTE: zap's implementation efficiently handles the case when the
		// number of fields is zero
		return l.With(fields...)
	}
}

// Builders represents an aggregation of individual Builder strategies.
// This slice is mutable, and as such any of the methods that modify the
// builders slice are not safe for concurrent access.
type Builders []Builder

// Build is itself a Builder that simply applies a sequence of strategies in order.
// This method is safe for concurrent access.
func (b Builders) Build(r *http.Request, l *zap.Logger) *zap.Logger {
	for _, f := range b {
		l = f(r, l)
	}

	return l
}

// Add appends more Builder strategies to this Builders sequence
func (b *Builders) Add(more ...Builder) {
	if len(more) > 0 {
		*b = append(*b, more...)
	}
}

// AddFields adds a number of fields under a single Builder
func (b *Builders) AddFields(more ...FieldBuilder) {
	if len(more) > 0 {
		*b = append(*b, Fields(more...))
	}
}

// Named returns a Builder that creates a sublogger with the given name.  Useful
// to separate functional HTTP areas, such as different handlers, servers, or packages.
func Named(name string) Builder {
	return func(_ *http.Request, l *zap.Logger) *zap.Logger {
		return l.Named(name)
	}
}

// DefaultFields is a Builder that simply appends the basic request fields implemented
// in this package under their default logging keys
func DefaultFields(r *http.Request, l *zap.Logger) *zap.Logger {
	return l.With(
		zap.String(DefaultMethodKey, r.Method),
		zap.String(DefaultURIKey, r.RequestURI),
		zap.String(DefaultRemoteAddrKey, r.RemoteAddr),
	)
}

// Method is a FieldBuilder that adds the request method under the DefaultMethodKey
func Method(r *http.Request, f []zap.Field) []zap.Field {
	return append(f, zap.String(DefaultMethodKey, r.Method))
}

// MethodCustom creates a FieldBuilder that adds the request methd under a custom key
func MethodCustom(key string) FieldBuilder {
	return func(r *http.Request, f []zap.Field) []zap.Field {
		return append(f, zap.String(key, r.Method))
	}
}

// URI is a FieldBuilder that adds the request URI under the DefaultURIKey
func URI(r *http.Request, f []zap.Field) []zap.Field {
	return append(f, zap.String(DefaultURIKey, r.RequestURI))
}

// URICustom creates a FieldBuilder that adds the request URI under a custom key
func URICustom(key string) FieldBuilder {
	return func(r *http.Request, f []zap.Field) []zap.Field {
		return append(f, zap.String(key, r.RequestURI))
	}
}

// RemoteAddr is a FieldBuilder that adds the request's remote address under the DefaultRemoteAddrKey
func RemoteAddr(r *http.Request, f []zap.Field) []zap.Field {
	return append(f, zap.String(DefaultRemoteAddrKey, r.RemoteAddr))
}

// RemoteAddrCustom creates a FieldBuilder that adds the remote address under a custom key
func RemoteAddrCustom(key string) FieldBuilder {
	return func(r *http.Request, f []zap.Field) []zap.Field {
		return append(f, zap.String(key, r.RemoteAddr))
	}
}
