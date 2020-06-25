package sallust

//go:generate stringer -type=Encoding

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

// Encoding is an enumerated type that specifies the supported encoding strategies
type Encoding int

const (
	// EncodingJSON refers to the built-in zapcore JSON Encoder.
	// This is the default encoding.
	//
	// See: https://pkg.go.dev/go.uber.org/zap@v1.15.0/zapcore?tab=doc#NewJSONEncoder
	EncodingJSON Encoding = iota

	// EncodingConsole refers to the built-in zapcore console Encoder,
	// which is like the JSON Encoder except with some human-readable tweaks.
	//
	// This does NOT refer to stdout; rather, it refers to an encoding appropriate
	// for a console.
	//
	// See: https://pkg.go.dev/go.uber.org/zap@v1.15.0/zapcore?tab=doc#NewConsoleEncoder
	EncodingConsole
)

// UnmarshalText implements encoding.TextUnmarshaler
func (e *Encoding) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	default:
		*e = EncodingJSON
	case "console":
		*e = EncodingConsole
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e Encoding) MarshalText() ([]byte, error) {
	switch e {
	default:
		return []byte("json"), nil

	case EncodingConsole:
		return []byte("console"), nil
	}
}

// NewEncoder creates the encoder appropriate for this Encoding value.
func (e Encoding) NewEncoder(ec zapcore.EncoderConfig) zapcore.Encoder {
	ec = ApplyEncoderConfigDefaults(ec)
	switch e {
	default:
		return zapcore.NewJSONEncoder(ec)
	case EncodingConsole:
		return zapcore.NewConsoleEncoder(ec)
	}
}

const (
	// DefaultMessageKey is the default value for zapcore.EncoderConfig.MessageKey
	DefaultMessageKey = "msg"

	// DefaultLevelKey is the default value for zapcore.EncoderConfig.LevelKey
	DefaultLevelKey = "level"

	// DefaultTimeKey is the default value for zapcore.EncoderConfig.TimeKey
	DefaultTimeKey = "ts"

	// DefaultNameKey is the default value for zapcore.EncoderConfig.NameKey
	DefaultNameKey = "logger"

	// DefaultCallerKey is the default value for zapcore.EncoderConfig.CallerKey
	DefaultCallerKey = "caller"

	// DefaultStacktraceKey is the default value for zapcore.EncoderConfig.StacktraceKey
	DefaultStacktraceKey = "stacktrace"
)

// ApplyEncoderConfigDefaults produces a clone of a given zapcore.EncoderConfig
// with any missing fields set to the defaults used by sallust.
func ApplyEncoderConfigDefaults(original zapcore.EncoderConfig) zapcore.EncoderConfig {
	clone := original
	if len(clone.MessageKey) == 0 {
		clone.MessageKey = DefaultMessageKey
	}

	if len(clone.LevelKey) == 0 {
		clone.LevelKey = DefaultLevelKey
	}

	if len(clone.TimeKey) == 0 {
		clone.TimeKey = DefaultTimeKey
	}

	if len(clone.NameKey) == 0 {
		clone.NameKey = DefaultNameKey
	}

	if len(clone.CallerKey) == 0 {
		clone.CallerKey = DefaultCallerKey
	}

	if len(clone.StacktraceKey) == 0 {
		clone.StacktraceKey = DefaultStacktraceKey
	}

	if len(clone.LineEnding) == 0 {
		clone.LineEnding = zapcore.DefaultLineEnding
	}

	if clone.EncodeLevel == nil {
		clone.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	if clone.EncodeTime == nil {
		clone.EncodeTime = zapcore.RFC3339TimeEncoder
	}

	if clone.EncodeDuration == nil {
		clone.EncodeDuration = zapcore.StringDurationEncoder
	}

	if clone.EncodeCaller == nil {
		clone.EncodeCaller = zapcore.ShortCallerEncoder
	}

	return clone
}
