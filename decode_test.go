package sallust

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func testStringToLevelHookNoConversion(t *testing.T) {
	var (
		assert      = assert.New(t)
		result, err = StringToLevelHook(
			reflect.TypeOf(int(0)),
			reflect.TypeOf(""),
			123,
		)
	)

	assert.Equal(123, result)
	assert.NoError(err)
}

func testStringToLevelHookToZapcoreLevel(t *testing.T) {
	var (
		assert   = assert.New(t)
		expected = zapcore.DebugLevel
	)

	result, err := StringToLevelHook(
		reflect.TypeOf(""),
		reflect.TypeOf(zapcore.Level(0)),
		"debug",
	)

	assert.Equal(expected, result)
	assert.NoError(err)

	result, err = StringToLevelHook(
		reflect.TypeOf(""),
		reflect.TypeOf(new(zapcore.Level)),
		"debug",
	)

	assert.Equal(&expected, result)
	assert.NoError(err)
}

func testStringToLevelHookToZapAtomicLevel(t *testing.T) {
	var (
		assert   = assert.New(t)
		expected = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	)

	result, err := StringToLevelHook(
		reflect.TypeOf(""),
		reflect.TypeOf(zap.AtomicLevel{}),
		"debug",
	)

	assert.Equal(expected, result)
	assert.NoError(err)

	result, err = StringToLevelHook(
		reflect.TypeOf(""),
		reflect.TypeOf(new(zap.AtomicLevel)),
		"debug",
	)

	assert.Equal(&expected, result)
	assert.NoError(err)
}

func TestStringToLevelHook(t *testing.T) {
	t.Run("NoConversion", testStringToLevelHookNoConversion)
	t.Run("ToZapcoreLevel", testStringToLevelHookToZapcoreLevel)
	t.Run("ToZapAtomicLevel", testStringToLevelHookToZapAtomicLevel)
}

func testStringToLevelEncoderHookNoConversion(t *testing.T) {
	var (
		assert      = assert.New(t)
		result, err = StringToLevelEncoderHook(
			reflect.TypeOf(int(0)),
			reflect.TypeOf(""),
			123,
		)
	)

	assert.Equal(123, result)
	assert.NoError(err)
}

func testStringToLevelEncoderHookSuccess(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = StringToLevelEncoderHook(
			reflect.TypeOf(""),
			reflect.TypeOf(zapcore.LevelEncoder(nil)),
			"capital",
		)
	)

	assert.NoError(err)

	levelEncoder, ok := result.(zapcore.LevelEncoder)
	require.True(ok)

	var output bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			LevelKey:    "level",
			EncodeLevel: levelEncoder,
		}),
		zapcore.AddSync(&output),
		zapcore.DebugLevel,
	)

	core.Write(zapcore.Entry{
		Level: zapcore.InfoLevel,
	}, nil)

	assert.Contains(output.String(), "INFO")
}

func TestStringToLevelEncoderHook(t *testing.T) {
	t.Run("NoConversion", testStringToLevelEncoderHookNoConversion)
	t.Run("Success", testStringToLevelEncoderHookSuccess)
}

func testStringToTimeEncoderHookNoConversion(t *testing.T) {
	var (
		assert      = assert.New(t)
		result, err = StringToTimeEncoderHook(
			reflect.TypeOf(int(0)),
			reflect.TypeOf(""),
			123,
		)
	)

	assert.Equal(123, result)
	assert.NoError(err)
}

func testStringToTimeEncoderHookSuccess(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		now         = time.Now()
		result, err = StringToTimeEncoderHook(
			reflect.TypeOf(""),
			reflect.TypeOf(zapcore.TimeEncoder(nil)),
			"RFC3339",
		)
	)

	assert.NoError(err)

	timeEncoder, ok := result.(zapcore.TimeEncoder)
	require.True(ok)

	var output bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:    "ts",
			EncodeTime: timeEncoder,
		}),
		zapcore.AddSync(&output),
		zapcore.DebugLevel,
	)

	core.Write(zapcore.Entry{
		Time: now,
	}, nil)

	assert.Contains(output.String(), now.Format(time.RFC3339))
}

