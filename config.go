// SPDX-FileCopyrightText: 2024 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package sallust

import (
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// DefaultMessageKey is the default value for EncoderConfig.MessageKey.  This value
	// is not used if EncoderConfig.DisableDefaultKeys is true.
	DefaultMessageKey = "msg"

	// DefaultLevelKey is the default value for EncoderConfig.LevelKey.  This value
	// is not used if EncoderConfig.DisableDefaultKeys is true.
	DefaultLevelKey = "level"

	// DefaultTimeKey is the default value for EncoderConfig.TimeKey.  This value
	// is not used if EncoderConfig.DisableDefaultKeys is true.
	DefaultTimeKey = "ts"

	// DefaultNameKey is the default value for EncoderConfig.NameKey.  This value
	// is not used if EncoderConfig.DisableDefaultKeys is true.
	DefaultNameKey = "name"

	// Stdout is the reserved zap output path name that corresponds to stdout.
	Stdout = "stdout"

	// Stderr is the reserved zap output path name that corresponds to stderr.
	Stderr = "stderr"
)

// EncoderConfig is an analog to zap.EncoderConfig.  This type is friendlier
// to unmarshaling from maps, as encoding.TextUnmarshaler is not honored in those cases.
//
// See: https://pkg.go.dev/go.uber.org/zap/zapcore#EncoderConfig
type EncoderConfig struct {
	// DisableDefaultKeys disables the convenience defaulting of certain log keys.
	// Useful when you want to turn off one of those keys, but explicitly set the others.
	DisableDefaultKeys bool `json:"disableDefaultKeys" yaml:"disableDefaultKeys"`

	// MessageKey is the logging key for the log message.  If unset and if DisableDefaultKeys is true,
	// messages are not inserted into log output.
	MessageKey string `json:"messageKey" yaml:"messageKey"`

	// LevelKey is the logging key for the log level.  If unset and if DisableDefaultKeys is true,
	// log levels are not inserted into log output.
	LevelKey string `json:"levelKey" yaml:"levelKey"`

	// TimeKey is the logging key for the log timestamp.  If unset and if DisableDefaultKeys is true,
	// timestamps are not inserted into log output.
	TimeKey string `json:"timeKey" yaml:"timeKey"`

	// NameKey is the logging key for the logger name.  If unset and if DisableDefaultKeys is true,
	//	logger names are not inserted into log output.
	NameKey string `json:"nameKey" yaml:"nameKey"`

	// CallerKey is the logging key for the caller of the logging method.  If unset, callers are not
	// inserted into log output.
	//
	// Note that Config.DisableCaller, if set, will also prevent callers in each log record.
	// This difference is that Config.DisableCaller shuts off the code that determines the caller,
	// while this field simply doesn't output the caller even though it may have been computed.
	CallerKey string `json:"callerKey" yaml:"callerKey"`

	// FunctionKey is the logging key for the function which called the logging method.  If unset, functions are not
	// inserted into log output.
	//
	// As with CallerKey, Config.DisableCaller also affects whether functions are output.
	FunctionKey string `json:"functionKey" yaml:"functionKey"`

	// StacktraceKey is the logging key for stacktraces for warn, error, and panics.  If unset,
	// stacktraces are never produced.
	StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`

	// LineEnding is the US-ASCII string that terminates each log record.  By default,
	// a single '\n' is used.
	LineEnding string `json:"lineEnding" yaml:"lineEnding"`

	// EncodeLevel determines how levels are represented.  If unset, LowercaseLevelEncoder is used.
	//
	// See: https://pkg.go.dev/go.uber.org/zap/zapcore#LowercaseLevelEncoder
	EncodeLevel string `json:"levelEncoder" yaml:"levelEncoder" mapstructure:"levelEncoder"`

	// EncodeTime determines how timestamps are represented.  If unset, RFC3339TimeEncoder is used.
	//
	// See: https://pkg.go.dev/go.uber.org/zap/zapcore#RFC3339TimeEncoder
	EncodeTime string `json:"timeEncoder" yaml:"timeEncoder" mapstructure:"timeEncoder"`

	// EncodeDuration determines how time durations are represented.  If unset,
	// StringDurationEncoder is used.
	//
	// See: https://pkg.go.dev/go.uber.org/zap/zapcore#StringDurationEncoder
	EncodeDuration string `json:"durationEncoder" yaml:"durationEncoder" mapstructure:"durationEncoder"`

	// EncodeCaller determines how callers are represented.  If unset,
	// FullCallerEncoder is used.
	//
	// See: https://pkg.go.dev/go.uber.org/zap/zapcore#FullCallerEncoder
	EncodeCaller string `json:"callerEncoder" yaml:"callerEncoder" mapstructure:"callerEncoder"`

	// EncodeName determines how logger names are represented.  If unset,
	// FullNameEncoder is used.
	//
	// See: https://pkg.go.dev/go.uber.org/zap/zapcore#FullNameEncoder
	EncodeName string `json:"nameEncoder" yaml:"nameEncoder" mapstructure:"nameEncoder"`

	// Configures the field separator used by the console encoder. Defaults
	// to tab.
	ConsoleSeparator string `json:"consoleSeparator" yaml:"consoleSeparator"`
}

func applyEncoderConfigDefaults(zec *zapcore.EncoderConfig) {
	if len(zec.MessageKey) == 0 {
		zec.MessageKey = "msg"
	}

	if len(zec.LevelKey) == 0 {
		zec.LevelKey = "level"
	}

	if len(zec.TimeKey) == 0 {
		zec.TimeKey = "ts"
	}

	if len(zec.NameKey) == 0 {
		zec.NameKey = "name"
	}
}

// NewZapcoreEncoderConfig converts this instance into a zapcore.EncoderConfig.
//
// In order to ease configuration of zap, this method implements a few conveniences
// on the returned zapcore.EncoderConfig:
//
// (1) Each of the EncodeXXX fields is defaulted to a sane value.  Leaving them unset
// does not raise an error.
//
// (2) Several logging key fields are defaulted.  This defaulting can be turned off by
// setting DisableDefaultKeys to true.  The fields that are defaulted are:  MessageKey,
// LevelKey, TimeKey, and NameKey.  The other logging key fields, such as CallerKey,
// are not defaulted.  It's common to leave them turned off for performance or preference.
//
// This method returns an error, as the various UnmarshalText methods it uses return errors.
// However, in actual practice this method returns a nil error since zapcore falls back to
// known defaults for any unrecognized text.  This method still checks errors internally,
// in case zapcore changes in the future.
func (ec EncoderConfig) NewZapcoreEncoderConfig() (zec zapcore.EncoderConfig, err error) {
	zec = zapcore.EncoderConfig{
		MessageKey:       ec.MessageKey,
		LevelKey:         ec.LevelKey,
		TimeKey:          ec.TimeKey,
		NameKey:          ec.NameKey,
		CallerKey:        ec.CallerKey,
		FunctionKey:      ec.FunctionKey,
		StacktraceKey:    ec.StacktraceKey,
		LineEnding:       ec.LineEnding,
		ConsoleSeparator: ec.ConsoleSeparator,
	}

	if !ec.DisableDefaultKeys {
		applyEncoderConfigDefaults(&zec)
	}

	if len(ec.EncodeLevel) > 0 {
		err = zec.EncodeLevel.UnmarshalText([]byte(ec.EncodeLevel))
	} else {
		zec.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	if err == nil {
		if len(ec.EncodeTime) > 0 {
			err = zec.EncodeTime.UnmarshalText([]byte(ec.EncodeTime))
		} else {
			zec.EncodeTime = zapcore.RFC3339TimeEncoder
		}
	}

	if err == nil {
		if len(ec.EncodeDuration) > 0 {
			err = zec.EncodeDuration.UnmarshalText([]byte(ec.EncodeDuration))
		} else {
			zec.EncodeDuration = zapcore.StringDurationEncoder
		}
	}

	if err == nil {
		if len(ec.EncodeCaller) > 0 {
			err = zec.EncodeCaller.UnmarshalText([]byte(ec.EncodeCaller))
		} else {
			zec.EncodeCaller = zapcore.FullCallerEncoder
		}
	}

	if err == nil {
		if len(ec.EncodeName) > 0 {
			err = zec.EncodeName.UnmarshalText([]byte(ec.EncodeName))
		} else {
			zec.EncodeName = zapcore.FullNameEncoder
		}
	}

	return
}

// Config describes the set of options for building a single zap.Logger.  Most of these
// fields correspond with zap.Config.  Use of this type is optional.  It simply provides
// easier configuration for certain features like log rotation.  This type is also easier
// to use with libraries like spf13/viper, which unmarshal from a map[string]interface{}
// instead of directly from a file.
//
// A Config instance is converted to a zap.Config by applying certain features,
// such as log rotation.  Ultimately, zap.Config.Build is used to actually construct
// the logger.
//
// See: https://pkg.go.dev/go.uber.org/zap?tab=doc#Config.Build
type Config struct {
	// Level is the log level, which is converted to a zap.AtomicLevel.  If unset,
	// info level is assumed.
	Level string `json:"level" yaml:"level"`

	// Development corresponds to zap.Config.Development
	Development bool `json:"development" yaml:"development"`

	// DisableCaller corresponds to zap.Config.DisableCaller
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`

	// DisableStacktrace corresponds to zap.Config.DisableStacktrace
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`

	// Sampling corresponds to zap.Config.Sampling.  No custom type is necessary
	// here because zap.SamplingConfig uses primitive types.
	Sampling *zap.SamplingConfig `json:"samplingConfig" yaml:"samplingConfig"`

	// Encoding corresponds to zap.Config.Encoding.  If this is unset, and if Development
	// is false, "json" is used.  "console" is the other built-in value for this field,
	// and other encodings can be registered via the zap package.
	//
	// See: https://pkg.go.dev/go.uber.org/zap#RegisterEncoder
	Encoding string `json:"encoding" yaml:"encoding"`

	// EncoderConfig corresponds to zap.Config.EncoderConfig.  A custom type is used
	// here to make integration with libraries like spf13/viper much easier.
	EncoderConfig EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`

	// OutputPaths are the set of sinks for log output.  This field corresponds to
	// zap.Config.OutputPaths.  If unset, all logging output is discarded.  There is
	// no default for this field.
	//
	// Each output path will have environment variable references expanded unless
	// DisablePathExpansion is true.
	//
	// If Rotation is set, then each output path that is a system file will undergo
	// log file rotation.
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`

	// ErrorOutputPaths are the set of sinks for zap's internal messages.  This field
	// corresponds to zap.Config.ErrorOutputPaths.  If unset, Stderr is assumed.
	//
	// As with OutputPaths, environment variable references in each path are expanded
	// unless DisablePathExpansion is true.
	//
	// If Rotation is set, then each output path that is a system file will undergo
	// log file rota
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`

	// InitialFields corresponds to zap.Config.InitialFields.  Note that when unmarshaling
	// from spf13/viper, all keys in this map will be lowercased.
	//
	// Any fields set here will be set on all loggers derived from this configuration.
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`

	// DisablePathExpansion controls whether the paths in OutputPaths and ErrorOutputPaths
	// are expanded.  If this field is set to true, Mapping is ignored and no
	// expansion, even with environment variables, is performed.
	DisablePathExpansion bool `json:"disablePathExpansion" yaml:"disablePathExpansion"`

	// Permissions is the optional nix-style file permissions to use when creating log files.
	// If supplied, this value must be parseable via ParsePermissions.  If this field is unset,
	// zap and lumberjack will control what permissions new log files have.
	Permissions string `json:"permissions" yaml:"permissions"`

	// Mapping is an optional strategy for expanding variables in output paths.
	// If not supplied, os.Getenv is used.
	Mapping func(string) string `json:"-" yaml:"-"`

	// Rotation describes the set of log file rotation options.  This field is optional,
	// and if unset log files are not rotated.
	Rotation *Rotation `json:"rotation,omitempty" yaml:"rotation,omitempty"`
}

func applyConfigDefaults(zc *zap.Config) {
	if len(zc.Encoding) == 0 {
		zc.Encoding = "json"
	}

	if zc.Development && len(zc.OutputPaths) == 0 {
		// NOTE: difference from zap ... in development they send output to stderr
		zc.OutputPaths = []string{Stdout}
	}

	if len(zc.ErrorOutputPaths) == 0 {
		zc.ErrorOutputPaths = []string{Stderr}
	}

	// NOTE: can't compare the Level with nil very easily, so just
	// unconditionally set this default.  It will be overwritten
	// by decoding code if appropriate.
	zc.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
}

// ensureExists makes sure the given path exists with the specified permissions.
// If the path has already been created or if perms is 0, this function won't do anything.
//
// The path is treated as a URI in a similar fashion to zap.Open.
func ensureExists(path string, perms fs.FileMode) (err error) {
	if perms == 0 {
		return
	}

	var f *os.File
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	switch {
	case path == Stdout:
		fallthrough

	case path == Stderr:
		break

	// Windows hack:  filepath.Abs will return false outside of Windows
	// for many paths.  This just makes sure we don't have to do a bunch
	// of platform-specific nonsense.
	case filepath.IsAbs(path):
		f, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, perms)

	default:
		var url *url.URL
		url, err = url.Parse(path)
		if err == nil {
			f, err = os.OpenFile(url.Path, os.O_CREATE|os.O_WRONLY, perms)
		}
	}

	return
}

// NewZapConfig creates a zap.Config enriched with features from these Options.
// Primarily, this involves creating lumberjack URLs so that the registered sink
// will create the appropriate infrastructure to do log file rotation.
//
// This method also enforces the Permissions field.  Any output or error path
// will be created initially with the configured file permissions.  This allows
// both zap's file sink and the custom lumberjack sink in this package to honor
// custom permissions.
func (c Config) NewZapConfig() (zc zap.Config, err error) {
	zc = zap.Config{
		Development:       c.Development,
		DisableCaller:     c.DisableCaller,
		DisableStacktrace: c.DisableStacktrace,
		Encoding:          c.Encoding,
		OutputPaths:       append([]string{}, c.OutputPaths...),
		ErrorOutputPaths:  append([]string{}, c.ErrorOutputPaths...),
	}

	if c.Sampling != nil {
		zc.Sampling = new(zap.SamplingConfig)
		*zc.Sampling = *c.Sampling
	}

	if len(c.InitialFields) > 0 {
		zc.InitialFields = make(map[string]interface{}, len(c.InitialFields))
		for k, v := range c.InitialFields {
			zc.InitialFields[k] = v
		}
	}

	applyConfigDefaults(&zc)

	if len(c.Level) > 0 {
		var l zapcore.Level
		err = l.UnmarshalText([]byte(c.Level))
		if err == nil {
			zc.Level = zap.NewAtomicLevelAt(l)
		}
	}

	var perms fs.FileMode
	perms, err = ParsePermissions(c.Permissions)

	if err == nil {
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
	}

	// Iterate over the transformed paths and ensure that any URIs that refer to
	// files are created with relevant permissions.
	for i := 0; err == nil && i < len(zc.OutputPaths); i++ {
		err = ensureExists(zc.OutputPaths[i], perms)
	}

	for i := 0; err == nil && i < len(zc.ErrorOutputPaths); i++ {
		err = ensureExists(zc.ErrorOutputPaths[i], perms)
	}

	if err == nil {
		zc.EncoderConfig, err = c.EncoderConfig.NewZapcoreEncoderConfig()
	}

	return
}

// Build behaves similarly to zap.Config.Build.  It uses the configuration created
// by NewZapConfig to build the root logger.
func (c Config) Build(opts ...zap.Option) (l *zap.Logger, err error) {
	var zc zap.Config
	zc, err = c.NewZapConfig()
	if err == nil {
		l, err = zc.Build(opts...)
	}

	return
}
