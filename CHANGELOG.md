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

## [v0.1.0] - 2022-04-xx

Initial release!

This release provides early versions of tooling used to evaluate Red Hat
Satellite (RSAT) instances. This evaluation is performed using official APIs.

### Added

- `lsrs`, a CLI app to list details from a Red Hat Satellite (RSAT) instance.

- `check_rsat_subscriptions`, a Nagios plugin used to monitor subscription
  thresholds for one or more RSAT Organizations.

[Unreleased]: https://github.com/atc0005/check-rsat/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/atc0005/check-rsat/releases/tag/v0.1.0
