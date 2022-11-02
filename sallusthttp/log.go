package sallusthttp

import (
	"net"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewConnStateLogger(logger *zap.Logger, lvl zapcore.Level, fs ...zap.Field) func(net.Conn, http.ConnState) {
	return func(c net.Conn, cs http.ConnState) {
		fs = append(fs, zap.String("connState", cs.String()))
		if addr := c.LocalAddr(); addr != nil {
			fs = append(fs, zap.String("localAddress", addr.String()))
		}

		logger.Log(
			lvl,
			"connState",
			fs...,
		)
	}
}
