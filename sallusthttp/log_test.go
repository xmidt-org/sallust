package sallusthttp

import (
	"bytes"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
)

func assertBufferContains(assert *assert.Assertions, verify *bytes.Buffer, values ...string) {
	text := verify.String()
	for _, value := range values {
		assert.Contains(text, value)
	}
}

func assertConnState(assert *assert.Assertions, verify *bytes.Buffer, connState func(net.Conn, http.ConnState)) {
	if assert.NotNil(connState) {
		conn1, conn2 := net.Pipe()
		defer conn1.Close()
		defer conn2.Close()

		assert.NotPanics(func() {
			connState(conn1, http.StateNew)
		})
		assert.NotPanics(func() {
			connState(conn1, http.StateNew)
		})
		if verify != nil {
			assertBufferContains(assert, verify, conn1.LocalAddr().String(), http.StateNew.String())
		}
	}
}

func TestNewConnStateLogger(t *testing.T) {
	var (
		assert    = assert.New(t)
		require   = require.New(t)
		v, l      = sallust.NewTestLogger(zap.DebugLevel)
		connState = NewConnStateLogger(l, "serverName", zap.DebugLevel)
	)

	require.NotNil(connState)
	assertConnState(assert, v, connState)
}
