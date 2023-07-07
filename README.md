<!-- omit in toc -->

# check-rsat

Go-based tooling to monitor Red Hat Satellite (RSAT) systems; NOT affiliated
with or endorsed by Red Hat, Inc.

<!-- omit in toc -->

[![Latest Release](https://img.shields.io/github/release/atc0005/check-rsat.svg?style=flat-square)](https://github.com/atc0005/check-rsat/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/check-rsat.svg)](https://pkg.go.dev/github.com/atc0005/check-rsat)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/atc0005/check-rsat)](https://github.com/atc0005/check-rsat)
[![Lint and Build](https://github.com/atc0005/check-rsat/actions/workflows/lint-and-build.yml/badge.svg)](https://github.com/atc0005/check-rsat/actions/workflows/lint-and-build.yml)
[![Project Analysis](https://github.com/atc0005/check-rsat/actions/workflows/project-analysis.yml/badge.svg)](https://github.com/atc0005/check-rsat/actions/workflows/project-analysis.yml)

<!-- omit in toc -->
## Table of Contents

- [check-rsat](#check-rsat)
  - [Project home](#project-home)
  - [Overview](#overview)
    - [Output](#output)
    - [`check_rsat_sync_plans`](#check_rsat_sync_plans)
      - [Performance Data](#performance-data)
    - [`lssp`](#lssp)
  - [Features](#features)
    - [`check_rsat_sync_plans`](#check_rsat_sync_plans-1)
    - [`lssp`](#lssp-1)
    - [common](#common)
  - [Changelog](#changelog)
  - [Requirements](#requirements)
    - [Building source code](#building-source-code)
    - [Running](#running)
  - [Installation](#installation)
    - [From source](#from-source)
    - [Using release binaries](#using-release-binaries)
    - [Deployment](#deployment)
  - [Configuration options](#configuration-options)
    - [Command-line arguments](#command-line-arguments)
      - [`check_rsat_sync_plans`](#check_rsat_sync_plans-2)
      - [`lssp`](#lssp-2)
    - [Configuration file](#configuration-file)
  - [Examples](#examples)
    - [`check_rsat_sync_plans` Nagios plugin](#check_rsat_sync_plans-nagios-plugin)
      - [CLI invocations](#cli-invocations)
        - [Omit OK sync plans](#omit-ok-sync-plans)
        - [List all sync plans](#list-all-sync-plans)
      - [Command definitions](#command-definitions)
    - [`lssp` CLI app](#lssp-cli-app)
      - [The `pretty-table` format (default)](#the-pretty-table-format-default)
      - [The `overview` format](#the-overview-format)
      - [The `verbose` format](#the-verbose-format)
      - [Other output formats](#other-output-formats)
  - [License](#license)
  - [References](#references)

## Project home

See [our GitHub repo][repo-url] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

Just to be 100% clear: this project is not affiliated with or endorsed by
Red Hat, Inc.

## Overview

This repo contains various tools and plugins used to monitor Red Hat Satellite
(RSAT) systems.

| Plugin or Tool Name     | Description                                                                        |
| ----------------------- | ---------------------------------------------------------------------------------- |
| `check_rsat_sync_plans` | Nagios plugin used to monitor for problematic Red Hat Satellite (RSAT) sync plans. |
| `lssp`                  | CLI app to list Red Hat Satellite sync plans.                                      |

### Output

The output for plugins in this project is designed to provide the one-line
summary needed by Nagios (and other monitoring systems) for quick
identification of a problem while providing longer, more detailed information
for display within the web UI, use in email and Teams notifications
([atc0005/send2teams](https://github.com/atc0005/send2teams)).

By default, output intended for processing by Nagios is sent to `stdout` and
output intended for troubleshooting by the sysadmin is sent to `stderr`.
Output emitted to `stderr` is  configurable via the `--log-level` flag.

For some monitoring systems or addons (e.g., Icinga Web 2, Nagios XI), the
`stderr` output is mixed in with the `stdout` output in the web UI for the
service check. This may add visual noise when viewing the service check
output. For those cases, you may wish to explicitly disable the output to
`stderr` by using the `--log-level "disabled"` CLI flag & value.

### `check_rsat_sync_plans`

Nagios plugin used to monitor for problematic Red Hat Satellite (RSAT) sync
plans.

See the [configuration options](#configuration-options) section for details
regarding supported flags and values.

See the [features list](#features) for the validation checks currently
supported by this plugin.

#### Performance Data

Initial support has been added for emitting Performance Data / Metrics, but
refinement suggestions are welcome.

Consult the table below for the metrics implemented thus far.

Please add to an existing
[Discussion](https://github.com/atc0005/check-rsat/discussions) thread
(if applicable) or [open a new
one](https://github.com/atc0005/check-rsat/discussions/new) with any
feedback that you may have. Thanks in advance!

| Emitted Performance Data / Metric | Meaning                                                             |
| --------------------------------- | ------------------------------------------------------------------- |
| `time`                            | Runtime for plugin                                                  |
| `organizations`                   | Number of organizations                                             |
| `sync_plans_total`                | Number of total sync plans                                          |
| `sync_plans_enabled`              | Number of sync plans in an enabled state                            |
| `sync_plans_disabled`             | Number of sync plans in an disabled state                           |
| `sync_plans_stuck`                | Number of sync plans in a "stuck" state                             |
| `sync_plans_problems`             | Number of sync plans in a non-OK (*needs sysadmin attention*) state |

### `lssp`

CLI app used to generate an overview of the Red Hat Satellite sync plans along
with their current state (e.g., disabled, enabled, next scheduled sync time,
overall status).

## Features

### `check_rsat_sync_plans`

- Evaluate sync plans from all Red Hat Satellite organizations
  - this includes monitoring for "stuck" sync plans. While these plans have an
    `enabled` state they have a `Sync Date` value set to a time in the past.
    These plans are effectively disabled until a sysadmin takes action to
    resolve the issue (e.g., create a new recurring logic & associate it with
    the sync plan).

### `lssp`

- List sync plans from all Red Hat Satellite organizations
  - multiple output formats
    - `overview`
    - `simple-table`
    - `pretty-table`
    - `verbose`

### common

Features common to all tools provided by this project.

- Optional, leveled logging using `rs/zerolog` package
  - [`logfmt`][logfmt] format output
    - to `stderr` for plugins
    - to `stdout` for CLI apps

- Optional, user-specified timeout value for plugin execution

- Optional override of network type
  - defaults to either of IPv4 and IPv6
  - optionally limited to IPv4-only or IPv6-only

- Optional, user-specified read limit
  - helps protect against excessive/unexpected input size

- Optional support for omitting sync plans in an `OK` state
  - help focus on just the sync plans with a "problem" status

- Optional support for accepting renegotiation requests from the Red Hat
  Satellite server
  - this support is disabled by default
  - renegotiation is not supported for TLS 1.3

- Optional use of specified CA certificate to validate Red Hat Satellite
  certificate chain

- Optional disabling of certificate validation
  - WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this
  option.

- Optional branding "signature"
  - appended at the end of plugin output
  - used to indicate what Nagios plugin (and what version) is responsible for
    the service check result

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Requirements

The following is a loose guideline. Other combinations of Go and operating
systems for building and running tools from this repo may work, but have not
been tested.

### Building source code

- Go
  - see this project's `go.mod` file for *preferred* version
  - this project tests against [officially supported Go
    releases][go-supported-releases]
    - the most recent stable release (aka, "stable")
    - the prior, but still supported release (aka, "oldstable")
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

### Running

- Windows 11
- Ubuntu Linux 22.04
- Red Hat Enterprise Linux 8

## Installation

### From source

1. [Download][go-docs-download] Go
1. [Install][go-docs-install] Go
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/check-rsat`
   1. `cd check-rsat`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc`
   - for CentOS Linux
     1. `sudo yum install make gcc`
1. Build
   - manually, explicitly specifying target OS and architecture
     - `GOOS=linux GOARCH=amd64 go build -mod=vendor ./cmd/check_rsat_sync_plans/`
       - most likely this is what you want (if building manually)
       - substitute `amd64` with the appropriate architecture if using
         different hardware (e.g., `arm64` or `386`)
     - `GOOS=linux GOARCH=amd64 go build -mod=vendor ./cmd/lssp/`
       - see notes above
   - using Makefile `linux` recipe
     - `make linux`
       - generates x86 and x64 binaries
   - using Makefile `release-build` recipe
     - `make release-build`
       - generates the same release assets as provided by this project's
         releases
1. Locate generated binaries
   - if using `Makefile`
     - look in `/tmp/check-rsat/release_assets/check_rsat_sync_plans/`
     - look in `/tmp/check-rsat/release_assets/lssp/`
   - if using `go build`
     - look in `/tmp/check-rsat/`
1. Copy the applicable binaries to whatever systems needs to run them so that
   they can be deployed

**NOTE**: Depending on which `Makefile` recipe you use the generated binary
may be compressed and have an `xz` extension. If so, you should decompress the
binary first before deploying it (e.g., `xz -d
check_rsat_sync_plans-linux-amd64.xz`).

### Using release binaries

1. Download the [latest release][repo-url] binaries
1. Decompress binaries
   - e.g., `xz -d check_rsat_sync_plans-linux-amd64.xz`
1. Copy the applicable binaries to whatever systems needs to run them so that
   they can be deployed

**NOTE**:

DEB and RPM packages are provided as an alternative to manually deploying
binaries.

### Deployment

1. Place `check_rsat_sync_plans` in a location where it can be executed by the
   monitoring agent
   - Usually the same place as other Nagios plugins
   - For example, on a default Red Hat Enterprise Linux system using
   `check_nrpe` the `check_rsat_sync_plans` plugin would be deployed to
   `/usr/lib64/nagios/plugins/check_rsat_sync_plans` or
   `/usr/local/nagios/libexec/check_rsat_sync_plans`
1. Place `lssp` in a location where it can be easily accessed
   - Usually the same place as other custom tools installed outside of your
     package manager's control
   - e.g., `/usr/local/bin/lssp`

**NOTE**:

DEB and RPM packages are provided as an alternative to manually deploying
binaries.

## Configuration options

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

#### `check_rsat_sync_plans`

| Flag                       | Required | Default   | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                      |
| -------------------------- | -------- | --------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`                 | No       | `false`   | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                             |
| `h`, `help`                | No       | `false`   | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                           |
| `v`, `version`             | No       | `false`   | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                    |
| `ll`, `log-level`          | No       | `info`    | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                              |
| `t`, `timeout`             | No       | `10`      | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                           |
| `omit-ok`                  | No       | `false`   | No     | `true`, `false`                                                         | Whether sync plans listed in plugin output should be limited to just those in a non-OK state.                                                                                                                                                                                                    |
| `read-limit`               | No       | `1048576` | No     | *valid whole number of bytes*                                           | Limit in bytes used to help prevent abuse when reading input that could be larger than expected. The default value is nearly 4x the largest observed (formatted) feed size.                                                                                                                      |
| `verbose`                  | No       | `false`   | No     | `true`, `false`                                                         | Whether to display verbose details in the final plugin output.                                                                                                                                                                                                                                   |
| `server`                   | Yes      | *empty*   | No     | *fully-qualified domain name or IP Address*                             | The Red Hat Satellite server FQDN or IP Address.                                                                                                                                                                                                                                                 |
| `username`                 | Yes      | *empty*   | No     | *valid user account*                                                    | The valid user for the given Red Hat Satellite server.                                                                                                                                                                                                                                           |
| `password`                 | Yes      | *empty*   | No     | *valid password*                                                        | The valid password for the specified user.                                                                                                                                                                                                                                                       |
| `port`                     | No       | `443`     | No     | *positive whole number between 1-65535, inclusive*                      | The port used by the Red Hat Satellite server API.                                                                                                                                                                                                                                               |
| `permit-tls-renegotiation` | No       | `false`   | No     | `true`, `false`                                                         | Whether support for accepting renegotiation requests from the Red Hat Satellite server are permitted. This support is disabled by default. Renegotiation is not supported for TLS 1.3.                                                                                                           |
| `trust-cert`               | No       | `false`   | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                            |
| `net-type`                 | No       | `auto`    | No     | `tcp4`, `tcp6`, `auto`                                                  | Limits network connections to one of tcp4 (IPv4-only), tcp6 (IPv6-only) or auto (either).                                                                                                                                                                                                        |
| `ca-cert`                  | No       | *empty*   | No     | *valid path to file*                                                    | CA Certificate used to validate the certificate chain used by the Red Hat Satellite server. This is usually the path to the CA cert provided by the `katello-ca-consumer-latest.noarch.rpm` package which is installed as part of registering a RHEL instance with a Red Hat Satellite instance. |

#### `lssp`

| Flag                       | Required | Default   | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                      |
| -------------------------- | -------- | --------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `h`, `help`                | No       | `false`   | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                           |
| `v`, `version`             | No       | `false`   | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                    |
| `ll`, `log-level`          | No       | `info`    | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                              |
| `t`, `timeout`             | No       | `10`      | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                           |
| `omit-ok`                  | No       | `false`   | No     | `true`, `false`                                                         | Whether sync plans listed in plugin output should be limited to just those in a non-OK state.                                                                                                                                                                                                    |
| `read-limit`               | No       | `1048576` | No     | *valid whole number of bytes*                                           | Limit in bytes used to help prevent abuse when reading input that could be larger than expected. The default value is nearly 4x the largest observed (formatted) feed size.                                                                                                                      |
| `output-format`            | No       | `table`   | No     | `overview`, `simple-table`, `pretty-table`, `verbose`                   | Sets output format. The default format is `pretty-table`.                                                                                                                                                                                                                                        |
| `server`                   | Yes      | *empty*   | No     | *fully-qualified domain name or IP Address*                             | The Red Hat Satellite server FQDN or IP Address.                                                                                                                                                                                                                                                 |
| `username`                 | Yes      | *empty*   | No     | *valid user account*                                                    | The valid user for the given Red Hat Satellite server.                                                                                                                                                                                                                                           |
| `password`                 | Yes      | *empty*   | No     | *valid password*                                                        | The valid password for the specified user.                                                                                                                                                                                                                                                       |
| `port`                     | No       | `443`     | No     | *positive whole number between 1-65535, inclusive*                      | The port used by the Red Hat Satellite server API.                                                                                                                                                                                                                                               |
| `permit-tls-renegotiation` | No       | `false`   | No     | `true`, `false`                                                         | Whether support for accepting renegotiation requests from the Red Hat Satellite server are permitted. This support is disabled by default. Renegotiation is not supported for TLS 1.3.                                                                                                           |
| `trust-cert`               | No       | `false`   | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                            |
| `net-type`                 | No       | `auto`    | No     | `tcp4`, `tcp6`, `auto`                                                  | Limits network connections to one of tcp4 (IPv4-only), tcp6 (IPv6-only) or auto (either).                                                                                                                                                                                                        |
| `ca-cert`                  | No       | *empty*   | No     | *valid path to file*                                                    | CA Certificate used to validate the certificate chain used by the Red Hat Satellite server. This is usually the path to the CA cert provided by the `katello-ca-consumer-latest.noarch.rpm` package which is installed as part of registering a RHEL instance with a Red Hat Satellite instance. |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Examples

Entries in this section attempt to provide a brief overview of usage. Please
file an issue if coverage is found to be unclear or incorrect.

### `check_rsat_sync_plans` Nagios plugin

#### CLI invocations

##### Omit OK sync plans

This example omits all of the OK status sync plans. We also omit the bulk of
the entries only showing the first 3 orgs and the last one.

```console
$ /usr/lib/nagios/plugins/check_rsat_sync_plans --server rsat.example.com --port 443 --username $RSAT_USER --password $RSAT_PASSWORD --ca-cert /etc/rhsm/ca/katello-server-ca.pem --timeout 240 --permit-tls-renegotiation --log-level info --omit-ok
OK: No sync plans with non-OK status detected for rsat.example.com (evaluated 20 orgs, 58 sync plans)


SYNC PLANS OVERVIEW

* Org1 (0 enabled, 3 disabled)
* Org2 (0 enabled, 3 disabled)
* Org3 (3 enabled, 0 disabled)

...

* Org20 (3 enabled, 0 disabled)


 | 'organizations'=20;;;; 'sync_plans_disabled'=32;;;; 'sync_plans_enabled'=26;;;; 'sync_plans_problems'=0;;;; 'sync_plans_stuck'=0;;;; 'sync_plans_total'=58;;;; 'time'=134627ms;;;;
```

##### List all sync plans

This example does not omit the OK status sync plans. We do however omit the
bulk of the entries only showing the first 3 and last org.

```console
$ /usr/lib/nagios/plugins/check_rsat_sync_plans --server rsat.example.com --port 443 --username $RSAT_USER --password $RSAT_PASSWORD --ca-cert /etc/rhsm/ca/katello-server-ca.pem --timeout 240 --permit-tls-renegotiation --log-level disabled
OK: No sync plans with non-OK status detected for rsat.example.com (evaluated 20 orgs, 58 sync plans)


SYNC PLANS OVERVIEW


Org1 (0 enabled, 3 disabled):
  * [Name: Base OS, Interval: hourly, Next Sync: N/A]
  * [Name: OS Related, Interval: daily, Next Sync: N/A]
  * [Name: Other, Interval: weekly, Next Sync: N/A]

Org2 (0 enabled, 3 disabled):
  * [Name: Base OS, Interval: hourly, Next Sync: N/A]
  * [Name: OS Related, Interval: daily, Next Sync: N/A]
  * [Name: Other, Interval: weekly, Next Sync: N/A]

Org3 (3 enabled, 0 disabled):
  * [Name: Base OS, Interval: hourly, Next Sync: 2023-07-06 07:43:00 CDT]
  * [Name: OS Related, Interval: daily, Next Sync: 2023-07-06 16:43:00 CDT]
  * [Name: Other, Interval: weekly, Next Sync: 2023-07-12 16:43:00 CDT]

Org20 (3 enabled, 0 disabled):
  * [Name: Base OS, Interval: hourly, Next Sync: 2023-07-06 08:12:00 CDT]
  * [Name: OS Related, Interval: daily, Next Sync: 2023-07-06 21:12:00 CDT]
  * [Name: Other, Interval: weekly, Next Sync: 2023-07-10 21:12:00 CDT]


 | 'organizations'=20;;;; 'sync_plans_disabled'=32;;;; 'sync_plans_enabled'=26;;;; 'sync_plans_problems'=0;;;; 'sync_plans_stuck'=0;;;; 'sync_plans_total'=58;;;; 'time'=128017ms;;;;
```

#### Command definitions

The command definition file below defines several commands for use in service
check definitions.

Each command explicitly excludes "OK" sync plans in order to keep the output
manageable. Remove the `omit-ok` flag if you wish to list the current status
of all sync plans.

```nagios
# /etc/nagios-plugins/config/rsat-sync-plans.cfg

define command{
    command_name    check_rsat_sync_plans
    command_line    /usr/lib/nagios/plugins/check_rsat_sync_plans --server '$ARG1$' --port '$ARG2$' --username '$ARG3$' --password '$ARG4$' --timeout 240 --ca-cert '/etc/rhsm/ca/katello-server-ca.pem' --permit-tls-renegotiation --omit-ok --log-level disabled
    }

define command{
    command_name    check_rsat_sync_plans_custom_timeout
    command_line    /usr/lib/nagios/plugins/check_rsat_sync_plans --server '$ARG1$' --port '$ARG2$' --username '$ARG3$' --password '$ARG4$' --timeout '$ARG5$' --ca-cert '/etc/rhsm/ca/katello-server-ca.pem' --permit-tls-renegotiation --omit-ok --log-level disabled
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each.

### `lssp` CLI app

#### The `pretty-table` format (default)

This example uses the default timeout of 300s (5m) and also specifies the CA
cert to perform certificate chain validation for the Red Hat Satellite API
endpoints. We permit TLS renegotiation as that has been found to be needed for
some Satellite instances.

The current behavior is to emit an empty table if no problems were detected
and the flag to omit OK sync plans was specified. This is subject to change in
a future version.

```console
$ /usr/local/bin/lssp --server rsat.example.com --port 443 --username $RSAT_USER --password $RSAT_PASSWORD --ca-cert /etc/rhsm/ca/katello-server-ca.pem --permit-tls-renegotiation --log-level info --omit-ok
7:36AM INF Attempting to load specified CA cert ca-cert=/etc/rhsm/ca/katello-server-ca.pem
7:36AM INF Successfully loaded CA cert
7:36AM INF Retrieving Red Hat Satellite sync plans (this may take a while) timeout=5m0s
7:38AM INF Retrieved sync plans organizations=20 sync_plans=58
7:38AM INF Evaluating sync plans
7:38AM INF No problems detected
7:38AM INF Generating sync plans report

SYNC PLANS OVERVIEW

┌────────────┬─────────────┬───────────┬────────────┬─────────────┬──────────┐
│  Org Name  │  Plan Name  │  Enabled  │  Interval  │  Next Sync  │  Status  │
├────────────┼─────────────┼───────────┼────────────┼─────────────┼──────────┤
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
│            │             │           │            │             │          │
└────────────┴─────────────┴───────────┴────────────┴─────────────┴──────────┘
```

Here is an example where we did not specify the flag to omit OK results:

```console
$ /usr/local/bin/lssp --server rsat.example.com --port 443 --username $RSAT_USER --password $RSAT_PASSWORD --ca-cert /etc/rhsm/ca/katello-server-ca.pem --permit-tls-renegotiation --log-level info
7:41AM INF Attempting to load specified CA cert ca-cert=/etc/rhsm/ca/katello-server-ca.pem
7:41AM INF Successfully loaded CA cert
7:41AM INF Retrieving Red Hat Satellite sync plans (this may take a while) timeout=5m0s
7:43AM INF Retrieved sync plans organizations=20 sync_plans=58
7:43AM INF Evaluating sync plans
7:43AM INF No problems detected
7:43AM INF Generating sync plans report

SYNC PLANS OVERVIEW

┌────────────┬────────────────────────────────────────┬───────────┬────────────┬───────────────────────────┬──────────┐
│  Org Name  │               Plan Name                │  Enabled  │  Interval  │         Next Sync         │  Status  │
├────────────┼────────────────────────────────────────┼───────────┼────────────┼───────────────────────────┼──────────┤
│  Org1      │  Base OS                               │  false    │  hourly    │  Not scheduled            │    ✔     │
│  Org1      │  OS Related                            │  false    │  daily     │  Not scheduled            │    ✔     │
│  Org1      │  Other                                 │  false    │  weekly    │  Not scheduled            │    ✔     │
│            │                                        │           │            │                           │          │
│  Org2      │  Base OS                               │  false    │  hourly    │  Not scheduled            │    ✔     │
│  Org2      │  OS Related                            │  false    │  daily     │  Not scheduled            │    ✔     │
│  Org2      │  Other                                 │  false    │  weekly    │  Not scheduled            │    ✔     │
│            │                                        │           │            │                           │          │
│  Org3     │  Base OS                               │  true     │  hourly    │  2023-07-06 07:43:00 CDT  │    ✔     │
│  Org3     │  OS Related                            │  true     │  daily     │  2023-07-06 16:43:00 CDT  │    ✔     │
│  Org3     │  Other                                 │  true     │  weekly    │  2023-07-12 16:43:00 CDT  │    ✔     │
│            │                                        │           │            │                           │          │
│  Org20    │  Base OS                               │  true     │  hourly    │  2023-07-06 08:12:00 CDT  │    ✔     │
│  Org20    │  OS Related                            │  true     │  daily     │  2023-07-06 21:12:00 CDT  │    ✔     │
│  Org20    │  Other                                 │  true     │  weekly    │  2023-07-10 21:12:00 CDT  │    ✔     │
└────────────┴────────────────────────────────────────┴───────────┴────────────┴───────────────────────────┴──────────┘
```

The formatting in the console is a little different than shown above (right
border is displayed accurately), but this gives an overall idea of the output.

There are also 20 organizations with the majority having 3 plans each. The
output above is trimmed to just the first 3 and last 1.

#### The `overview` format

This example uses the default timeout of 300s (5m) and also specifies the CA
cert to perform certificate chain validation for the Red Hat Satellite API
endpoints. We permit TLS renegotiation as that has been found to be needed for
some Satellite instances.

```console
$ /usr/local/bin/lssp --server rsat.example.com --port 443 --username $RSAT_USER --password $RSAT_PASSWORD --ca-cert /etc/rhsm/ca/katello-server-ca.pem --permit-tls-renegotiation --log-level info --output-format overview
7:51AM INF Attempting to load specified CA cert ca-cert=/etc/rhsm/ca/katello-server-ca.pem
7:51AM INF Successfully loaded CA cert
7:51AM INF Retrieving Red Hat Satellite sync plans (this may take a while) timeout=5m0s
7:54AM INF Retrieved sync plans organizations=20 sync_plans=58
7:54AM INF Evaluating sync plans
7:54AM INF No problems detected
7:54AM INF Generating sync plans report

SYNC PLANS OVERVIEW

* Org1 (0 problems, 0 enabled, 3 disabled)
* Org2 (0 problems, 0 enabled, 3 disabled)
* Org3 (0 problems, 3 enabled, 0 disabled)

...

* Org20 (0 problems, 3 enabled, 0 disabled)
```

#### The `verbose` format

The behavior between this format and the `overview` format is nearly identical
when using the flag to omit OK sync plans. If not opting to omit OK sync plans
(as the output below shows), the sync plan name, interval and next scheduled
sync time is shown. This can be useful for diagnosing "stuck" sync plans and
seeing at a glance which orgs have disabled sync plans.

```console
$ /usr/local/bin/lssp --server rsat.example.com --port 443 --username $RSAT_USER --password $RSAT_PASSWORD --ca-cert /etc/rhsm/ca/katello-server-ca.pem --permit-tls-renegotiation --log-level info --output-format verbose
7:58AM INF Attempting to load specified CA cert ca-cert=/etc/rhsm/ca/katello-server-ca.pem
7:58AM INF Successfully loaded CA cert
7:58AM INF Retrieving Red Hat Satellite sync plans (this may take a while) timeout=5m0s
8:00AM INF Retrieved sync plans organizations=20 sync_plans=58
8:00AM INF Evaluating sync plans
8:00AM INF No problems detected
8:00AM INF Generating sync plans report

SYNC PLANS OVERVIEW

Org1 (0 enabled, 3 disabled):
  * [Name: Base OS, Interval: hourly, Next Sync: N/A]
  * [Name: OS Related, Interval: daily, Next Sync: N/A]
  * [Name: Other, Interval: weekly, Next Sync: N/A]

Org2 (0 enabled, 3 disabled):
  * [Name: Base OS, Interval: hourly, Next Sync: N/A]
  * [Name: OS Related, Interval: daily, Next Sync: N/A]
  * [Name: Other, Interval: weekly, Next Sync: N/A]

Org3 (3 enabled, 0 disabled):
  * [Name: Base OS, Interval: hourly, Next Sync: 2023-07-06 08:43:00 CDT]
  * [Name: OS Related, Interval: daily, Next Sync: 2023-07-06 16:43:00 CDT]
  * [Name: Other, Interval: weekly, Next Sync: 2023-07-12 16:43:00 CDT]

Org20 (3 enabled, 0 disabled):
  * [Name: Base OS, Interval: hourly, Next Sync: 2023-07-06 08:12:00 CDT]
  * [Name: OS Related, Interval: daily, Next Sync: 2023-07-06 21:12:00 CDT]
  * [Name: Other, Interval: weekly, Next Sync: 2023-07-10 21:12:00 CDT]
```

#### Other output formats

Other output formats are also available. See the [configuration
options](#configuration-options) for more information.

## License

From the [LICENSE](LICENSE) file:

```license
MIT License

Copyright (c) 2023 Adam Chalkley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## References

- Other monitoring/plugin projects
  - <https://github.com/atc0005/check-cert>
  - <https://github.com/atc0005/check-illiad>
  - <https://github.com/atc0005/check-mail>
  - <https://github.com/atc0005/check-process>
  - <https://github.com/atc0005/check-restart>
  - <https://github.com/atc0005/check-ssh>
  - <https://github.com/atc0005/check-statuspage>
  - <https://github.com/atc0005/check-vmware>
  - <https://github.com/atc0005/check-whois>
  - <https://github.com/atc0005/send2teams>
  - <https://github.com/atc0005/nagios-debug>
  - <https://github.com/atc0005/go-nagios>

- Red Hat Satellite
  - <https://access.redhat.com/documentation/en-us/red_hat_satellite>
  - [katello-server-ca.pem certificate missing under /etc/rhsm/ca on a content
    host registered to Red Hat Satellite
    6](https://access.redhat.com/solutions/3484571)
  - <https://access.redhat.com/documentation/en-us/red_hat_satellite/6.0/html/installation_guide/sect-red_hat_satellite-installation_guide-configuring_rednbsphat_satellite_with_a_custom_server_certificate>
  - <https://access.redhat.com/documentation/en-us/red_hat_satellite/6.7/html/administering_red_hat_satellite/chap-red_hat_satellite-administering_red_hat_satellite-accessing_red_hat_satellite>

- Logging
  - <https://github.com/rs/zerolog>

- Nagios
  - <https://github.com/atc0005/go-nagios>
  - <https://nagios-plugins.org/doc/guidelines.html>
  - <https://www.monitoring-plugins.org/doc/guidelines.html>
  - <https://icinga.com/docs/icinga-2/latest/doc/05-service-monitoring/>

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/check-rsat>  "This project's GitHub repo"

[go-docs-download]: <https://golang.org/dl>  "Download Go"

[go-docs-install]: <https://golang.org/doc/install>  "Install Go"

[go-supported-releases]: <https://go.dev/doc/devel/release#policy> "Go Release Policy"

[logfmt]: <https://brandur.org/logfmt>

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
