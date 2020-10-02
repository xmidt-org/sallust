package sallust

import (
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var stringType = reflect.TypeOf("")

var (
	zapcoreLevelType      = reflect.TypeOf(zapcore.Level(0))
	zapcoreLevelPtrType   = reflect.PtrTo(zapcoreLevelType)
	zapAtomicLevelType    = reflect.TypeOf(zap.AtomicLevel{})
	zapAtomicLevelPtrType = reflect.PtrTo(zapAtomicLevelType)
)

// StringToLevelHook is a mapstructure DecodeHookFuncType implementation
// that handles converting from a string to a zap logging level.  This function will
// convert to any of the following types, based on the to parameter:
//
//   zapcore.Level
//   *zapcore.Level
//   zap.AtomicLevel
//   *zap.AtomicLevel
//
// To use this with spf13/viper, create a viper.DecoderConfigOption and apply it:
//
//   v := viper.New()
//   v.Unmarshal(&myType, func(cfg *mapstructure.DecoderConfig) {
//     if cfg.DecodeHook != nil {
//       cfg.DecodeHook = mapstructure.ComposeDecodeHookFunc(
//         cfg.DecodeHook,
//         StringToLevelHookFunc,
//       )
//     } else {
//       cfg.DecodeHook = StringToLevelHookFunc
//     }
//   })
//
// The UnmarshalText implementation of both level types is used to perform the conversion.
//
// This function is necessary because libraries like spf13/viper do not directly unmarshal
// into types.  Rather, they unmarshal first to a neutral format, e.g. map[string]interface{},
// and then use mapstructure to take that format and convert it into a struct.
func StringToLevelHook(from, to reflect.Type, src interface{}) (interface{}, error) {
	if from == stringType {
		switch to {
		case zapcoreLevelType:
			var l zapcore.Level
			err := l.UnmarshalText([]byte(src.(string)))
			return l, err

		case zapcoreLevelPtrType:
			l := new(zapcore.Level)
			err := l.UnmarshalText([]byte(src.(string)))
			return l, err

		case zapAtomicLevelType:
			l := zap.NewAtomicLevel()
			err := l.UnmarshalText([]byte(src.(string)))
			return l, err

		case zapAtomicLevelPtrType:
			l := zap.NewAtomicLevel()
			err := l.UnmarshalText([]byte(src.(string)))
			return &l, err
		}
	}

	return src, nil
}

var levelEncoderType = reflect.TypeOf(zapcore.LevelEncoder(nil))

// StringToLevelEncoderHook is a mapstructure DecodeHookFuncType that converts from a string
// to a zapcore.LevelEncoder via the zapcore.LevelEncoder.UnmarshalText method.
func StringToLevelEncoderHook(from, to reflect.Type, src interface{}) (interface{}, error) {
	if from == stringType && to == levelEncoderType {
		var le zapcore.LevelEncoder
		err := le.UnmarshalText([]byte(src.(string)))
		return le, err
	}

	return src, nil
}

var timeEncoderType = reflect.TypeOf(zapcore.TimeEncoder(nil))

// StringToTimeEncoderHook is a mapstructure DecodeHookFuncType that converts from a string
// to a zapcore.TimeEncoder via the zapcore.TimeEncoder.UnmarshalText method.
func StringToTimeEncoderHook(from, to reflect.Type, src interface{}) (interface{}, error) {
	if from == stringType && to == timeEncoderType {
		var te zapcore.TimeEncoder
		err := te.UnmarshalText([]byte(src.(string)))
		return te, err
	}

	return src, nil
}

var durationEncoderType = reflect.TypeOf(zapcore.DurationEncoder(nil))

// StringToDurationEncoderHook is a mapstructure DecodeHookFuncType that converts from a string
// to a zapcore.DurationEncoder via the zapcore.DurationEncoder.UnmarshalText method.
func StringToDurationEncoderHook(from, to reflect.Type, src interface{}) (interface{}, error) {
	if from == stringType && to == durationEncoderType {
		var de zapcore.DurationEncoder
		err := de.UnmarshalText([]byte(src.(string)))
		return de, err
	}

	return src, nil
}

var callerEncoderType = reflect.TypeOf(zapcore.CallerEncoder(nil))

// StringToCallerEncoderHook is a mapstructure DecodeHookFuncType that converts from a string
// to a zapcore.CallerEncoder via the zapcore.CallerEncoder.UnmarshalText method.
func StringToCallerEncoderHook(from, to reflect.Type, src interface{}) (interface{}, error) {
	if from == stringType && to == callerEncoderType {
		var ce zapcore.CallerEncoder
		err := ce.UnmarshalText([]byte(src.(string)))
		return ce, err
	}

	return src, nil
}

var nameEncoderType = reflect.TypeOf(zapcore.NameEncoder(nil))

// StringToNameEncoderHook is a mapstructure DecodeHookFuncType that converts from a string
// to a zapcore.NameEncoder via the zapcore.NameEncoder.UnmarshalText method.
func StringToNameEncoderHook(from, to reflect.Type, src interface{}) (interface{}, error) {
	if from == stringType && to == nameEncoderType {
		var ne zapcore.NameEncoder
		err := ne.UnmarshalText([]byte(src.(string)))
		return ne, err
	}

	return src, nil
}
