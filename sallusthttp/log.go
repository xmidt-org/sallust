package sallusthttp

import (
	"net"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewConnStateLogger(logger *zap.Logger, serverName string, lvl zapcore.Level) func(net.Conn, http.ConnState) {
	return func(c net.Conn, cs http.ConnState) {
		logger.Log(
			lvl,
			"connState",
			zap.String("serverName", serverName),
			zap.String("localAddress", c.LocalAddr().String()),
			zap.String("state", cs.String()),
		)
	}
}
