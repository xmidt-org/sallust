package sallustkit

import (
	"testing"

	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogMethodForTestSuite struct {
	GokitTestSuite
}

var _ suite.SetupTestSuite = (*LogMethodForTestSuite)(nil)

func (suite *LogMethodForTestSuite) testLevel(expected zapcore.Level, actual level.Value) {
	lm := LogMethodFor(suite.zapLogger, actual)
	suite.Require().NotNil(lm)

	lm("test message", zap.Int("value", 123))
	suite.assertOneMessage(
		func(e zapcore.Entry, fs []zapcore.Field) error {
			suite.Equal("test message", e.Message)
			suite.Equal(expected, e.Level)

			if suite.Len(fs, 1) {
				suite.Equal(zap.Int("value", 123), fs[0])
			}

			return nil
		},
	)
}

func (suite *LogMethodForTestSuite) TestDebugLevel() {
	suite.testLevel(zapcore.DebugLevel, level.DebugValue())
}

func (suite *LogMethodForTestSuite) TestInfoLevel() {
	suite.testLevel(zapcore.InfoLevel, level.InfoValue())
}

func (suite *LogMethodForTestSuite) TestWarnLevel() {
	suite.testLevel(zapcore.WarnLevel, level.WarnValue())
}

func (suite *LogMethodForTestSuite) TestErrorLevel() {
	suite.testLevel(zapcore.ErrorLevel, level.ErrorValue())
}

func TestLogMethodFor(t *testing.T) {
	suite.Run(t, new(LogMethodForTestSuite))
}
