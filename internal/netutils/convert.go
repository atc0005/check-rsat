// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package netutils

import (
	"fmt"
	"net"
	"strings"

	"github.com/rs/zerolog"
)

// networkTypeToIPTypeStr resolves a network type or name (e.g., tcp4, tcp6)
// to a human readable IP Address type.
func networkTypeToIPTypeStr(netType string) string {
	switch strings.ToLower(netType) {
	case NetTypeTCP4:
		return "IPv4"
	case NetTypeTCP6:
		return "IPv6"
	default:
		return "IPv4 or IPv6"
	}
}

func ipStringsToNetIPs(ipStrings []string, logger zerolog.Logger) ([]net.IP, error) {
	ips := make([]net.IP, 0, len(ipStrings))

	logger.Debug().Msg("converting DNS lookup results to net.IP values for net type validation")

	for i := range ipStrings {
		ip := net.ParseIP(ipStrings[i])
		if ip == nil {
			return nil, fmt.Errorf(
				"error parsing %s: %w",
				ipStrings[i],
				ErrIPAddressParsingFailed,
			)
		}

		ips = append(ips, ip)
	}

	// FIXME: Is this length check really needed? Presumably if there were
	// zero results from the parsing attempt an error would have ready been
	// returned?
	switch {
	case len(ips) < 1:
		errMsg := fmt.Sprintf(
			"failed to to convert DNS lookup results to net.IP values after receiving %d DNS lookup results ([%s])",
			len(ipStrings),
			strings.Join(ipStrings, ", "),
		)

		logger.Error().Msg(errMsg)

		return nil, fmt.Errorf(
			"%s: %w",
			errMsg,
			ErrIPAddressParsingFailed,
		)

	default:
		logger.Debug().Msg("successfully converted DNS lookup results to net.IP values")
	}

	return ips, nil
}

func netIPsToIPStrings(netIPs []net.IP) []string {
	ipStrs := make([]string, len(netIPs))
	for i := range netIPs {
		ipStrs[i] = netIPs[i].String()
	}

	return ipStrs
}
