# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.2.2]

### Added
- [Configurable file permissions #51](https://github.com/xmidt-org/sallust/issues/51)

## [v0.2.1]
- [Pruned deprecated code #43](https://github.com/xmidt-org/sallust/pull/43)

## [v0.2.0]
- [Migrate Useful webpa-common/logging Utilities #37](https://github.com/xmidt-org/sallust/issues/37)
- [Enable & Fix Linter #36](https://github.com/xmidt-org/sallust/issues/36)

## [v0.1.6]

### Fixed
- levelEncoder, timeEncoder, durationEncoder, callerEncoder, nameEncoder marshaling
  due to upstream viper change to support mapstructure.
### Added
- Bootstrapping a zap logger for a fx application, including the fxevent.Logger.
- SyncOnShutdown added enabling sync of the logs prior to app shutting down.
### Deprecated
- Deprecated Buffer and CaptureCore.
### Changed
- Dependencies have been updated.

## [v0.1.5]

### Added
- sallustkit package adapts go-kit's logging onto zap [#15](https://github.com/xmidt-org/sallust/issues/15)

## [v0.1.4]
- sonar integration
- use a custom Config and EncoderConfig that are friendlier to libraries like viper

## [v0.1.3]
- Sane defaults for fields in zap.Config and zap.EncoderConfig
- Rename NewLogger to Build to properly override zap's behavior

## [v0.1.2]
- Added a mapstructure DecodeHook for zap and zapcore types used in configuration

## [v0.1.1]
- Rename Options to Config and use embedding to tidy up

## [v0.1.0]
- First release

[Unreleased]: https://github.com/xmidt-org/sallust/compare/v0.2.2..HEAD
[v0.2.2]: https://github.com/xmidt-org/sallust/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/xmidt-org/sallust/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/xmidt-org/sallust/compare/v0.1.6...v0.2.0
[v0.1.6]: https://github.com/xmidt-org/sallust/compare/v0.1.5...v0.1.6
[v0.1.5]: https://github.com/xmidt-org/sallust/compare/v0.1.4...v0.1.5
[v0.1.4]: https://github.com/xmidt-org/sallust/compare/v0.1.3...v0.1.4
[v0.1.3]: https://github.com/xmidt-org/sallust/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/xmidt-org/sallust/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/xmidt-org/sallust/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/xmidt-org/sallust/compare/v0.0.0...v0.1.0
