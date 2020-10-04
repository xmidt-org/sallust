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
			config: Config{},
			expected: zap.Config{
				Level:            zap.NewAtomicLevelAt(zapcore.ErrorLevel),
				Encoding:         "json",
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
			},
		},
		{
			config: Config{
				Config: zap.Config{
					EncoderConfig: zapcore.EncoderConfig{
						MessageKey:       "msg",
						LevelKey:         "level",
						TimeKey:          "ts",
						NameKey:          "name",
						CallerKey:        "caller",
						FunctionKey:      "function",
						StacktraceKey:    "trace",
						LineEnding:       "--",
						ConsoleSeparator: ".",
					},
					OutputPaths: []string{"/log.json"},
				},
			},
			expected: zap.Config{
				Level:            zap.NewAtomicLevelAt(zapcore.ErrorLevel),
				Encoding:         "json",
				OutputPaths:      []string{"/log.json"},
				ErrorOutputPaths: []string{"stderr"},
				EncoderConfig: zapcore.EncoderConfig{
					MessageKey:       "msg",
					LevelKey:         "level",
					TimeKey:          "ts",
					NameKey:          "name",
					CallerKey:        "caller",
					FunctionKey:      "function",
					StacktraceKey:    "trace",
					LineEnding:       "--",
					ConsoleSeparator: ".",
					// function types omitted
				},
			},
		},
		{
			config: Config{
				Config: zap.Config{
					OutputPaths: []string{"/log.json"},
				},
				Rotation: &Rotation{
					MaxAge: 10,
				},
			},
			expected: zap.Config{
				Level:            zap.NewAtomicLevelAt(zapcore.ErrorLevel),
				Encoding:         "json",
				OutputPaths:      []string{"lumberjack:///log.json?maxAge=10"},
				ErrorOutputPaths: []string{"stderr"},
				// function types in EncoderConfig omitted
			},
		},
		{
			config: Config{
				Config: zap.Config{
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
				EncoderConfig:     zapcore.EncoderConfig{ /* function types omitted */ },
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
			assert.Equal(record.expected.OutputPaths, actual.OutputPaths)
			assert.Equal(record.expected.ErrorOutputPaths, actual.ErrorOutputPaths)
			assert.Equal(record.expected.InitialFields, actual.InitialFields)

			assert.Equal(record.expected.EncoderConfig.MessageKey, actual.EncoderConfig.MessageKey)
			assert.Equal(record.expected.EncoderConfig.LevelKey, actual.EncoderConfig.LevelKey)
			assert.Equal(record.expected.EncoderConfig.TimeKey, actual.EncoderConfig.TimeKey)
			assert.Equal(record.expected.EncoderConfig.NameKey, actual.EncoderConfig.NameKey)
			assert.Equal(record.expected.EncoderConfig.CallerKey, actual.EncoderConfig.CallerKey)
			assert.Equal(record.expected.EncoderConfig.FunctionKey, actual.EncoderConfig.FunctionKey)
			assert.Equal(record.expected.EncoderConfig.StacktraceKey, actual.EncoderConfig.StacktraceKey)
			assert.Equal(record.expected.EncoderConfig.LineEnding, actual.EncoderConfig.LineEnding)
			assert.Equal(record.expected.EncoderConfig.ConsoleSeparator, actual.EncoderConfig.ConsoleSeparator)

			assert.NotNil(actual.EncoderConfig.EncodeLevel)
			assert.NotNil(actual.EncoderConfig.EncodeTime)
			assert.NotNil(actual.EncoderConfig.EncodeDuration)
			assert.NotNil(actual.EncoderConfig.EncodeCaller)
			assert.NotNil(actual.EncoderConfig.EncodeName)
		})
	}
}

func testConfigNewZapConfigBadOutputPath(t *testing.T) {
	var (
		assert = assert.New(t)
		config = Config{
			Rotation: new(Rotation),
			Config: zap.Config{
				OutputPaths: []string{"#%@(&%(@%XX"},
			},
		}
	)

	_, err := config.NewZapConfig()
	assert.Error(err)
}

func testConfigNewZapConfigBadErrorOutputPath(t *testing.T) {
	var (
		assert = assert.New(t)
		config = Config{
			Rotation: new(Rotation),
			Config: zap.Config{
				ErrorOutputPaths: []string{"#%@(&%(@%XX"},
			},
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
			Config: zap.Config{
				OutputPaths: []string{file},
			},
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
			Config: zap.Config{
				OutputPaths: []string{file},
			},
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
