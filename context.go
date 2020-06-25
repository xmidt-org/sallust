package sallust

import (
	"context"

	"go.uber.org/zap"
)

// contextKey is the internal key type used to access a zap.Logger
// within a context.Context instance
type contextKey struct{}

// With places a zap.Logger into the context.  If the given logger is nil,
// this function returns the parent as-is.  Since the Get functions return
// a nop logger when there is no logger in the context, a nil logger
// can be safely ignored.
//
// See: https://pkg.go.dev/go.uber.org/zap?tab=doc#Logger
func With(parent context.Context, logger *zap.Logger) context.Context {
	if logger != nil {
		return context.WithValue(parent, contextKey{}, logger)
	}

	return parent
}

// Get returns the zap.Logger from the given context.  If no zap.Logger
// exists, this function returns a nop logger.
//
// See: https://pkg.go.dev/go.uber.org/zap?tab=doc#Logger
// See: https://pkg.go.dev/go.uber.org/zap?tab=doc#NewNop
func Get(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(contextKey{}).(*zap.Logger); ok {
		return l
	}

	return zap.NewNop()
}

// GetDefault attempts to find a zap.Logger in the given context.  If none is
// found, the given default is returned.  If the given default is nil, then
// a nop logger is returned instead.
//
// See: https://pkg.go.dev/go.uber.org/zap?tab=doc#Logger
func GetDefault(ctx context.Context, def *zap.Logger) *zap.Logger {
	if l, ok := ctx.Value(contextKey{}).(*zap.Logger); ok {
		return l
	}

	if def != nil {
		return def
	}

	return zap.NewNop()
}
