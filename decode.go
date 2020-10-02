package sallust

import (
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	stringType = reflect.TypeOf("")

	levelType          = reflect.TypeOf(zapcore.Level(0))
	levelPtrType       = reflect.PtrTo(levelType)
	atomicLevelType    = reflect.TypeOf(zap.AtomicLevel{})
	atomicLevelPtrType = reflect.PtrTo(atomicLevelType)

	levelEncoderType    = reflect.TypeOf(zapcore.LevelEncoder(nil))
	timeEncoderType     = reflect.TypeOf(zapcore.TimeEncoder(nil))
	durationEncoderType = reflect.TypeOf(zapcore.DurationEncoder(nil))
	callerEncoderType   = reflect.TypeOf(zapcore.CallerEncoder(nil))
	nameEncoderType     = reflect.TypeOf(zapcore.NameEncoder(nil))
)

// DecodeHook is an all-in-one mapstructure DecodeHookFunc that converts from
// a string (typically unmarshaled in something like spf13/viper) into the appropriate
// configuration field required by zapcore.
//
// The from type must refer exactly to a string, not to a type derived from string.
//
// The to type may be one of:
//
//   zapcore.Level
//   *zapcore.Level
//   zap.AtomicLevel
//   *zap.AtomicLevel
//   zapcore.LevelEncoder
//   zapcore.TimeEncoder
//   zapcore.DurationEncoder
//   zapcore.CallerEncoder
//   zapcore.NameEncoder
//
// The UnmarshalText method of the to type is used to do the conversion.
//
// Any other from or to type will cause the function to do no conversion and
// return the src as is with no error.
func DecodeHook(from, to reflect.Type, src interface{}) (interface{}, error) {
	if from != stringType {
		return src, nil
	}

	switch to {
	case levelType:
		var l zapcore.Level
		err := l.UnmarshalText([]byte(src.(string)))
		return l, err

	case levelPtrType:
		l := new(zapcore.Level)
		err := l.UnmarshalText([]byte(src.(string)))
		return l, err

	case atomicLevelType:
		l := zap.NewAtomicLevel()
		err := l.UnmarshalText([]byte(src.(string)))
		return l, err

	case atomicLevelPtrType:
		l := zap.NewAtomicLevel()
		err := l.UnmarshalText([]byte(src.(string)))
		return &l, err

	case levelEncoderType:
		var le zapcore.LevelEncoder
		err := le.UnmarshalText([]byte(src.(string)))
		return le, err

	case timeEncoderType:
		var te zapcore.TimeEncoder
		err := te.UnmarshalText([]byte(src.(string)))
		return te, err

	case durationEncoderType:
		var de zapcore.DurationEncoder
		err := de.UnmarshalText([]byte(src.(string)))
		return de, err

	case callerEncoderType:
		var ce zapcore.CallerEncoder
		err := ce.UnmarshalText([]byte(src.(string)))
		return ce, err

	case nameEncoderType:
		var ne zapcore.NameEncoder
		err := ne.UnmarshalText([]byte(src.(string)))
		return ne, err

	default:
		return src, nil
	}
}
