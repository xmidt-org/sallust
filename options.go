package sallust

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options describes the set of options for building a single zap.Logger.  Most of these
// fields are exactly the same as zap.Config.
//
// An Options instance is converted to a zap.Config by applying certain features,
// such as log rotation.  Ultimately, zap.Config.Build is used to actually construct
// the logger.
//
// See: https://pkg.go.dev/go.uber.org/zap?tab=doc#Config.Build
type Options struct {
	// Level is the dynamic log level.  Unlike zap, this field is defaulted to zapcore.ErrorLevel.
	// No error will be returned if this field is left unset.
	Level zap.AtomicLevel `json:"level" yaml:"level"`

	// Development is the same as zap.Config.Development
	Development bool `json:"development" yaml:"development"`

	// DisableCaller is the same as zap.Config.DisableCaller
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`

	// DisableStacktrace is the same as zap.Config.DisableStacktrace
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`

	// Sampling is the same as zap.Config.Sampling
	Sampling *zap.SamplingConfig `json:"sampling" yaml:"sampling"`

	// Encoding is the same as zap.Config.Encoding
	Encoding string `json:"encoding" yaml:"encoding"`

	// EncoderConfig is the same as zap.Config.EncoderConfig
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`

	// OutputPaths is the same as zap.Config.OutputPaths.  If the Rotation field is
	// specified, this slice is preprocessed to produce lumberjack-based zap.Sink objects.
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`

	// ErrorOutputPaths is the same as zap.Config.ErrorOutputPaths
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`

	// InitialFields is the same as zap.Config.InitialFields
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`

	// DisablePathExpansion controls whether the paths in OutputPaths and ErrorOutputPaths
	// are expanded.  If this field is set to true, Mapping is ignored and no
	// expansion, even with environment variables, is performed.
	DisablePathExpansion bool `json:"disablePathExpansion" yaml:"disablePathExpansion"`

	// Mapping is an optional strategy for expanding variables in output paths.
	// If not supplied, os.Getenv is used.
	Mapping func(string) string `json:"-" yaml:"-"`

	// Rotation describes the set of log file rotation options.  This field is optional,
	// and if unset log files are not rotated.
	Rotation *Rotation `json:"rotation,omitempty" yaml:"rotation,omitempty"`
}

// NewZapConfig creates a zap.Config enriched with features from these Options.
// Primarily, this involves creating lumberjack URLs so that the registered sink
// will create the appropriate infrastructure to do log file rotation.
func (o Options) NewZapConfig() (zap.Config, error) {
	pt := PathTransformer{
		Rotation: o.Rotation,
	}

	if !o.DisablePathExpansion {
		pt.Mapping = o.Mapping
		if pt.Mapping == nil {
			pt.Mapping = os.Getenv
		}
	}

	outputPaths, err := ApplyTransform(pt.Transform, o.OutputPaths...)
	if err != nil {
		return zap.Config{}, err
	}

	errorOutputPaths, err := ApplyTransform(pt.Transform, o.ErrorOutputPaths...)
	if err != nil {
		return zap.Config{}, err
	}

	level := o.Level
	if level == (zap.AtomicLevel{}) {
		// difference from zap:  we let this be unset, and default it to ErrorLevel
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	return zap.Config{
		Level:             level,
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling:          o.Sampling,
		Encoding:          o.Encoding,
		EncoderConfig:     o.EncoderConfig,
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  errorOutputPaths,
		InitialFields:     o.InitialFields,
	}, nil
}

// NewLogger behaves similarly to zap.Config.Build.  It uses the configuration created
// by NewZapConfig to build the root logger.
func (o Options) NewLogger(opts ...zap.Option) (*zap.Logger, error) {
	zapConfig, err := o.NewZapConfig()
	if err != nil {
		return nil, err
	}

	return zapConfig.Build(opts...)
}
