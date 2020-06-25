package sallust

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

// Rotater is implemented by objects which can rotate logs
type Rotater interface {
	Rotate() error
}

// Lumberjack is a zapcore.WriteSyncer that writes to a lumberjack writer.
// This type also implements Rotater and io.Closer.
//
// A Lumberjack is safe for concurrent writes.  No additional synchronization
// is required.
type Lumberjack struct {
	w *lumberjack.Logger
}

// Write implements io.Writer
func (lj *Lumberjack) Write(p []byte) (int, error) {
	return lj.w.Write(p)
}

// Sync is a nop, and implements zapcore.WriteSyncer
func (lj *Lumberjack) Sync() error {
	return nil
}

// Close closes the underlying lumberjack logger.  Note that subsequent
// writes will cause any logfile to be reopened and possibly rotated.
func (lj *Lumberjack) Close() error {
	return lj.w.Close()
}

// Rotate implements Rotater.  It rotates the underlying lumberjack logger.
func (lj *Lumberjack) Rotate() error {
	return lj.w.Rotate()
}

// NewLumberjack returns a Lumberjack which wraps the given logger.  The
// returned instances may be used as a zapcore.WriteSyncer.
//
// Once passed to this function, the supplied lumberjack Logger must not
// be modified.
func NewLumberjack(w *lumberjack.Logger) *Lumberjack {
	return &Lumberjack{
		w: w,
	}
}
