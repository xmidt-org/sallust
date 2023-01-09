package sallust

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func assertZapcoreEncoderConfigDefaults(assert *assert.Assertions, zec zapcore.EncoderConfig) {
	assert.Equal(DefaultMessageKey, zec.MessageKey)
	assert.Equal(DefaultLevelKey, zec.LevelKey)
	assert.Equal(DefaultTimeKey, zec.TimeKey)
	assert.Equal(DefaultNameKey, zec.NameKey)
	assert.Empty(zec.CallerKey)
	assert.Empty(zec.FunctionKey)
	assert.Empty(zec.StacktraceKey)
	assert.NotNil(zec.EncodeLevel)
	assert.NotNil(zec.EncodeTime)
	assert.NotNil(zec.EncodeDuration)
	assert.NotNil(zec.EncodeCaller)
	assert.NotNil(zec.EncodeName)
}

func testEncoderConfigDefaults(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		ec      EncoderConfig
	)

	zec, err := ec.NewZapcoreEncoderConfig()
	require.NoError(err)
	assertZapcoreEncoderConfigDefaults(assert, zec)
}

func testEncoderConfigCustom(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		ec      = EncoderConfig{
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
	)

	zec, err := ec.NewZapcoreEncoderConfig()
	require.NoError(err)

	assert.Equal("message_key", zec.MessageKey)
	assert.Equal("level_key", zec.LevelKey)
	assert.Equal("time_key", zec.TimeKey)
	assert.Equal("name_key", zec.NameKey)
	assert.Equal("caller_key", zec.CallerKey)
	assert.Equal("function_key", zec.FunctionKey)
	assert.Equal("stacktrace_key", zec.StacktraceKey)
	assert.NotNil(zec.EncodeLevel)
	assert.NotNil(zec.EncodeTime)
	assert.NotNil(zec.EncodeDuration)
	assert.NotNil(zec.EncodeCaller)
	assert.NotNil(zec.EncodeName)

	assert.Equal("foo", zec.LineEnding)
	assert.Equal("bar", zec.ConsoleSeparator)
}

func testEncoderConfigDisableDefaultKeys(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		ec      = EncoderConfig{
			DisableDefaultKeys: true,
		}
	)

	zec, err := ec.NewZapcoreEncoderConfig()
	require.NoError(err)

	assert.Empty(zec.MessageKey)
	assert.Empty(zec.LevelKey)
	assert.Empty(zec.TimeKey)
	assert.Empty(zec.NameKey)
	assert.Empty(zec.CallerKey)
	assert.Empty(zec.FunctionKey)
	assert.Empty(zec.StacktraceKey)
	assert.NotNil(zec.EncodeLevel)
	assert.NotNil(zec.EncodeTime)
	assert.NotNil(zec.EncodeDuration)
	assert.NotNil(zec.EncodeCaller)
	assert.NotNil(zec.EncodeName)
}

func TestEncoderConfig(t *testing.T) {
	t.Run("Defaults", testEncoderConfigDefaults)
	t.Run("Custom", testEncoderConfigCustom)
	t.Run("DisableDefaultKeys", testEncoderConfigDisableDefaultKeys)
}

func testConfigNewZapConfigDefaults(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)

		c Config
	)

	zc, err := c.NewZapConfig()
	require.NoError(err)

	assert.Equal(zapcore.InfoLevel, zc.Level.Level())
	assert.False(zc.Development)
	assert.False(zc.DisableCaller)
	assert.False(zc.DisableStacktrace)
	assert.Equal("json", zc.Encoding)
	assert.Empty(zc.OutputPaths)
	assert.Equal([]string{"stderr"}, zc.ErrorOutputPaths)
	assert.Nil(zc.Sampling)
	assert.Empty(zc.InitialFields)

	zec, err := c.EncoderConfig.NewZapcoreEncoderConfig()
	require.NoError(err)
	assertZapcoreEncoderConfigDefaults(assert, zec)
}

func testConfigNewZapConfigCustom(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)

		c = Config{
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
	)

	zc, err := c.NewZapConfig()
	require.NoError(err)

	assert.Equal(zapcore.DebugLevel, zc.Level.Level())
	assert.True(zc.Development)
	assert.True(zc.DisableCaller)
	assert.True(zc.DisableStacktrace)
	assert.Equal("console", zc.Encoding)
	assert.Equal([]string{"/var/log/test/test.log"}, zc.OutputPaths)
	assert.Equal([]string{"stdout"}, zc.ErrorOutputPaths)
	assert.Equal(
		zap.SamplingConfig{
			Initial:    1,
			Thereafter: 10,
		},
		*zc.Sampling,
	)

	assert.Equal(
		map[string]interface{}{
			"name":  "value",
			"slice": []string{"1", "2"},
		},
		zc.InitialFields,
	)

	zec, err := c.EncoderConfig.NewZapcoreEncoderConfig()
	require.NoError(err)
	assertZapcoreEncoderConfigDefaults(assert, zec)
}

func testConfigNewZapConfigDevelopmentDefaults(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)

		c = Config{
			Development: true,
		}
	)

	zc, err := c.NewZapConfig()
	require.NoError(err)

	assert.Equal(zapcore.InfoLevel, zc.Level.Level())
	assert.True(zc.Development)
	assert.False(zc.DisableCaller)
	assert.False(zc.DisableStacktrace)
	assert.Equal("json", zc.Encoding)
	assert.Equal([]string{"stdout"}, zc.OutputPaths)
	assert.Equal([]string{"stderr"}, zc.ErrorOutputPaths)
	assert.Nil(zc.Sampling)
	assert.Empty(zc.InitialFields)

	zec, err := c.EncoderConfig.NewZapcoreEncoderConfig()
	require.NoError(err)
	assertZapcoreEncoderConfigDefaults(assert, zec)
}

func testConfigNewZapConfig(t *testing.T) {
	t.Run("Defaults", testConfigNewZapConfigDefaults)
	t.Run("Custom", testConfigNewZapConfigCustom)
	t.Run("DevelopmentDefaults", testConfigNewZapConfigDevelopmentDefaults)
}

func testConfigBuildSimple(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)

		buffer bytes.Buffer

		c = Config{
			Development: true,
		}
	)

	// create an encoder config to replace the one created by the zap package
	// so that we can run assertions
	zec, err := EncoderConfig{}.NewZapcoreEncoderConfig()
	require.NoError(err)

	l, err := c.Build(
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(zec),
				zapcore.AddSync(&buffer),
				zapcore.DebugLevel,
			)
		}),
	)

	require.NoError(err)
	require.NotNil(l)
	l.Info("test message")
	assert.Greater(buffer.Len(), 0)
}

func testConfigBuild(t *testing.T) {
	t.Run("Simple", testConfigBuildSimple)
}

func TestConfig(t *testing.T) {
	t.Run("NewZapConfig", testConfigNewZapConfig)
	t.Run("Build", testConfigBuild)
}
