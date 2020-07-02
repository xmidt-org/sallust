package sallust

import (
	"net/url"
	"strconv"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// LumberjackScheme is the URL Scheme for lumberjack-rotatable files
	LumberjackScheme = "lumberjack"

	// MaxSizeParameter is the URL parameter that corresponds to lumberjack.Logger.MaxSize
	MaxSizeParameter = "maxSize"

	// MaxAgeParameter is the URL parameter that corresponds to lumberjack.Logger.MaxAge
	MaxAgeParameter = "maxAge"

	// MaxBackupsParameter is the URL parameter that corresponds to lumberjack.Logger.MaxBackups
	MaxBackupsParameter = "maxBackups"

	// LocalTimeParameter is the URL parameter that corresponds to lumberjack.Logger.LocalTime
	LocalTimeParameter = "localTime"

	// CompressParameter is the URL parameter that corresponds to lumberjack.Logger.Compress
	CompressParameter = "compress"
)

func init() {
	zap.RegisterSink(LumberjackScheme, NewLumberjackSink)
}

// Rotater is implemented by objects which can rotate logs
type Rotater interface {
	Rotate() error
}

// Rotation describes the set of configurable options for log file rotation.
// This configuration, if supplied, is only applied to file outputs.
//
// The fields in this struct correspond exactly to lumberjack.Logger.
//
// See: https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2?tab=doc#Logger
type Rotation struct {
	// MaxSize corresponds to lumberjack.Logger.MaxSize
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	// MaxAge corresponds to lumberjack.Logger.MaxAge
	MaxAge int `json:"maxage" yaml:"maxage"`

	// MaxBackups corresponds to lumberjack.Logger.MaxBackups
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// LocalTime corresponds to lumberjack.Logger.LocalTime
	LocalTime bool `json:"localtime" yaml:"localtime"`

	// Compress corresponds to lumberjack.Logger.Compress
	Compress bool `json:"compress" yaml:"compress"`
}

// AddQuery adds the set of URL query parameters for these Rotation options
func (r Rotation) AddQuery(v url.Values) {
	if r.MaxSize > 0 {
		v.Set(MaxSizeParameter, strconv.Itoa(r.MaxSize))
	}

	if r.MaxAge > 0 {
		v.Set(MaxAgeParameter, strconv.Itoa(r.MaxAge))
	}

	if r.MaxBackups > 0 {
		v.Set(MaxBackupsParameter, strconv.Itoa(r.MaxBackups))
	}

	if r.LocalTime {
		v.Set(LocalTimeParameter, strconv.FormatBool(r.LocalTime))
	}

	if r.Compress {
		v.Set(CompressParameter, strconv.FormatBool(r.Compress))
	}
}

// NewURL creates a URL object that represents a lumberjack-rotatable file
// using this set of Rotation options
func (r Rotation) NewURL(path string) *url.URL {
	v := url.Values{}
	r.AddQuery(v)

	return &url.URL{
		Scheme:   LumberjackScheme,
		Path:     path,
		RawQuery: v.Encode(),
	}
}

// Lumberjack is a zap.Sink that writes to a lumberjack writer.
// This type also implements Rotater and io.Closer.
//
// A Lumberjack is safe for concurrent writes.  No additional synchronization
// is required.
type Lumberjack struct {
	*lumberjack.Logger
}

var _ zap.Sink = Lumberjack{}

// Sync is a nop, and implements zapcore.WriteSyncer
func (lj Lumberjack) Sync() error {
	return nil
}

// NewLumberjackSink creates a zap.Sink which rotates its corresponding file
func NewLumberjackSink(u *url.URL) (zap.Sink, error) {
	lj := &Lumberjack{
		Logger: &lumberjack.Logger{
			Filename: u.Path,
		},
	}

	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	if v := values.Get(MaxSizeParameter); len(v) > 0 {
		lj.MaxSize, err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	if v := values.Get(MaxAgeParameter); len(v) > 0 {
		lj.MaxAge, err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	if v := values.Get(MaxBackupsParameter); len(v) > 0 {
		lj.MaxBackups, err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	if v := values.Get(LocalTimeParameter); len(v) > 0 {
		lj.LocalTime, err = strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
	}

	if v := values.Get(CompressParameter); len(v) > 0 {
		lj.Compress, err = strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}
	}

	return lj, nil
}
