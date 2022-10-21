package sallust

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerFunc is a strategy for adding key/value pairs (possibly) based on an HTTP request.
// Functions of this type must append key/value pairs to the supplied slice and then return
// the new slice.
type LoggerFunc func([]zap.Field, *http.Request) []zap.Field

func NewTestLogger() (*bytes.Buffer, *zap.Logger) {
	b := &bytes.Buffer{}
	return b, zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zap.CombineWriteSyncers(os.Stderr, zapcore.AddSync(b)),
		zapcore.InfoLevel,
	))
}

// NewServerLogger creates a new zap.Logger appropriate for http.Server.ErrorLog
func NewServerLogger(serverName string, logger *zap.Logger) *log.Logger {
	if logger == nil {
		logger = Default()
	}

	return log.New(
		zap.NewStdLog(logger).Writer(),
		serverName,
		log.LstdFlags|log.LUTC,
	)
}
