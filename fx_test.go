// SPDX-FileCopyrightText: 2024 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package sallust

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FxSuite struct {
	suite.Suite
}

func (suite *FxSuite) testWithLoggerDefault() {
	var logger *zap.Logger
	app := fxtest.New(
		suite.T(),
		WithLogger(),
		fx.Populate(&logger),
	)

	app.RequireStart()
	app.RequireStop()
	suite.Require().NotNil(logger)
	logger.Error("discarded")
	suite.NoError(logger.Sync())
}

func (suite *FxSuite) testWithLoggerFull() {
	var (
		logger *zap.Logger

		hooksCalled []string

		suppliedHook = func(e zapcore.Entry) error {
			hooksCalled = append(hooksCalled, "suppliedHook")
			return nil
		}

		injectedHook = func(e zapcore.Entry) error {
			hooksCalled = append(hooksCalled, "injectedHook")
			return nil
		}

		app = fxtest.New(
			suite.T(),
			fx.Supply(
				Config{
					OutputPaths: []string{"stdout"},
				},
				[]zap.Option{
					zap.Hooks(injectedHook),
				},
			),
			WithLogger(
				zap.Hooks(suppliedHook),
			),
			fx.Populate(&logger),
		)
	)

	app.RequireStart()
	app.RequireStop()
	suite.Require().NotNil(logger)
	logger.Error("expected error")

	// verify that supplied options are applied first
	// each hook should have been called for each entry, which will include fx.App startup messages
	suite.Require().Zero(len(hooksCalled) % 2)
	for i := 0; i < len(hooksCalled); i += 2 {
		suite.Equal(
			[]string{"suppliedHook", "injectedHook"},
			hooksCalled[i:i+2],
		)
	}
}

func (suite *FxSuite) TestWithLogger() {
	suite.Run("Default", suite.testWithLoggerDefault)
	suite.Run("Full", suite.testWithLoggerFull)
}

type syncCore struct {
	zapcore.Core

	synced bool
}

func (sc *syncCore) Sync() error {
	sc.synced = true
	return sc.Core.Sync()
}

func (suite *FxSuite) TestSyncOnShutdown() {
	var (
		logger *zap.Logger

		app = fxtest.New(
			suite.T(),
			WithLogger(
				zap.WrapCore(
					func(c zapcore.Core) zapcore.Core {
						return &syncCore{
							Core: c,
						}
					},
				),
			),
			fx.Populate(&logger),
			SyncOnShutdown(),
		)
	)

	app.RequireStart()
	app.RequireStop()

	suite.Require().NotNil(logger)
	logger.Error("expected error")
	suite.True(logger.Core().(*syncCore).synced)
}

func TestFx(t *testing.T) {
	suite.Run(t, new(FxSuite))
}
