package sallust

import (
	"log"

	"go.uber.org/zap"
)

// NewErrorLog creates a new logging.Logger appropriate for http.Server.ErrorLog
func NewErrorLog(serverName string, logger *zap.Logger) *log.Logger {
	return log.New(
		zap.NewStdLog(logger).Writer(),
		serverName,
		log.LstdFlags|log.LUTC,
	)
}
