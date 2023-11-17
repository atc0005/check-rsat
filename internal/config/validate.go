// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"strings"

	"github.com/atc0005/check-rsat/internal/textutils"
)

// validate verifies all Config struct fields have been provided acceptable
// values.
func (c Config) validate(appType AppType) error {

	// Shared validation
	switch {
	case strings.TrimSpace(c.Server) == "":
		return fmt.Errorf(
			"%w: missing server FQDN or IP Address",
			ErrUnsupportedOption,
		)

	case strings.TrimSpace(c.Username) == "":
		return fmt.Errorf(
			"%w: missing username",
			ErrUnsupportedOption,
		)

	case strings.TrimSpace(c.Password) == "":
		return fmt.Errorf(
			"%w: missing password",
			ErrUnsupportedOption,
		)

	// TCP Port 0 is used by server applications to indicate that they should
	// bind to an available port. Specifying port 0 for a client application
	// is not useful.
	case c.TCPPort <= 0:
		return fmt.Errorf(
			"%w: invalid TCP port number %d",
			ErrUnsupportedOption,
			c.TCPPort,
		)

	case c.Timeout() <= 0:
		return fmt.Errorf(
			"invalid timeout value %d provided: %w",
			c.Timeout(),
			ErrUnsupportedOption,
		)

	case c.ReadLimit <= 0:
		return fmt.Errorf(
			"invalid read limit value %d provided: %w",
			c.ReadLimit,
			ErrUnsupportedOption,
		)

	case c.TrustCert && c.CACertificate != "":
		return fmt.Errorf(
			"invalid combination of flags; only one of %s or %s flags are permitted: %w",
			TrustCertFlagLong,
			CACertificateFlagLong,
			ErrUnsupportedOption,
		)

	case !textutils.InList(c.NetworkType, supportedNetworkTypes(), true):
		return fmt.Errorf(
			"%w: invalid network type; got %v, expected one of %v",
			ErrUnsupportedOption,
			c.NetworkType,
			supportedNetworkTypes(),
		)

	case !textutils.InList(c.LoggingLevel, supportedLogLevels(), true):
		return fmt.Errorf(
			"%w: invalid logging level; got %v, expected one of %v",
			ErrUnsupportedOption,
			c.LoggingLevel,
			supportedLogLevels(),
		)
	}

	switch {
	case appType.Inspector:

		supportedFormats := supportedInspectorOutputFormats()
		if !textutils.InList(c.InspectorOutputFormat, supportedFormats, true) {
			return fmt.Errorf(
				"%w: invalid output format; got %v, expected one of %v",
				ErrUnsupportedOption,
				c.InspectorOutputFormat,
				supportedFormats,
			)
		}

	case appType.Plugin:

		// Placeholder for future plugin-specific validation.

	}

	// Optimist
	return nil
}
