# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
- sonar integration

## [v0.1.3]
- Sane defaults for fields in zap.Config and zap.EncoderConfig
- Rename NewLogger to Build to properly override zap's behavior

## [v0.1.2]
- Added a mapstructure DecodeHook for zap and zapcore types used in configuration

## [v0.1.1]
- Rename Options to Config and use embedding to tidy up

## [v0.1.0]
- First release

[Unreleased]: https://github.com/xmidt-org/sallust/compare/v0.1.3..HEAD
[v0.1.3]: https://github.com/xmidt-org/sallust/compare/0.1.2...v0.1.3
[v0.1.2]: https://github.com/xmidt-org/sallust/compare/0.1.1...v0.1.2
[v0.1.1]: https://github.com/xmidt-org/sallust/compare/0.1.0...v0.1.1
[v0.1.0]: https://github.com/xmidt-org/sallust/compare/0.0.0...v0.1.0
