package sallustkit

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.uber.org/zap"
)

const (
	// NotAString is used in log output when a key or value was expected
	// to be a string but was not.
	NotAString = "(KEY IS NOT A STRING)"

	// NoLogMessage is used as the log message when no message could be
	// found in the go-kit keyvals.
	NoLogMessage = "No log message found"

	// DefaultMessageKey is the key assumed to hold the log message when no
	// MessageKey is set.
	DefaultMessageKey = "msg"
)

// toString attempts to cast v to a string, returning notAString
// if it isn't a string.
func toString(v interface{}) string {
	s, ok := v.(string)
	if ok {
		return s
	}

	return NotAString
}

// Logger is a go-kit logger that adapts its output onto a zap logger.
type Logger struct {
	// Zap is the zap Logger to which output is sent.  This field is required,
	// and Log will panic if it is not set.
	Zap *zap.Logger

	// MessageKey is the go-kit logging key which holds the log message.
	// If unset, DefaultMessageKey is used.
	//
	// This key is used to pull out the message so that it can be passed
	// as the first argument to a zap logger's method, e.g. Error, Info, etc.
	MessageKey string

	// DefaultLevel is the go-kit level to use when no level
	// is supplied in the keyvals.  If unset, Error is used.
	DefaultLevel level.Value
}

var _ log.Logger = Logger{}

// Log accepts key/value pairs in the typical go-kit fashion and parses them
// to use with the configured zap logger.  This method always returns nil.
//
// Each key/value pair is examined and used to build up a method call to
// the configured zap logger using the following basic steps:
//
//   - Any key that is not a string results in a NotAString name/value pair in the zap output
//   - The value for any key that equals the configured MessageKey (or, DefaultMessageKey
//     if that field is unset) is used as the first parameter to the zap logger method.
//   - Any value that is a defined go-kit level.Value is used to determine which zap
//     logger method is invoked, e.g. level.DebugValue() results in the Debug method, etc.
//   - Any key/value not matching the above steps is passed to the zap logger method
//     as a zap.Any field.
//
// Examples:
//
//   given:
//
//   l := zap.NewDevelopment() // or any *zap.Logger
//   gokit := sallustkit.Logger{
//     Zap: l,
//     // take the defaults for the other fields
//   }
//
//   then:
//
//   this:     gokit.Log("msg", "hi there", "value", 123)
//   becomes:  l.Error("hi there", zap.Any("value", 123)) // no level, can change this by setting gokit.DefaultLevel
//
//   this:     gokit.Log("msg", "more values", "name1", "value1", "name2", 45.6)
//   becomes:  l.Error("more values", zap.Any("name1", "value1"), zap.Any("name2", 45.6))
//
//   this:     gokit.Log(level.Key(), level.InfoValue(), "value", 123)
//   becomes:  l.Info("No log message found", zap.Any("value", 123))
//
//   this:     gokit.Log("msg", "hi there", "this key doesn't matter", level.DebugValue())
//   becomes:  l.Debug("hi there") // if a value is a gokit level.Value, the key is ignored
//
func (l Logger) Log(keyvals ...interface{}) error {
	if len(keyvals) > 0 {
		var (
			message    = NoLogMessage
			messageKey = DefaultMessageKey
			lvl        = l.DefaultLevel
			fields     = make([]zap.Field, 0, len(keyvals)/2)
		)

		if len(l.MessageKey) > 0 {
			messageKey = l.MessageKey
		}

		for i, j := 0, 1; j < len(keyvals); i, j = i+2, j+2 {
			key := toString(keyvals[i])
			value := keyvals[j]

			if key == messageKey {
				message = toString(value)
				continue
			}

			if fieldLevel, ok := value.(level.Value); ok {
				lvl = fieldLevel
				continue
			}

			fields = append(fields, zap.Any(key, value))
		}

		if len(keyvals)%2 != 0 {
			// odd number of keyvals ...
			fields = append(fields, zap.NamedError(
				toString(keyvals[len(keyvals)-1]),
				log.ErrMissingValue,
			))
		}

		LogMethodFor(l.Zap, lvl)(message, fields...)
	}

	return nil
}