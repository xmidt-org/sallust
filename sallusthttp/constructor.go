package sallusthttp

import (
	"net/http"

	"go.uber.org/zap"
)

// Constructor is responsible for decorating http.Handler instances with
// logging information
type Constructor struct {
	// Base is the base zap.Logger from which request loggers are derived.
	// If this field is nil, a Nop logger is used instead.
	Base *zap.Logger

	// Builders is the sequence of Builder strategies used to tailor the Base logger
	Builders Builders
}

// Decorate is a middleware function for augmenting request contexts with
// loggers.  If next is nil, then this function decorates http.DefaultServeMux.
//
// This function may be used with gorilla/mux, e.g.:
//
//   var c Constructor
//   c.Builders.Add(Named("myHandler"), DefaultFields)
//   r := mux.NewRouter()
//   r.UseMiddleware(c.Decorate)
//   r.Handle("/", MyHandler)
func (c Constructor) Decorate(next http.Handler) http.Handler {
	// keep a similar behavior to justinas/alice:
	if next == nil {
		next = http.DefaultServeMux
	}

	base := c.Base
	if base == nil {
		base = zap.NewNop()
	}

	return &handler{
		next:    next,
		base:    base,
		builder: c.Builders.Build,
	}
}

// DecorateFunc is syntactic sugar for decorating an HTTP handler function.
// If the given function is nil, this function decorates http.DefaultServeMux.
func (c Constructor) DecorateFunc(next http.HandlerFunc) http.Handler {
	if next == nil {
		// ensure a true nil gets passed
		return c.Decorate(nil)
	}

	return c.Decorate(next)
}
