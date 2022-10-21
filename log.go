package sallust

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A Valuer generates a log value. When passed to With
// in a value element (odd indexes), it represents a dynamic
// value which is re-evaluated with each log event.
type Valuer func() interface{}

// LoggerFunc is a strategy for adding key/value pairs (possibly) based on an HTTP request.
// Functions of this type must append key/value pairs to the supplied slice and then return
// the new slice.
type LoggerFunc func([]zap.Field, *http.Request) []zap.Field

var (

	// DefaultCaller is a Valuer that returns the file and line where the Log
	// method was invoked. It can only be used with log.With.
	DefaultCaller = Caller(3)
)

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

// Caller returns a Valuer that returns a file and line from a specified depth
// in the callstack. Users will probably want to use DefaultCaller.
func Caller(depth int) Valuer {
	return func() interface{} {
		_, file, line, _ := runtime.Caller(depth)
		idx := strings.LastIndexByte(file, '/')
		// using idx+1 below handles both of following cases:
		// idx == -1 because no "/" was found, or
		// idx >= 0 and we want to start at the character after the found "/".
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}
