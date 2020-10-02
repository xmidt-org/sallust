package sallust

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestStringToLevelHookFunc(t *testing.T) {
	expectedZapcoreLevel := zapcore.DebugLevel
	expectedZapLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)

	testData := []struct {
		from  reflect.Type
		to    reflect.Type
		value interface{}

		expected    interface{}
		expectedErr error
	}{
		{
			from:     reflect.TypeOf(int(0)),
			to:       reflect.TypeOf(""),
			value:    123,
			expected: 123,
		},
		{
			from:     reflect.TypeOf(""),
			to:       reflect.TypeOf(float64(0.0)),
			value:    "test",
			expected: "test",
		},
		{
			from:     reflect.TypeOf(""),
			to:       reflect.TypeOf(zapcore.Level(0)),
			value:    "debug",
			expected: zapcore.DebugLevel,
		},
		{
			from:     reflect.TypeOf(""),
			to:       reflect.PtrTo(reflect.TypeOf(zapcore.Level(0))),
			value:    "debug",
			expected: &expectedZapcoreLevel,
		},
		{
			from:     reflect.TypeOf(""),
			to:       reflect.TypeOf(zap.AtomicLevel{}),
			value:    "debug",
			expected: expectedZapLevel,
		},
		{
			from:     reflect.TypeOf(""),
			to:       reflect.PtrTo(reflect.TypeOf(zap.AtomicLevel{})),
			value:    "debug",
			expected: &expectedZapLevel,
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert            = assert.New(t)
				actual, actualErr = StringToLevelHookFunc(record.from, record.to, record.value)
			)

			assert.Equal(record.expected, actual)
			assert.Equal(record.expectedErr, actualErr)
		})
	}
}
