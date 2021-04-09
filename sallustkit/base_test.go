package sallustkit

import (
	"io/ioutil"

	"github.com/stretchr/testify/suite"
	"github.com/xmidt-org/sallust"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GokitTestSuite is embedded by other suites to inherit
// the test setup logic.
type GokitTestSuite struct {
	suite.Suite

	core      *sallust.CaptureCore
	zapLogger *zap.Logger
}

var _ suite.SetupTestSuite = (*GokitTestSuite)(nil)

func (suite *GokitTestSuite) SetupTest() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	suite.core = sallust.NewCaptureCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(ioutil.Discard),
		zap.DebugLevel, // turn on everything
	)

	suite.zapLogger = zap.New(suite.core)
}

func (suite *GokitTestSuite) assertOneMessage(f func(zapcore.Entry, []zap.Field) error) {
	c, err := suite.core.EachMessage(f)
	suite.NoError(err, "unexpected error from EachMessage")
	suite.Equal(1, c, "wrong number of captured messages")
}
