package sallust

import (
	"os"
	"sort"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	stdout = zapcore.Lock(os.Stdout)
	stderr = zapcore.Lock(os.Stderr)
)

const (
	// Stdout is the special location indicating stdout.  Additionally,
	// if a sink location is the empty string, stdout is assumed.
	Stdout = "stdout"

	// Stderr is the special location indicating stderr
	Stderr = "stderr"
)

// Sink describes a destination for logging.
// Lumberjack is used to do the actual rotation.
//
// See: https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2?tab=doc#Logger
type Sink struct {
	// DisableRotation indicates that the log file in question should not be rotated.
	// This option is implied for any terminal-based logging, such as stdout.
	//
	// If this field is true, all the rotation fields of this struct are ignored.
	DisableRotation bool `json:"disableRotation" yaml:"disableRotation"`

	// MaxSize is the maximum size of the log file before rotation, in megabytes.
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files
	MaxAge int `json:"maxage" yaml:"maxage"`

	// MaxBackups is the maximum number of backups to retain
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// LocalTime indicates that timestamps in log file names should be in local time.
	// The default is false, indicating that log file names have UTC timestamps.
	//
	// NOTE: This field does NOT control timestamps in log messages.  It is used to
	// generate backup log file names when log files are rotated.
	LocalTime bool `json:"localtime" yaml:"localtime"`

	// Compress indicates that log files should be compressed using gzip.
	// The default is false, indicating no compression.
	Compress bool `json:"compress" yaml:"compress"`

	// Level is the threshold for logging.  If unset, the default is zapcore.ErrorLevel.
	//
	// NOTE: Unlike zap, sallust does not return an error if this field is unset.
	Level zap.AtomicLevel `json:"level" yaml:"level"`

	// Encoding specifies the zap.Encoder for this sink's log output.  The default
	// is to use JSON.
	//
	// See: https://pkg.go.dev/go.uber.org/zap@v1.15.0/zapcore?tab=doc#NewConsoleEncoder
	// See: https://pkg.go.dev/go.uber.org/zap@v1.15.0/zapcore?tab=doc#NewJSONEncoder
	Encoding Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`

	// EncoderConfig is the zapcore.EncoderConfig used to configure the Encoder.  Unlike zap itself,
	// sallust does not allow fields to be disabled.  There are default values for all fields, used when
	// the corresponding field is not set.
	//
	// See: https://pkg.go.dev/go.uber.org/zap@v1.15.0/zapcore?tab=doc#EncoderConfig
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
}

// Sinks is a mapping between locations, such as file system paths, and Sink configurations
type Sinks map[string]Sink

// Get retrieves the Sink associated with the given location.
// If no such location is registered, the second return value will be false.
func (s Sinks) Get(location string) (Sink, bool) {
	v, ok := s[location]
	return v, ok
}

// Set establishes the Sink configuration for a location.  Any existing
// Sink for that location will be overwritten.
func (s *Sinks) Set(location string, sink Sink) {
	if *s == nil {
		*s = make(Sinks)
	}

	(*s)[location] = sink
}

// Options describes the set of options for building a single zap.Logger
type Options struct {
	// Development sets the development flag of the zap.Logger, which controls the
	// behavior of DPanic.
	//
	// See: https://pkg.go.dev/go.uber.org/zap?tab=doc#Development
	Development bool `json:"development" yaml:"development"`

	// DisableCaller disables the caller output in log messages.
	// This field corresponds to zap.Config.DisableCaller.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`

	// DisableStacktrace disables stacktrace output.
	// This field corresponds to zap.Config.DisableStacktrace.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`

	// InitialFields specifies the logging fields for the root logger created using these
	// Options.  This field is optional, and there is no default.
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`

	// Sinks is the mapping of locations to which log output is sent.  If
	// this sinks is empty, then a nop logger is assumed.
	Sinks Sinks `json:"sinks,omitempty" yaml:"sinks,omitempty"`
}

// NewWriteSyncer produces a zapcore.WriteSyncer using the supplied
// Sink configuration.
//
// The location parameter may be a file system path or the Stdout or Stderr values.
func NewWriteSyncer(location string, sink Sink) (zapcore.WriteSyncer, error) {
	switch {
	case len(location) == 0 || location == Stdout:
		return stdout, nil

	case location == Stderr:
		return stderr, nil

	case sink.DisableRotation:
		f, err := os.OpenFile(location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		return zapcore.Lock(f), nil

	default:
		return NewLumberjack(&lumberjack.Logger{
			Filename:   location,
			MaxSize:    sink.MaxSize,
			MaxAge:     sink.MaxAge,
			MaxBackups: sink.MaxBackups,
			LocalTime:  sink.LocalTime,
			Compress:   sink.Compress,
		}), nil
	}
}

// NewSingleCore constructs a zapcore.Core logger from the supplied Sink configuration
//
// See: https://pkg.go.dev/go.uber.org/zap@v1.15.0/zapcore?tab=doc#NewCore
func NewSingleCore(location string, sink Sink) (zapcore.Core, error) {
	ws, err := NewWriteSyncer(location, sink)
	if err != nil {
		return zapcore.NewNopCore(), err
	}

	encoder := sink.Encoding.NewEncoder(
		sink.EncoderConfig,
	)

	var level zapcore.LevelEnabler = sink.Level
	if level == (zap.AtomicLevel{}) {
		level = zapcore.ErrorLevel
	}

	return zapcore.NewCore(encoder, ws, level), nil
}

// NewCore constructs a zapcore.Core that emits output to all the given Sink
// locations. zapcore.NewTee is used to produce an aggregate Core logger.
//
// This function never returns a nil Core.  If an error occurs, it returns
// a nop Core.
//
// See: https://pkg.go.dev/go.uber.org/zap@v1.15.0/zapcore?tab=doc#NewTee
// See: https://pkg.go.dev/go.uber.org/zap@v1.15.0/zapcore?tab=doc#NewNopCore
func NewCore(sinks Sinks) (zapcore.Core, error) {
	cores := make([]zapcore.Core, 0, len(sinks))
	for location, sink := range sinks {
		core, err := NewSingleCore(location, sink)
		if err != nil {
			return zapcore.NewNopCore(), err
		}

		cores = append(cores, core)
	}

	return zapcore.NewTee(cores...), nil
}

// New marries the zap package with sallust.  It uses a mapping of Sinks
// to create a Core, then creates a zap.Logger with that core and an optional
// set of zap options.
//
// This function never returns a nil logger.  If an error occurs, a
// nop logger is returned.
func New(o Options, zapOpts ...zap.Option) (*zap.Logger, error) {
	core, err := NewCore(o.Sinks)
	if err != nil {
		return zap.NewNop(), err
	}

	// see: https://github.com/uber-go/zap/blob/v1.15.0/config.go#L196
	opts := []zap.Option{zap.ErrorOutput(stderr)}
	if o.Development {
		opts = append(opts, zap.Development())
	}

	if !o.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	if !o.DisableStacktrace {
		stackLevel := zap.ErrorLevel
		if o.Development {
			stackLevel = zap.WarnLevel
		}

		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if len(o.InitialFields) > 0 {
		// the following mimics the zap package's algorithm:
		// produce a set of zap.Field in sorted order, sorted by key
		var (
			fields = make([]zap.Field, 0, len(o.InitialFields))
			keys   = make([]string, 0, len(o.InitialFields))
		)

		for k := range o.InitialFields {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		for _, k := range keys {
			fields = append(fields, zap.Any(k, o.InitialFields[k]))
		}

		opts = append(opts, zap.Fields(fields...))
	}

	opts = append(opts, zapOpts...)
	return zap.New(core, opts...), nil
}
