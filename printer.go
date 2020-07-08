package sallust

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Printer adapts the zap.SugaredLogger type onto a style of logging more like the fmt package.
// The stdlib's log package uses this approach.  Also, go.uber.org/fx can use this type
// as a logger for DI container information:
//
//   logger := zap.NewDevelopment()
//   fx.New(
//     fx.Logger(sallust.Printer{SugaredLogger: logger.Sugar()}),
//     // etc
//   )
type Printer struct {
	// SugaredLogger is the zap.SugaredLogger being adapted.  A Printer still allows all the
	// usual logging methods to be called.
	*zap.SugaredLogger

	// Level is the log level used to emit print messages.
	// If unset, zapcore.InfoLevel is used.
	//
	// The panic and fatal levels in the zapcore are not supported.
	// If this field is set to one of those levels, InfoLevel is used instead.
	Level zapcore.Level
}

// Printf invokes the appropriate "*f" method on the SugaredLogger, based
// on the Level.  The go.uber.org/fx package can use this method as an fx.Logger.
// The stdlib log package also exposes a method with this same signature.
func (p Printer) Printf(format string, args ...interface{}) {
	switch p.Level {
	case zapcore.DebugLevel:
		p.SugaredLogger.Debugf(format, args...)

	case zapcore.WarnLevel:
		p.SugaredLogger.Warnf(format, args...)

	case zapcore.ErrorLevel:
		p.SugaredLogger.Errorf(format, args...)

	case zapcore.InfoLevel:
		fallthrough

	default:
		p.SugaredLogger.Infof(format, args...)
	}
}

// Print invokes the appropriate leveled fmt.Print-style method on the SugaredLogger,
// based on the Level.
func (p Printer) Print(args ...interface{}) {
	switch p.Level {
	case zapcore.DebugLevel:
		p.SugaredLogger.Debug(args...)

	case zapcore.WarnLevel:
		p.SugaredLogger.Warn(args...)

	case zapcore.ErrorLevel:
		p.SugaredLogger.Error(args...)

	case zapcore.InfoLevel:
		fallthrough

	default:
		p.SugaredLogger.Info(args...)
	}
}
