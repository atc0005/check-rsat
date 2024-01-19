# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/check-rsat/issues) for
any deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.1.5] - 2024-01-19

### Changed

#### Dependency Updates

- (GH-125) ghaw: bump github/codeql-action from 2 to 3
- (GH-131) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.2 to go-ci-oldstable-build-v0.14.3 in /dependabot/docker/builds
- (GH-127) go.mod: bump golang.org/x/sys from 0.15.0 to 0.16.0
- (GH-129) canary: bump golang from 1.20.12 to 1.20.13 in /dependabot/docker/go

## [v0.1.4] - 2023-12-08

### Changed

#### Dependency Updates

- (GH-117) canary: bump golang from 1.20.11 to 1.20.12 in /dependabot/docker/go
- (GH-119) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.1 to go-ci-oldstable-build-v0.14.2 in /dependabot/docker/builds
- (GH-116) go.mod: bump golang.org/x/sys from 0.14.0 to 0.15.0
- (GH-114) go.mod: bump zgo.at/termtext from 1.1.0 to 1.2.0

## [v0.1.3] - 2023-11-20

### Changed

#### Dependency Updates

- (GH-103) canary: bump golang from 1.20.10 to 1.20.11 in /dependabot/docker/go
- (GH-69) canary: bump golang from 1.20.7 to 1.20.8 in /dependabot/docker/go
- (GH-92) canary: bump golang from 1.20.8 to 1.20.10 in /dependabot/docker/go
- (GH-94) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.10 to go-ci-oldstable-build-v0.13.12 in /dependabot/docker/builds
- (GH-105) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.12 to go-ci-oldstable-build-v0.14.1 in /dependabot/docker/builds
- (GH-60) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.4 to go-ci-oldstable-build-v0.13.5 in /dependabot/docker/builds
- (GH-61) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.5 to go-ci-oldstable-build-v0.13.6 in /dependabot/docker/builds
- (GH-63) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.6 to go-ci-oldstable-build-v0.13.7 in /dependabot/docker/builds
- (GH-71) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.7 to go-ci-oldstable-build-v0.13.8 in /dependabot/docker/builds
- (GH-76) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.8 to go-ci-oldstable-build-v0.13.9 in /dependabot/docker/builds
- (GH-81) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.9 to go-ci-oldstable-build-v0.13.10 in /dependabot/docker/builds
- (GH-67) ghaw: bump actions/checkout from 3 to 4
- (GH-98) go.mod: bump github.com/mattn/go-isatty from 0.0.19 to 0.0.20
- (GH-84) go.mod: bump github.com/rs/zerolog from 1.30.0 to 1.31.0
- (GH-66) go.mod: bump golang.org/x/sys from 0.11.0 to 0.12.0
- (GH-89) go.mod: bump golang.org/x/sys from 0.12.0 to 0.13.0
- (GH-102) go.mod: bump golang.org/x/sys from 0.13.0 to 0.14.0

### Fixed

- (GH-108) Fix goconst linting errors

## [v0.1.2] - 2023-08-17

### Added

- (GH-28) Add initial automated release notes config
- (GH-30) Add initial automated release build workflow

### Changed

- Dependencies
  - `Go`
    - `1.19.11` to `1.20.7`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.11.3` to `go-ci-oldstable-build-v0.13.4`
  - `rs/zerolog`
    - `v1.29.1` to `v1.30.0`
  - `mattn/go-runewidth`
    - `v0.0.14` to `v0.0.15`
  - `golang.org/x/sys`
    - `v0.10.0` to `v0.11.0`
- (GH-32) Update Dependabot config to monitor both branches
- (GH-54) Update project to Go 1.20 series

## [v0.1.1] - 2023-07-13

### Overview

- Dependency updates
- built using Go 1.19.11
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.19.10` to `1.19.11`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.11.0` to `go-ci-oldstable-build-v0.11.3`

## [v0.1.0] - 2023-07-07

Initial release!

This release provides early versions of tooling used to evaluate Red Hat
Satellite (RSAT) instances. This evaluation is performed using official APIs.

### Added

- `lssp`, a CLI app to list Red Hat Satellite sync plans.
- `check_rsat_sync_plans`, a Nagios plugin used to monitor for problematic Red
  Hat Satellite (RSAT) sync plans.

Just to be 100% clear: this project is not affiliated with or
endorsed by  Red Hat, Inc.

[Unreleased]: https://github.com/atc0005/check-rsat/compare/v0.1.5...HEAD
[v0.1.5]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.5
[v0.1.4]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.4
[v0.1.3]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.3
[v0.1.2]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.0
