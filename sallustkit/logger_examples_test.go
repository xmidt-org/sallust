package sallustkit

import (
	"time"

	"github.com/go-kit/log/level"
	"go.uber.org/zap"
)

func ExampleLogger() {
	l := Logger{
		Zap: zap.NewExample(),
	}

	l.Log("msg", "hello, world!", "duration", 15*time.Second)
	l.Log("msg", "this entry supplies a level", "value", 123, "this key doesn't matter", level.DebugValue())

	// Output:
	// {"level":"error","msg":"hello, world!","duration":"15s"}
	// {"level":"debug","msg":"this entry supplies a level","value":123}
}

func ExampleLogger_custom() {
	l := Logger{
		Zap:          zap.NewExample(),
		MessageKey:   "message",
		DefaultLevel: level.InfoValue(),
	}

	// the message will still be output to whatever zap is configured with
	l.Log("message", "hello, world!", "duration", 15*time.Second)

	// override the DefaultLevel with what's passed to Log
	l.Log("message", "this entry supplies a level", "value", 123, "this key doesn't matter", level.DebugValue())

	// Output:
	// {"level":"info","msg":"hello, world!","duration":"15s"}
	// {"level":"debug","msg":"this entry supplies a level","value":123}
}
