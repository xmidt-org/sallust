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

func decodeLevel(text string) (l zapcore.Level, err error) {
	err = l.UnmarshalText([]byte(text))
	return
}

func decodeLevelPointer(text string) (l *zapcore.Level, err error) {
	l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(text))
	return
}

func decodeAtomicLevel(text string) (al zap.AtomicLevel, err error) {
	al = zap.NewAtomicLevel()
	err = al.UnmarshalText([]byte(text))
	return
}

func decodeAtomicLevelPointer(text string) (p *zap.AtomicLevel, err error) {
	al := zap.NewAtomicLevel()
	err = al.UnmarshalText([]byte(text))
	p = &al
	return
}

func decodeLevelEncoder(text string) (le zapcore.LevelEncoder, err error) {
	err = le.UnmarshalText([]byte(text))
	return
}

func decodeTimeEncoder(text string) (te zapcore.TimeEncoder, err error) {
	err = te.UnmarshalText([]byte(text))
	return
}

func decodeDurationEncoder(text string) (de zapcore.DurationEncoder, err error) {
	err = de.UnmarshalText([]byte(text))
	return
}

func decodeCallerEncoder(text string) (ce zapcore.CallerEncoder, err error) {
	err = ce.UnmarshalText([]byte(text))
	return
}

func decodeNameEncoder(text string) (ne zapcore.NameEncoder, err error) {
	err = ne.UnmarshalText([]byte(text))
	return
}

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

	text := src.(string)

	switch to {
	case levelType:
		return decodeLevel(text)

	case levelPtrType:
		return decodeLevelPointer(text)

	case atomicLevelType:
		return decodeAtomicLevel(text)

	case atomicLevelPtrType:
		return decodeAtomicLevelPointer(text)

	case levelEncoderType:
		return decodeLevelEncoder(text)

	case timeEncoderType:
		return decodeTimeEncoder(text)

	case durationEncoderType:
		return decodeDurationEncoder(text)

	case callerEncoderType:
		return decodeCallerEncoder(text)

	case nameEncoderType:
		return decodeNameEncoder(text)

	default:
		return src, nil
	}
}
