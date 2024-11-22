// SPDX-FileCopyrightText: 2024 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package sallust

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func testWithNil(t *testing.T) {
	var (
		assert    = assert.New(t)
		parentCtx = context.Background()
		ctx       = With(parentCtx, nil)
	)

	assert.Equal(parentCtx, ctx)
}

func testWithNonNil(t *testing.T) {
	var (
		assert    = assert.New(t)
		parentCtx = context.Background()
		ctx       = With(parentCtx, new(zap.Logger))
	)

	assert.NotEqual(parentCtx, ctx)
}

func TestWith(t *testing.T) {
	t.Run("Nil", testWithNil)
	t.Run("NonNil", testWithNonNil)
}

func testGetNoLogger(t *testing.T) {
	var (
		assert  = assert.New(t)
		require = require.New(t)
		actual  = Get(context.Background())
	)

	require.NotNil(actual)
	assert.Equal(Default(), actual)
}

func testGetWithLogger(t *testing.T) {
	var (
		assert   = assert.New(t)
		require  = require.New(t)
		expected = new(zap.Logger)
		actual   = Get(With(context.Background(), expected))
	)

	require.NotNil(actual)
	assert.Equal(expected, actual)
}

func TestGet(t *testing.T) {
	t.Run("NoLogger", testGetNoLogger)
	t.Run("WithLogger", testGetWithLogger)
}

func testGetDefaultNoLogger(t *testing.T) {
	t.Run("NoDefault", func(t *testing.T) {
		var (
			assert  = assert.New(t)
			require = require.New(t)
			actual  = GetDefault(context.Background(), nil)
		)

		require.NotNil(actual)
		assert.Equal(Default(), actual)
	})

	t.Run("WithDefault", func(t *testing.T) {
		var (
			assert   = assert.New(t)
			require  = require.New(t)
			expected = new(zap.Logger)
			actual   = GetDefault(context.Background(), expected)
		)

		require.NotNil(actual)
		assert.Equal(expected, actual)
	})
}

func testGetDefaultWithLogger(t *testing.T) {
	var (
		assert   = assert.New(t)
		require  = require.New(t)
		expected = new(zap.Logger)
		actual   = GetDefault(With(context.Background(), expected), new(zap.Logger))
	)

	require.NotNil(actual)
	assert.Equal(expected, actual)
}

func TestGetDefault(t *testing.T) {
	t.Run("NoLogger", testGetDefaultNoLogger)
	t.Run("WithLogger", testGetDefaultWithLogger)
}
