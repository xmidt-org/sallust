package sallust

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config describes the set of options for building a single zap.Logger.  Most of these
// fields are exactly the same as zap.Config.  Use of this type is optional.  It simply provides
// easier configuration for certain features like log rotation.
//
// A Config instance is converted to a zap.Config by applying certain features,
// such as log rotation.  Ultimately, zap.Config.Build is used to actually construct
// the logger.
//
// See: https://pkg.go.dev/go.uber.org/zap?tab=doc#Config.Build
type Config struct {
	// Config embeds all the usual fields from zap.  This is marked to squash so that
	// these fields don't have to be nested.
	zap.Config `mapstructure:",squash"`

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
func (c Config) NewZapConfig() (zc zap.Config, err error) {
	zc = c.Config

	pt := PathTransformer{
		Rotation: c.Rotation,
	}

	if !c.DisablePathExpansion {
		pt.Mapping = c.Mapping
		if pt.Mapping == nil {
			pt.Mapping = os.Getenv
		}
	}

	zc.OutputPaths, err = ApplyTransform(pt.Transform, zc.OutputPaths...)
	if err == nil {
		zc.ErrorOutputPaths, err = ApplyTransform(pt.Transform, zc.ErrorOutputPaths...)
	}

	if err == nil {
		// apply certain defaults.  zap will error out if most of these are unset,
		// but we want to add sane defaults rather than require someone to configure
		// every field.
		if zc.Level == (zap.AtomicLevel{}) {
			zc.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		}

		if len(zc.Encoding) == 0 {
			zc.Encoding = "json"
		}

		if len(zc.OutputPaths) == 0 {
			// NOTE: difference from zap ... in development they send output to stderr
			zc.OutputPaths = []string{"stdout"}
		}

		if len(zc.ErrorOutputPaths) == 0 {
			zc.ErrorOutputPaths = []string{"stderr"}
		}

		if zc.EncoderConfig.EncodeLevel == nil {
			zc.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}

		if zc.EncoderConfig.EncodeTime == nil {
			zc.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		}

		if zc.EncoderConfig.EncodeDuration == nil {
			zc.EncoderConfig.EncodeDuration = zapcore.NanosDurationEncoder
		}

		if zc.EncoderConfig.EncodeCaller == nil {
			zc.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		}

		if zc.EncoderConfig.EncodeName == nil {
			// technically, we don't need to do this as zap will fall back to FullNameEncoder.
			// but, this keeps this field consistent with the other EncodeXXX fields
			zc.EncoderConfig.EncodeName = zapcore.FullNameEncoder
		}
	}

	return
}

// NewLogger behaves similarly to zap.Config.Build.  It uses the configuration created
// by NewZapConfig to build the root logger.
func (c Config) NewLogger(opts ...zap.Option) (*zap.Logger, error) {
	zapConfig, err := c.NewZapConfig()
	if err != nil {
		return nil, err
	}

	return zapConfig.Build(opts...)
}
