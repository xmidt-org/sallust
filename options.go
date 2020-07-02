package sallust

import (
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

	// Rotation describes the set of log file rotation options.  This field is optional,
	// and if unset log files are not rotated.
	Rotation *Rotation `json:"rotation,omitempty" yaml:"rotation,omitempty"`
}

// NewZapConfig creates a zap.Config enriched with features from these Options.
// Primarily, this involves creating lumberjack URLs so that the registered sink
// will create the appropriate infrastructure to do log file rotation.
func (o Options) NewZapConfig() (zap.Config, error) {
	zapConfig := zap.Config{
		Level:             o.Level,
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling:          o.Sampling,
		Encoding:          o.Encoding,
		EncoderConfig:     o.EncoderConfig,
		OutputPaths:       o.OutputPaths,
		ErrorOutputPaths:  o.ErrorOutputPaths,
		InitialFields:     o.InitialFields,
	}

	if zapConfig.Level == (zap.AtomicLevel{}) {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	if o.Rotation != nil {
		zapConfig.OutputPaths = o.Rotation.TransformURLs(o.OutputPaths...)
		zapConfig.ErrorOutputPaths = o.Rotation.TransformURLs(o.ErrorOutputPaths...)
	}

	return zapConfig, nil
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
