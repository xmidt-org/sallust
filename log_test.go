package sallust

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestNewServerLogger(t *testing.T) {
	testLog := "foobar"
	sn := "serverName"
	require := require.New(t)
	assert := assert.New(t)
	verify, logger := NewTestLogger(zapcore.InfoLevel)
	l := NewServerLogger(sn, logger)
	require.NotNil(l)
	l.Print(testLog)
	vstring := verify.String()
	for _, tlog := range []string{sn, testLog} {
		assert.Contains(vstring, tlog)
	}
}
