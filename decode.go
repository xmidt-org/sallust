package sallust

import (
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapcoreLevelType      = reflect.TypeOf(zapcore.Level(0))
	zapcoreLevelPtrType   = reflect.PtrTo(zapcoreLevelType)
	zapAtomicLevelType    = reflect.TypeOf(zap.AtomicLevel{})
	zapAtomicLevelPtrType = reflect.PtrTo(zapAtomicLevelType)
)

// StringToLevelHookFunc is a mapstructure DecodeHookFuncType implementation
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
func StringToLevelHookFunc(from, to reflect.Type, src interface{}) (interface{}, error) {
	if text, ok := src.(string); ok {
		switch {
		case to == zapcoreLevelType:
			var l zapcore.Level
			err := l.UnmarshalText([]byte(text))
			return l, err

		case to == zapcoreLevelPtrType:
			l := new(zapcore.Level)
			err := l.UnmarshalText([]byte(text))
			return l, err

		case to == zapAtomicLevelType:
			l := zap.NewAtomicLevel()
			err := l.UnmarshalText([]byte(text))
			return l, err

		case to == zapAtomicLevelPtrType:
			l := zap.NewAtomicLevel()
			err := l.UnmarshalText([]byte(text))
			return &l, err
		}
	}

	return src, nil
}