func TestStringToTimeEncoderHook(t *testing.T) {
	t.Run("NoConversion", testStringToTimeEncoderHookNoConversion)
	t.Run("Success", testStringToTimeEncoderHookSuccess)
}

func testStringToDurationEncoderHookNoConversion(t *testing.T) {
	var (
		assert      = assert.New(t)
		result, err = StringToDurationEncoderHook(
			reflect.TypeOf(int(0)),
			reflect.TypeOf(""),
			123,
		)
	)

	assert.Equal(123, result)
	assert.NoError(err)
}

func testStringToDurationEncoderHookSuccess(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = StringToDurationEncoderHook(
			reflect.TypeOf(""),
			reflect.TypeOf(zapcore.DurationEncoder(nil)),
			"string",
		)
	)

	assert.NoError(err)

	durationEncoder, ok := result.(zapcore.DurationEncoder)
	require.True(ok)

	var output bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			EncodeDuration: durationEncoder,
		}),
		zapcore.AddSync(&output),
		zapcore.DebugLevel,
	)

	core.Write(zapcore.Entry{}, []zapcore.Field{
		{
			Key:     "test",
			Type:    zapcore.DurationType,
			Integer: int64(17 * time.Minute),
		},
	})

	assert.Contains(output.String(), "17m")
}

func TestStringToDurationEncoderHook(t *testing.T) {
	t.Run("NoConversion", testStringToDurationEncoderHookNoConversion)
	t.Run("Success", testStringToDurationEncoderHookSuccess)
}

func testStringToCallerEncoderHookNoConversion(t *testing.T) {
	var (
		assert      = assert.New(t)
		result, err = StringToCallerEncoderHook(
			reflect.TypeOf(int(0)),
			reflect.TypeOf(""),
			123,
		)
	)

	assert.Equal(123, result)
	assert.NoError(err)
}

func testStringToCallerEncoderHookSuccess(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = StringToCallerEncoderHook(
			reflect.TypeOf(""),
			reflect.TypeOf(zapcore.CallerEncoder(nil)),
			"short",
		)
	)

	assert.NoError(err)

	callerEncoder, ok := result.(zapcore.CallerEncoder)
	require.True(ok)

	var output bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			CallerKey:    "test",
			EncodeCaller: callerEncoder,
		}),
		zapcore.AddSync(&output),
		zapcore.DebugLevel,
	)

	core.Write(zapcore.Entry{
		Caller: zapcore.EntryCaller{
			Defined: true,
			File:    "foo/bar.go",
			Line:    123,
		},
	}, nil)

	assert.Contains(output.String(), "foo/bar.go:123")
}

func TestStringToCallerEncoderHook(t *testing.T) {
	t.Run("NoConversion", testStringToCallerEncoderHookNoConversion)
	t.Run("Success", testStringToCallerEncoderHookSuccess)
}

func testStringToNameEncoderHookNoConversion(t *testing.T) {
	var (
		assert      = assert.New(t)
		result, err = StringToNameEncoderHook(
			reflect.TypeOf(int(0)),
			reflect.TypeOf(""),
			123,
		)
	)

	assert.Equal(123, result)
	assert.NoError(err)
}

func testStringToNameEncoderHookSuccess(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = StringToNameEncoderHook(
			reflect.TypeOf(""),
			reflect.TypeOf(zapcore.NameEncoder(nil)),
			"full",
		)
	)

	assert.NoError(err)

	nameEncoder, ok := result.(zapcore.NameEncoder)
	require.True(ok)

	var output bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			NameKey:    "test",
			EncodeName: nameEncoder,
		}),
		zapcore.AddSync(&output),
		zapcore.DebugLevel,
	)

	core.Write(zapcore.Entry{
		LoggerName: "foo.bar",
	}, nil)

	assert.Contains(output.String(), "foo.bar")
}

func TestStringToNameEncoderHook(t *testing.T) {
	t.Run("NoConversion", testStringToNameEncoderHookNoConversion)
	t.Run("Success", testStringToNameEncoderHookSuccess)
}
