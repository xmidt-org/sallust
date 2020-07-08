package sallust

import (
	"sync"

	"go.uber.org/zap/zapcore"
)

// message is a captured log message
type message struct {
	entry  zapcore.Entry
	fields []zapcore.Field
}

// CaptureCore is a zapcore.Core which captures each log entry and makes
// it available for inspection in it's structured form.  This type allows
// tests to examine log output programmatically in a way that is much easier
// than simply parsing text.
type CaptureCore struct {
	// Core is the zap.Core to which output is delegated
	zapcore.Core

	lock     sync.RWMutex
	messages []message
	with     []zapcore.Field
}

// Check delegates to the embedded Core, then adds this CaptureCore
// if appropriate.
func (cc *CaptureCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	// this basic approach is at: https://github.com/uber-go/zap/blob/v1.15.0/zapcore/core.go#L78
	if cc.Core.Enabled(e.Level) {
		// add this decorated Core, rather than the delegate,
		// so the decorated Write method gets called
		return ce.AddCore(e, cc)
	}

	return ce
}

// With first delegates to the embedded Core, then returns a CaptureCore
// which decorates the new Core.  The fields are captured, and are added
// to every captured log message.
func (cc *CaptureCore) With(f []zapcore.Field) zapcore.Core {
	c := cc.Core.With(f)

	cc.lock.RLock()
	decorated := &CaptureCore{
		Core:     c,
		messages: append([]message{}, cc.messages...),
		with: append(
			append([]zapcore.Field{}, cc.with...),
			f...,
		),
	}

	cc.lock.RUnlock()
	return decorated
}

// Write captures the log message, then delegates to its embedded Core
// for output
func (cc *CaptureCore) Write(e zapcore.Entry, f []zapcore.Field) error {
	cc.lock.Lock()
	cc.messages = append(cc.messages,
		message{
			entry: e,
			fields: append(
				append([]zapcore.Field{}, f...),
				cc.with...,
			),
		},
	)
	cc.lock.Unlock()

	return cc.Core.Write(e, f)
}

// ClearMessages wipes out the captured log messages
func (cc *CaptureCore) ClearMessages() {
	cc.lock.Lock()
	cc.messages = nil
	cc.lock.Unlock()
}

// EachMessage applies a visitor closure to each capture Message.
// This method returns the number of messages visited.  Any fields
// added via With are included in the fields for each entry.
//
// If the supplied closure returns an error, visitation is halted
// and that error is returned.
func (cc *CaptureCore) EachMessage(f func(zapcore.Entry, []zapcore.Field) error) (int, error) {
	defer cc.lock.RUnlock()
	cc.lock.RLock()

	for i, m := range cc.messages {
		if err := f(m.entry, m.fields); err != nil {
			return i, err
		}
	}

	return len(cc.messages), nil
}

// Len returns the current count of captured messages
func (cc *CaptureCore) Len() (n int) {
	cc.lock.RLock()
	n = len(cc.messages)
	cc.lock.RUnlock()
	return
}

// Capture decorates a given core and captures messages written to it.
// This function may be used with zap.WrapCore.
//
//   var l *zap.Logger = ...
//   l = l.WithOptions(sallust.Capture)
//   cc := l.Core().(*sallust.CaptureCore)
func Capture(c zapcore.Core) zapcore.Core {
	return &CaptureCore{
		Core: c,
	}
}

// NewCaptureCore is an analog to zapcore.NewCore.  It produces a CaptureCore
// that captures log messages and delegates to a Core with the given configuration.
// This can be used directly to create a logger:
//
//   l := zap.New(sallust.NewCaptureCore(...))
//   cc := l.Core().(*sallust.CaptureCore)
func NewCaptureCore(enc zapcore.Encoder, ws zapcore.WriteSyncer, enab zapcore.LevelEnabler) *CaptureCore {
	return &CaptureCore{
		Core: zapcore.NewCore(enc, ws, enab),
	}
}
