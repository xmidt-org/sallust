package sallustkit

import (
	"github.com/go-kit/log/level"
	"go.uber.org/zap"
)

// LogMethod refers to a method of a zap Logger.  Used to define
// which method should received go-kit levelled logging, e.g. Error, Info, etc.
type LogMethod func(string, ...zap.Field)

// LogMethodFor returns the method of a zap Logger that corresponds to
// a given go-kit level.  If v is unrecognized, l.Error is returned.
// This function never returns nil.
func LogMethodFor(l *zap.Logger, v level.Value) (lm LogMethod) {
	switch v {
	case level.DebugValue():
		lm = l.Debug

	case level.InfoValue():
		lm = l.Info

	case level.WarnValue():
		lm = l.Warn

	default:
		lm = l.Error
	}

	return
}
