package sallust

import (
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
	Config  Config       `optional:"true"`
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
