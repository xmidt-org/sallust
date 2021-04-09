package sallustkit

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerTestSuite struct {
	GokitTestSuite
}

var _ suite.SetupTestSuite = (*LoggerTestSuite)(nil)

func (suite *LoggerTestSuite) TestNoSuppliedLevel() {
	suite.Run("NoDefaultLevel", func() {
		suite.core.ClearMessages()
		l := Logger{
			Zap: suite.zapLogger,
		}

		l.Log(DefaultMessageKey, "test message", "value", 123)
		suite.assertOneMessage(func(e zapcore.Entry, fs []zap.Field) error {
			suite.Equal("test message", e.Message)
			suite.Equal(zap.ErrorLevel, e.Level)
			suite.Equal(
				[]zap.Field{zap.Int("value", 123)},
				fs,
			)

			return nil
		})
	})

	suite.Run("WithDefaultLevel", func() {
		suite.core.ClearMessages()
		l := Logger{
			Zap:          suite.zapLogger,
			DefaultLevel: level.InfoValue(),
		}

		l.Log(DefaultMessageKey, "test message", "value", 123)
		suite.assertOneMessage(func(e zapcore.Entry, fs []zap.Field) error {
			suite.Equal("test message", e.Message)
			suite.Equal(zap.InfoLevel, e.Level)
			suite.Equal(
				[]zap.Field{zap.Int("value", 123)},
				fs,
			)

			return nil
		})
	})
}

func (suite *LoggerTestSuite) TestLevelSupplied() {
	l := Logger{
		Zap:          suite.zapLogger,
		DefaultLevel: level.WarnValue(), // verify that this is overwritten
	}

	l.Log("value", 123, DefaultMessageKey, "test message", "doesn't matter", level.InfoValue())
	suite.assertOneMessage(func(e zapcore.Entry, fs []zap.Field) error {
		suite.Equal("test message", e.Message)
		suite.Equal(zap.InfoLevel, e.Level)
		suite.Equal(
			[]zap.Field{zap.Int("value", 123)},
			fs,
		)

		return nil
	})
}

func (suite *LoggerTestSuite) TestNoMessageSupplied() {
	l := Logger{
		Zap: suite.zapLogger,
	}

	l.Log("value", 123)
	suite.assertOneMessage(func(e zapcore.Entry, fs []zap.Field) error {
		suite.Equal(NoLogMessage, e.Message)
		suite.Equal(zap.ErrorLevel, e.Level)
		suite.Equal(
			[]zap.Field{zap.Int("value", 123)},
			fs,
		)

		return nil
	})
}

func (suite *LoggerTestSuite) TestCustomMessageKey() {
	l := Logger{
		Zap:        suite.zapLogger,
		MessageKey: "custom",
	}

	l.Log("value", 123, "custom", "test message")
	suite.assertOneMessage(func(e zapcore.Entry, fs []zap.Field) error {
		suite.Equal("test message", e.Message)
		suite.Equal(zap.ErrorLevel, e.Level)
		suite.Equal(
			[]zap.Field{zap.Int("value", 123)},
			fs,
		)

		return nil
	})
}

func (suite *LoggerTestSuite) TestOddKeyvals() {
	l := Logger{
		Zap: suite.zapLogger,
	}

	l.Log("value", 123, DefaultMessageKey, "test message", "dangling")
	suite.assertOneMessage(func(e zapcore.Entry, fs []zap.Field) error {
		suite.Equal("test message", e.Message)
		suite.Equal(zap.ErrorLevel, e.Level)
		suite.ElementsMatch(
			[]zap.Field{zap.Int("value", 123), zap.NamedError("dangling", log.ErrMissingValue)},
			fs,
		)

		return nil
	})
}

func (suite *LoggerTestSuite) TestNotAString() {
	l := Logger{
		Zap: suite.zapLogger,
	}

	l.Log(45.6, 123, DefaultMessageKey, "test message")
	suite.assertOneMessage(func(e zapcore.Entry, fs []zap.Field) error {
		suite.Equal("test message", e.Message)
		suite.Equal(zap.ErrorLevel, e.Level)
		suite.ElementsMatch(
			[]zap.Field{zap.Any(NotAString, 123)},
			fs,
		)

		return nil
	})
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
