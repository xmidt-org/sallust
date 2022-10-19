package sallust

import (
	"go.uber.org/zap"
)

// Contextual describes an object which can describe itself with metadata for logging.
// Implementing this interface allows code to carry logging context data across API
// boundaries without compromising encapsulation.
type Contextual interface {
	Metadata() map[string]interface{}
}

// enrich is the helper function that emits contextual information into its logger argument.
func enrich(wither func(fields ...zap.Field) *zap.Logger, logger *zap.Logger, objects []interface{}) *zap.Logger {
	var kvs []zap.Field
	for _, e := range objects {
		switch m := e.(type) {
		case Contextual:
			for k, v := range m.Metadata() {
				kvs = append(kvs, zap.Any(k, v))
			}

		case map[string]interface{}:
			for k, v := range m {
				kvs = append(kvs, zap.Any(k, v))
			}

		case map[string]string:
			for k, v := range m {
				kvs = append(kvs, zap.Any(k, v))
			}
		}
	}

	if len(kvs) > 0 {
		return wither(kvs...)
	}

	return logger
}

// Enrich uses zap.Logger.With to add contextual information to a logger.  The given set of objects are examined to see if they contain
// any metadata.  Objects that do not contain metadata are simply ignored.
//
// An object contains metadata if it implements Contextual, is a map[string]interface{}, or is a map[string]string.  In those cases,
// the key/value pairs are present in the returned logger.
func Enrich(logger *zap.Logger, objects ...interface{}) *zap.Logger {
	return enrich(logger.With, logger, objects)
}
