package sallust

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapcoreSuite is an embeddable suite that contains common functionality
// for the test suites involving configuration.
type ZapcoreSuite struct {
	suite.Suite
}

func (suite *ZapcoreSuite) assertEncoderConfigDefaults(zec zapcore.EncoderConfig) {
	suite.Equal(DefaultMessageKey, zec.MessageKey)
	suite.Equal(DefaultLevelKey, zec.LevelKey)
	suite.Equal(DefaultTimeKey, zec.TimeKey)
	suite.Equal(DefaultNameKey, zec.NameKey)
	suite.Empty(zec.CallerKey)
	suite.Empty(zec.FunctionKey)
	suite.Empty(zec.StacktraceKey)
	suite.NotNil(zec.EncodeLevel)
	suite.NotNil(zec.EncodeTime)
	suite.NotNil(zec.EncodeDuration)
	suite.NotNil(zec.EncodeCaller)
	suite.NotNil(zec.EncodeName)
}

type EncoderConfigSuite struct {
	ZapcoreSuite
}

func (suite *EncoderConfigSuite) TestDefaults() {
	var (
		ec       EncoderConfig
		zec, err = ec.NewZapcoreEncoderConfig()
	)

	suite.Require().NoError(err)
	suite.assertEncoderConfigDefaults(zec)
}

func (suite *EncoderConfigSuite) TestCustom() {
	var (
		ec = EncoderConfig{
			MessageKey:     "message_key",
			LevelKey:       "level_key",
			TimeKey:        "time_key",
			NameKey:        "name_key",
			CallerKey:      "caller_key",
			FunctionKey:    "function_key",
			StacktraceKey:  "stacktrace_key",
			EncodeLevel:    "capital",
			EncodeTime:     "iso8601",
			EncodeDuration: "nanos",  // doesn't matter, zapcore will unmarshal anything but "string" to nanos
			EncodeCaller:   "short",  // doesn't matter, zapcore will unmarshal anything but "full" to short
			EncodeName:     "custom", // doesn't matter, zapcore unmarshals everything as full name

			LineEnding:       "foo",
			ConsoleSeparator: "bar",
		}

		zec, err = ec.NewZapcoreEncoderConfig()
	)

	suite.Require().NoError(err)

	suite.Equal("message_key", zec.MessageKey)
	suite.Equal("level_key", zec.LevelKey)
	suite.Equal("time_key", zec.TimeKey)
	suite.Equal("name_key", zec.NameKey)
	suite.Equal("caller_key", zec.CallerKey)
	suite.Equal("function_key", zec.FunctionKey)
	suite.Equal("stacktrace_key", zec.StacktraceKey)
	suite.NotNil(zec.EncodeLevel)
	suite.NotNil(zec.EncodeTime)
	suite.NotNil(zec.EncodeDuration)
	suite.NotNil(zec.EncodeCaller)
	suite.NotNil(zec.EncodeName)

	suite.Equal("foo", zec.LineEnding)
	suite.Equal("bar", zec.ConsoleSeparator)
}

func (suite *EncoderConfigSuite) TestDisableDefaultKeys() {
	var (
		ec = EncoderConfig{
			DisableDefaultKeys: true,
		}

		zec, err = ec.NewZapcoreEncoderConfig()
	)

	suite.Require().NoError(err)

	suite.Empty(zec.MessageKey)
	suite.Empty(zec.LevelKey)
	suite.Empty(zec.TimeKey)
	suite.Empty(zec.NameKey)
	suite.Empty(zec.CallerKey)
	suite.Empty(zec.FunctionKey)
	suite.Empty(zec.StacktraceKey)
	suite.NotNil(zec.EncodeLevel)
	suite.NotNil(zec.EncodeTime)
	suite.NotNil(zec.EncodeDuration)
	suite.NotNil(zec.EncodeCaller)
	suite.NotNil(zec.EncodeName)
}

func TestEncoderConfig(t *testing.T) {
	suite.Run(t, new(EncoderConfigSuite))
}

type ConfigSuite struct {
	ZapcoreSuite
}

func (suite *ConfigSuite) TestDefaults() {
	var c Config

	zc, err := c.NewZapConfig()
	suite.Require().NoError(err)

	suite.Equal(zapcore.InfoLevel, zc.Level.Level())
	suite.False(zc.Development)
	suite.False(zc.DisableCaller)
	suite.False(zc.DisableStacktrace)
	suite.Equal("json", zc.Encoding)
	suite.Empty(zc.OutputPaths)
	suite.Equal([]string{"stderr"}, zc.ErrorOutputPaths)
	suite.Nil(zc.Sampling)
	suite.Empty(zc.InitialFields)

	zec, err := c.EncoderConfig.NewZapcoreEncoderConfig()
	suite.Require().NoError(err)
	suite.assertEncoderConfigDefaults(zec)
}

func (suite *ConfigSuite) TestCustom() {
	c := Config{
		Level:             "debug",
		Development:       true,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    1,
			Thereafter: 10,
		},
		Encoding:         "console",
		OutputPaths:      []string{"/var/log/test/test.log"},
		ErrorOutputPaths: []string{"stdout"},
		InitialFields: map[string]interface{}{
			"name":  "value",
			"slice": []string{"1", "2"},
		},
	}

	zc, err := c.NewZapConfig()
	suite.Require().NoError(err)

	suite.Equal(zapcore.DebugLevel, zc.Level.Level())
	suite.True(zc.Development)
	suite.True(zc.DisableCaller)
	suite.True(zc.DisableStacktrace)
	suite.Equal("console", zc.Encoding)
	suite.Equal([]string{"/var/log/test/test.log"}, zc.OutputPaths)
	suite.Equal([]string{"stdout"}, zc.ErrorOutputPaths)
	suite.Equal(
		zap.SamplingConfig{
			Initial:    1,
			Thereafter: 10,
		},
		*zc.Sampling,
	)

	suite.Equal(
		map[string]interface{}{
			"name":  "value",
			"slice": []string{"1", "2"},
		},
		zc.InitialFields,
	)

	zec, err := c.EncoderConfig.NewZapcoreEncoderConfig()
	suite.Require().NoError(err)
	suite.assertEncoderConfigDefaults(zec)
}

func (suite *ConfigSuite) TestDevelopmentDefaults() {
	c := Config{
		Development: true,
	}

	zc, err := c.NewZapConfig()
	suite.Require().NoError(err)

	suite.Equal(zapcore.InfoLevel, zc.Level.Level())
	suite.True(zc.Development)
	suite.False(zc.DisableCaller)
	suite.False(zc.DisableStacktrace)
	suite.Equal("json", zc.Encoding)
	suite.Equal([]string{"stdout"}, zc.OutputPaths)
	suite.Equal([]string{"stderr"}, zc.ErrorOutputPaths)
	suite.Nil(zc.Sampling)
	suite.Empty(zc.InitialFields)

	zec, err := c.EncoderConfig.NewZapcoreEncoderConfig()
	suite.Require().NoError(err)
	suite.assertEncoderConfigDefaults(zec)
}

func (suite *ConfigSuite) TestBuildSimple() {
	var (
		buffer bytes.Buffer

		c = Config{
			Development: true,
		}
	)

	// create an encoder config to replace the one created by the zap package
	// so that we can run assertions
	zec, err := EncoderConfig{}.NewZapcoreEncoderConfig()
	suite.Require().NoError(err)

	l, err := c.Build(
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(zec),
				zapcore.AddSync(&buffer),
				zapcore.DebugLevel,
			)
		}),
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(l)
	l.Info("test message")
	suite.Greater(buffer.Len(), 0)
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
