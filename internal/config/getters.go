// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"time"
)

// Timeout converts the user-specified connection timeout value in seconds to
// an appropriate time duration value for use with setting a timeout value.
func (c Config) Timeout() time.Duration {
	return time.Duration(c.timeout) * time.Second
}

// supportedLogLevels returns a list of valid log levels supported by tools in
// this project.
func supportedLogLevels() []string {
	return []string{
		LogLevelDisabled,
		LogLevelPanic,
		LogLevelFatal,
		LogLevelError,
		LogLevelWarn,
		LogLevelInfo,
		LogLevelDebug,
		LogLevelTrace,
	}
}

// supportedNetworkTypes returns a list of valid network types.
func supportedNetworkTypes() []string {
	return []string{
		netTypeTCPAuto,
		netTypeTCP4,
		netTypeTCP6,
	}
}

// supportedInspectorOutputFormats returns a list of valid output formats used
// by Inspector type applications in this project. This list is intended to be
// used for validating the user-specified output format.
func supportedInspectorOutputFormats() []string {
	return []string{
		InspectorOutputFormatOverview,
		InspectorOutputFormatSimpleTable,
		InspectorOutputFormatPrettyTable,
		InspectorOutputFormatVerbose,
	}
}

// UserAgent returns a string usable as-is as a custom user agent for plugins
// provided by this project.
func (c Config) UserAgent() string {
	// Default User Agent: (Go-http-client/1.1)
	// https://datatracker.ietf.org/doc/html/draft-ietf-httpbis-p2-semantics-22#section-5.5.3
	return fmt.Sprintf(
		"%s/%s",
		myAppName,
		version,
	)
}
