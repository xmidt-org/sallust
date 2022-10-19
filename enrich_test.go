package sallust

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testContextual struct {
	key, value string
}

func (tc testContextual) Metadata() map[string]interface{} {
	return map[string]interface{}{tc.key: tc.value}
}

func testEnrichNoObjects(t *testing.T) {
	var (
		require  = require.New(t)
		delegate = new(mockLogger)
	)

	logger := Enrich(delegate)
	require.NotNil(logger)

	delegate.On("Log", []interface{}{"message", "foobar"}).Return(error(nil)).Once()
	logger.Log("message", "foobar")
	delegate.AssertExpectations(t)
}

func testEnrichWithObjects(t *testing.T) {
	var (
		require  = require.New(t)
		delegate = new(mockLogger)
	)

	logger := Enrich(delegate, map[string]string{"key1": "value1"}, nil, map[string]interface{}{"key2": "value2"}, 27, testContextual{"key3", "value3"})
	require.NotNil(logger)

	delegate.On("Log", []interface{}{"key1", "value1", "key2", "value2", "key3", "value3", "message", "foobar"}).Return(error(nil)).Once()
	logger.Log("message", "foobar")
	delegate.AssertExpectations(t)
}

func TestEnrich(t *testing.T) {
	t.Run("NoObjects", testEnrichNoObjects)
	t.Run("WithObjects", testEnrichWithObjects)
}
