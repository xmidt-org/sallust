package sallusthttp

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
)

func TestWith(t *testing.T) {
	var (
		assert   = assert.New(t)
		require  = require.New(t)
		logger   = new(zap.Logger)
		original = httptest.NewRequest("PUT", "/test", nil)
	)

	child := With(original, logger)
	require.NotNil(child)
	assert.Equal(original.Method, child.Method)
	assert.Equal(original.URL.String(), child.URL.String())
	ctx := child.Context()
	assert.Equal(logger, sallust.Get(ctx))
}

func TestGet(t *testing.T) {
	var (
		assert   = assert.New(t)
		original = httptest.NewRequest("GET", "/test", nil)
	)

	logger := Get(original)
	assert.Equal(sallust.Default(), logger)

	def := new(zap.Logger)
	logger = GetDefault(original, def)
	assert.Equal(def, logger)

	logger = new(zap.Logger)
	child := With(original, logger)
	assert.Equal(logger, Get(child))
}
