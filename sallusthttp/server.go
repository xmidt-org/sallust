package sallusthttp

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

// NewConnStateLogger produces an http.Server.ConnState closure that logs the connection
// state to the supplied zap logger function.
//
//   l := zap.NewDevelopment()
//   server.ConnState = NewConnStateLogger(l.Debug)
func NewConnStateLogger(l func(string, ...zap.Field)) func(net.Conn, http.ConnState) {
	return func(c net.Conn, cs http.ConnState) {
		l("connection state changed", zap.Stringer("remoteAddr", c.RemoteAddr()), zap.Stringer("state", cs))
	}
}
