package sallust

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func testConfigNewZapConfigSuccess(t *testing.T) {
	testData := []struct {
		config   Config
		expected zap.Config
	}{
		{
			config: Config{
				OutputPaths: []string{"/log.json"},
			},
			expected: zap.Config{
				Level:       zap.NewAtomicLevelAt(zapcore.ErrorLevel),
				Encoding:    "json",
				OutputPaths: []string{"/log.json"},
			},
		},
		{
			config: Config{
				OutputPaths: []string{"/log.json"},
				Rotation: &Rotation{
					MaxAge: 10,
				},
			},
			expected: zap.Config{
				Level:       zap.NewAtomicLevelAt(zapcore.ErrorLevel),
				Encoding:    "json",
				OutputPaths: []string{"lumberjack:///log.json?maxAge=10"},
			},
		},
		{
			config: Config{
				Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
				Development:       true,
				DisableCaller:     true,
				DisableStacktrace: true,
				Sampling:          &zap.SamplingConfig{},
				Encoding:          "console",
				EncoderConfig:     zapcore.EncoderConfig{},
				OutputPaths:       []string{"stdout", "file:///log.json"},
				ErrorOutputPaths:  []string{"stderr"},
				InitialFields: map[string]interface{}{
					"foo": "bar",
				},
				Rotation: &Rotation{
					MaxAge: 10,
				},
			},
			expected: zap.Config{
				Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
				Development:       true,
				DisableCaller:     true,
				DisableStacktrace: true,
				Sampling:          &zap.SamplingConfig{},
				Encoding:          "console",
				EncoderConfig:     zapcore.EncoderConfig{},
				OutputPaths:       []string{"stdout", "lumberjack:///log.json?maxAge=10"},
				ErrorOutputPaths:  []string{"stderr"},
				InitialFields: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert      = assert.New(t)
				require     = require.New(t)
				actual, err = record.config.NewZapConfig()
			)

			require.NoError(err)
			assert.Equal(record.expected.Level, actual.Level)
			assert.Equal(record.expected.Development, actual.Development)
			assert.Equal(record.expected.DisableCaller, actual.DisableCaller)
			assert.Equal(record.expected.DisableStacktrace, actual.DisableStacktrace)
			assert.Equal(record.expected.Sampling, actual.Sampling)
			assert.Equal(record.expected.Encoding, actual.Encoding)
			assert.Equal(record.expected.EncoderConfig, actual.EncoderConfig)
			assert.Equal(record.expected.OutputPaths, actual.OutputPaths)
			assert.Equal(record.expected.ErrorOutputPaths, actual.ErrorOutputPaths)
			assert.Equal(record.expected.InitialFields, actual.InitialFields)
		})
	}
}

func testConfigNewZapConfigBadOutputPath(t *testing.T) {
	var (
		assert = assert.New(t)
		config = Config{
			Rotation:    new(Rotation),
			OutputPaths: []string{"#%@(&%(@%XX"},
		}
	)

	_, err := config.NewZapConfig()
	assert.Error(err)
}

func testConfigNewZapConfigBadErrorOutputPath(t *testing.T) {
	var (
		assert = assert.New(t)
		config = Config{
			Rotation:         new(Rotation),
			ErrorOutputPaths: []string{"#%@(&%(@%XX"},
		}
	)

	_, err := config.NewZapConfig()
	assert.Error(err)
}

func testConfigNewLoggerSuccess(t *testing.T) {
	var (
		assert = assert.New(t)
		file   = filepath.Join(os.TempDir(), "sallust-test.json")
		config = Config{
			Rotation: &Rotation{
				MaxSize: 100,
			},
			OutputPaths: []string{file},
		}
	)

	defer os.Remove(file)
	l, err := config.NewLogger()
	assert.NoError(err)
	assert.NotNil(l)
}

func testConfigNewLoggerBadOutputPath(t *testing.T) {
	var (
		assert = assert.New(t)
		file   = filepath.Join(os.TempDir(), "#^@*&^$*%XX")
		config = Config{
			Rotation: &Rotation{
				MaxSize: 100,
			},
			OutputPaths: []string{file},
		}
	)

	defer os.Remove(file)
	l, err := config.NewLogger()
	assert.Error(err)
	assert.Nil(l)
}

func TestConfig(t *testing.T) {
	t.Run("NewZapConfig", func(t *testing.T) {
		t.Run("Success", testConfigNewZapConfigSuccess)
		t.Run("BadOutputPath", testConfigNewZapConfigBadOutputPath)
		t.Run("BadErrorOutputPath", testConfigNewZapConfigBadErrorOutputPath)
	})

	t.Run("NewLogger", func(t *testing.T) {
		t.Run("Success", testConfigNewLoggerSuccess)
		t.Run("BadOutputPath", testConfigNewLoggerBadOutputPath)
	})
}
