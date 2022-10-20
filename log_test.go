package sallust

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func newTestLogger(t *testing.T) (*zapcore.Entry, *zap.Logger) {
	var verify zapcore.Entry
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
		func(e zapcore.Entry) error {
			verify = e
			return nil
		})))
	return &verify, logger
}

func TestNewErrorLog(t *testing.T) {
	testLog := "foobar"
	sn := "serverName"
	require := require.New(t)
	assert := assert.New(t)
	verify, logger := newTestLogger(t)
	l := NewErrorLog(sn, logger)
	require.NotNil(l)
	l.Print(testLog)

	for _, tlog := range []string{sn, testLog} {
		assert.Contains(verify.Message, tlog)
	}
}
