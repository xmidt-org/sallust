package sallust

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

type testContextual struct {
	key, value string
}

func (tc testContextual) Metadata() map[string]interface{} {
	return map[string]interface{}{tc.key: tc.value}
}

func TestEnrich(t *testing.T) {
	tests := []struct {
		description string
		objects     []interface{}
		expectedErr error
	}{
		{
			description: "Validate enrich logging without objects",
		},
		{
			description: "Validate enrich logging with objects",
			objects:     []interface{}{"key1", "value1", "key2", "value2", "key3", "value3", "message", "foobar"},
		},
		{
			description: "Validate enrich logging with objects",
			objects:     []interface{}{[]interface{}{"key1", "value1", "message", "foobar"}, 27, testContextual{"key3", "value3"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)
			logCount := 0
			logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
				func(e zapcore.Entry) error {
					logCount++
					return nil
				})))
			logger = Enrich(logger, tc.objects...)

			require.NotNil(logger)
			logger.Info("foobar")
			assert.Equal(1, logCount)
		})
	}
}
