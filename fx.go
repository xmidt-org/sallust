package sallust

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// LoggerIn describes the dependencies used to bootstrap a zap logger within
// an fx application.
type LoggerIn struct {
	fx.In

	// Config is the sallust configuration for the logger.  This component is optional,
	// and if not supplied a default zap logger will be created.
	Config Config `optional:"true"`

	// Options are the optional zap options, injected from the enclosing fx.App.
	// If supplied, these will be appended to the options supplied directly via WithLogger.
	Options []zap.Option `optional:"true"`
}

// WithLogger bootstraps a go.uber.org/zap logger together with an fxevent.Logger,
// using the dependencies described in LoggerIn.
//
// If any zap.Options are supplied to this function, they take precedence over any
// options injected via LoggerIn.
func WithLogger(options ...zap.Option) fx.Option {
	return fx.Options(
		fx.Provide(
			func(in LoggerIn) (*zap.Logger, error) {
				merged := make([]zap.Option, 0, len(options)+len(in.Options))

				// options passed to this function take preceence over options
				// injected from the fx.App.
				merged = append(merged, options...)
				merged = append(merged, in.Options...)

				return in.Config.Build(merged...)
			},
		),
		fx.WithLogger(
			func(l *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{
					Logger: l,
				}
			},
		),
	)
}

// SyncOnShutdown adds an fx lifecycle hook that invokes Sync on the application's logger.
// Generally, this option should be placed as an fx.Invoke last in the set of options.
// That ensures that log entries from other lifecycle OnStop hooks are written to log sinks.
//
//   fx.New(
//     // all other options come first ...
//
//     sallust.SyncOnShutdown(),
//   )
func SyncOnShutdown() fx.Option {
	return fx.Invoke(
		func(logger *zap.Logger, lifecycle fx.Lifecycle) {
			lifecycle.Append(fx.Hook{
				OnStop: func(context.Context) error {
					logger.Sync()

					// NOTE: do NOT return the error from Sync.
					// A non-nil error may short-circuit app shutdown,
					// and logger errors during shutdown are never fatal.
					return nil
				},
			})
		},
	)
}
