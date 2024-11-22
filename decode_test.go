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
			reflect.TypeOf(int(0)),
			reflect.TypeOf(""),
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
			reflect.TypeOf(""),
			reflect.TypeOf(float64(8.9)),
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
		reflect.TypeOf(""),
		reflect.TypeOf(zapcore.Level(0)),
		"debug",
	)

	assert.Equal(expected, result)
	assert.NoError(err)

	result, err = DecodeHook(
		reflect.TypeOf(""),
		reflect.TypeOf(new(zapcore.Level)),
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
		reflect.TypeOf(""),
		reflect.TypeOf(zap.AtomicLevel{}),
		"debug",
	)

	assert.Equal(expected, result)
	assert.NoError(err)

	result, err = DecodeHook(
		reflect.TypeOf(""),
		reflect.TypeOf(new(zap.AtomicLevel)),
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

func testDecodeHookToTimeEncoder(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		now         = time.Now()
		result, err = DecodeHook(
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

func testDecodeHookToDurationEncoder(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = DecodeHook(
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

func testDecodeHookToCallerEncoder(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = DecodeHook(
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

func testDecodeHookToNameEncoder(t *testing.T) {
	var (
		assert      = assert.New(t)
		require     = require.New(t)
		result, err = DecodeHook(
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
