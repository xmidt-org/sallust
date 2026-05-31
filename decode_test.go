// SPDX-FileCopyrightText: 2024 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
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

func testDecodeHookNotAString(t *testing.T) {
	var (
		assert      = assert.New(t)
		result, err = DecodeHook(
			reflect.TypeFor[int](),
			reflect.TypeFor[string](),
			123,
		)
	)

	assert.Equal(123, result)
	assert.NoError(err)
}

func testDecodeHookUnsupported(t *testing.T) {
	var (
		assert      = assert.New(t)
		result, err = DecodeHook(
			reflect.TypeFor[string](),
			reflect.TypeFor[float64](),
			"test",
		)
	)

	assert.Equal("test", result)
	assert.NoError(err)
}

func testDecodeHookToLevel(t *testing.T) {
	var (
		assert   = assert.New(t)
		expected = zapcore.DebugLevel
	)

	result, err := DecodeHook(
		reflect.TypeFor[string](),
		reflect.TypeFor[zapcore.Level](),
		"debug",
	)

	assert.Equal(expected, result)
	assert.NoError(err)

	result, err = DecodeHook(
		reflect.TypeFor[string](),
		reflect.TypeFor[*zapcore.Level](),
		"debug",
	)

	assert.Equal(&expected, result)
	assert.NoError(err)
}

func testDecodeHookToAtomicLevel(t *testing.T) {
	var (
		assert   = assert.New(t)
		expected = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	)

	result, err := DecodeHook(
		reflect.TypeFor[string](),
		reflect.TypeFor[zap.AtomicLevel](),
		"debug",
	)

	assert.Equal(expected, result)
	assert.NoError(err)

	result, err = DecodeHook(
		reflect.TypeFor[string](),
		reflect.TypeFor[*zap.AtomicLevel](),
		"debug",
	)

	assert.Equal(&expected, result)
	assert.NoError(err)
}

func testDecodeHookToLevelEncoder(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = DecodeHook(
			reflect.TypeFor[string](),
			reflect.TypeFor[zapcore.LevelEncoder](),
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

func testDecodeHookToTimeEncoder(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		now         = time.Now()
		result, err = DecodeHook(
			reflect.TypeFor[string](),
			reflect.TypeFor[zapcore.TimeEncoder](),
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

func testDecodeHookToDurationEncoder(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = DecodeHook(
			reflect.TypeFor[string](),
			reflect.TypeFor[zapcore.DurationEncoder](),
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

func testDecodeHookToCallerEncoder(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = DecodeHook(
			reflect.TypeFor[string](),
			reflect.TypeFor[zapcore.CallerEncoder](),
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

func testDecodeHookToNameEncoder(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = DecodeHook(
			reflect.TypeFor[string](),
			reflect.TypeFor[zapcore.NameEncoder](),
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

func TestDecodeHook(t *testing.T) {
	t.Run("NotAString", testDecodeHookNotAString)
	t.Run("Unsupported", testDecodeHookUnsupported)
	t.Run("ToLevel", testDecodeHookToLevel)
	t.Run("ToAtomicLevel", testDecodeHookToAtomicLevel)
	t.Run("ToLevelEncoder", testDecodeHookToLevelEncoder)
	t.Run("ToTimeEncoder", testDecodeHookToTimeEncoder)
	t.Run("ToDurationEncoder", testDecodeHookToDurationEncoder)
	t.Run("ToCallerEncoder", testDecodeHookToCallerEncoder)
	t.Run("ToNameEncoder", testDecodeHookToNameEncoder)
}
