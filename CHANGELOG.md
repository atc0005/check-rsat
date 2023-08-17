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

[Unreleased]: https://github.com/atc0005/check-rsat/compare/v0.1.2...HEAD
[v0.1.2]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.2
[v0.1.1]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.0
