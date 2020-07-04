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

func testOptionsNewZapConfigSuccess(t *testing.T) {
	testData := []struct {
		options  Options
		expected zap.Config
	}{
		{
			options: Options{
				OutputPaths: []string{"/log.json"},
			},
			expected: zap.Config{
				Level:       zap.NewAtomicLevelAt(zapcore.ErrorLevel),
				Encoding:    "json",
				OutputPaths: []string{"/log.json"},
			},
		},
		{
			options: Options{
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
			options: Options{
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
				actual, err = record.options.NewZapConfig()
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

func testOptionsNewZapConfigBadOutputPath(t *testing.T) {
	var (
		assert  = assert.New(t)
		options = Options{
			Rotation:    new(Rotation),
			OutputPaths: []string{"#%@(&%(@%XX"},
		}
	)

	_, err := options.NewZapConfig()
	assert.Error(err)
}

func testOptionsNewZapConfigBadErrorOutputPath(t *testing.T) {
	var (
		assert  = assert.New(t)
		options = Options{
			Rotation:         new(Rotation),
			ErrorOutputPaths: []string{"#%@(&%(@%XX"},
		}
	)

	_, err := options.NewZapConfig()
	assert.Error(err)
}

func testOptionsNewLoggerSuccess(t *testing.T) {
	var (
		assert  = assert.New(t)
		file    = filepath.Join(os.TempDir(), "sallust-test.json")
		options = Options{
			Rotation: &Rotation{
				MaxSize: 100,
			},
			OutputPaths: []string{file},
		}
	)

	defer os.Remove(file)
	l, err := options.NewLogger()
	assert.NoError(err)
	assert.NotNil(l)
}

func testOptionsNewLoggerBadOutputPath(t *testing.T) {
	var (
		assert  = assert.New(t)
		file    = filepath.Join(os.TempDir(), "#^@*&^$*%XX")
		options = Options{
			Rotation: &Rotation{
				MaxSize: 100,
			},
			OutputPaths: []string{file},
		}
	)

	defer os.Remove(file)
	l, err := options.NewLogger()
	assert.Error(err)
	assert.Nil(l)
}

func TestOptions(t *testing.T) {
	t.Run("NewZapConfig", func(t *testing.T) {
		t.Run("Success", testOptionsNewZapConfigSuccess)
		t.Run("BadOutputPath", testOptionsNewZapConfigBadOutputPath)
		t.Run("BadErrorOutputPath", testOptionsNewZapConfigBadErrorOutputPath)
	})

	t.Run("NewLogger", func(t *testing.T) {
		t.Run("Success", testOptionsNewLoggerSuccess)
		t.Run("BadOutputPath", testOptionsNewLoggerBadOutputPath)
	})
}
