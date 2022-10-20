package sallust

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newTestLogger(t *testing.T) (*bytes.Buffer, *zap.Logger) {
	b := &bytes.Buffer{}
	return b, zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zap.CombineWriteSyncers(os.Stderr, zapcore.AddSync(b)),
		zapcore.InfoLevel,
	))
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
	vstring := verify.String()
	for _, tlog := range []string{sn, testLog} {
		assert.Contains(vstring, tlog)
	}
}
