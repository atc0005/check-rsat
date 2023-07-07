// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package netutils

import (
	"errors"
)

//
// FIXME: Move general network related functionality to a separate package.
//

// Known, named networks used for TCP/IP connections. These names match the
// network names used by the `net` standard library package.
const (

	// NetTypeTCPAuto indicates that either of IPv4 or IPv6 will be used to
	// establish a connection depending on the specified IP Address.
	NetTypeTCPAuto string = "tcp"

	// NetTypeTCP4 indicates an IPv4-only network.
	NetTypeTCP4 string = "tcp4"

	// NetTypeTCP6 indicates an IPv6-only network.
	NetTypeTCP6 string = "tcp6"
)

var (
	// ErrMissingValue indicates that an expected value was missing.
	ErrMissingValue = errors.New("missing expected value")

	// ErrDNSLookupFailed indicates a failure to resolve a hostname to an IP
	// Address.
	ErrDNSLookupFailed = errors.New("failed to resolve hostname")

	// ErrIPAddressParsingFailed indicates a failure to parse a given value as
	// an IP Address.
	ErrIPAddressParsingFailed = errors.New("failed to parse IP Address")

	// ErrNoIPAddressesForChosenNetworkType indicates a failure to obtain any
	// IP Addresses of the specified network type (e.g., IPv4 vs IPv6) for a
	// given hostname.
	ErrNoIPAddressesForChosenNetworkType = errors.New("no resolved IP Addresses for chosen network type")

	// ErrNetworkConnectionFailed indicates a failure to establish a network
	// connection to the specified host.
	ErrNetworkConnectionFailed = errors.New("failed to establish network connection")
)
