package sallusthttp

import (
	"net/http"

	"go.uber.org/zap"
)

const (
	// MethodKey is the logging key used by DefaultFields for the HTTP Method
	MethodKey = "method"

	// URIKey is the logging key used by DefaultFields for the RequestURI
	URIKey = "uri"

	// RemoteAddrKey is the logging key used by DefaultFields for the RemoteAddr
	RemoteAddrKey = "remoteAddr"
)

// Builder is a strategy for augmenting a zap.Logger for an HTTP request.
// A Builder may copy information from the request, or it my inject information
// from other sources.
type Builder func(*http.Request, *zap.Logger) *zap.Logger

// Builders represents an aggregation of individual Builder strategies
type Builders []Builder

// Build is itself a Builder that simply apply a sequence of strategies in order
func (b Builders) Build(r *http.Request, l *zap.Logger) *zap.Logger {
	for _, f := range b {
		l = f(r, l)
	}

	return l
}

// Append adds more Builder strategies to this Builders sequence.  This method
// is valid to use on an uninitialized Builders slice.
func (b *Builders) Append(more ...Builder) {
	*b = append(*b, more...)
}

// Named returns a Builder that creates a sublogger with the given name.  Useful
// to separate functional HTTP areas, such as different handlers, servers, or packages.
func Named(name string) Builder {
	return func(_ *http.Request, l *zap.Logger) *zap.Logger {
		return l.Named(name)
	}
}

// DefaultFields is a Builder that returns a logger with method, URI, and content length
// fields added.
func DefaultFields(r *http.Request, l *zap.Logger) *zap.Logger {
	return l.With(
		zap.String(MethodKey, r.Method),
		zap.String(URIKey, r.RequestURI),
		zap.String(RemoteAddrKey, r.RemoteAddr),
	)
}
