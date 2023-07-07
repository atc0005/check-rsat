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

func filterNetIPsToIPv4(netIPs []net.IP, logger zerolog.Logger) []net.IP {
	filteredIPs := make([]net.IP, 0, len(netIPs))

	for i := range netIPs {
		if netIPs[i].To4() != nil {
			logger.Debug().
				Str("ipv4_address", netIPs[i].String()).
				Msg("matched IPv4 address")

			filteredIPs = append(filteredIPs, netIPs[i])
		}
	}

	return filteredIPs
}

func filterNetIPsToIPv6(netIPs []net.IP, logger zerolog.Logger) []net.IP {
	filteredIPs := make([]net.IP, 0, len(netIPs))

	for i := range netIPs {
		if netIPs[i].To4() == nil {
			// If earlier attempts to parse the IP Address succeeded (by way
			// of it being a net.IP value), but this is not considered an IPv4
			// address, we will consider it a valid IPv6 address.
			logger.Debug().
				Str("ipv6_address", netIPs[i].String()).
				Msg("matched IPv6 address")

			filteredIPs = append(filteredIPs, netIPs[i])
		}
	}

	return filteredIPs
}

func filterNetIPsToNetworkType(netIPs []net.IP, netType string, logger zerolog.Logger) ([]net.IP, error) {
	var filteredIPs []net.IP

	// Flag validation ensures that we see valid named networks as supported
	// by the `net` stdlib package, along with the "auto" keyword. Here we pay
	// attention to only the valid named networks. Since we're working with
	// user specified keywords, we compare case-insensitively.
	switch strings.ToLower(netType) {
	case NetTypeTCP4:
		logger.Debug().Msg("user opted for IPv4-only connectivity, gathering only IPv4 addresses")

		filteredIPs = filterNetIPsToIPv4(netIPs, logger)

	case NetTypeTCP6:
		logger.Debug().Msg("user opted for IPv6-only connectivity, gathering only IPv6 addresses")

		filteredIPs = filterNetIPsToIPv6(netIPs, logger)

	// either of IPv4 or IPv6 is acceptable
	default:
		logger.Debug().Msg("auto behavior enabled, gathering all addresses")

		filteredIPs = netIPs
	}

	// No IPs remain after filtering against IPv4-only or IPv6-only
	// requirement.
	switch {
	case len(filteredIPs) < 1:
		errMsg := fmt.Sprintf(
			"failed to gather IP Addresses when filtering %d IPs by specified network type %s ([%s])",
			len(netIPs),
			netType,
			strings.Join(netIPsToIPStrings(netIPs), ", "),
		)

		logger.Error().Msg(errMsg)

		return nil, fmt.Errorf(
			"%s: %w",
			errMsg,
			ErrNoIPAddressesForChosenNetworkType,
		)

	default:
		logger.Debug().
			Int("num_input_ips", len(netIPs)).
			Int("num_remaining_ips", len(filteredIPs)).
			Str("network_type", netType).
			Str("ips", strings.Join(netIPsToIPStrings(filteredIPs), ", ")).
			Msg("successfully gathered IP Addresses for specified network type")
	}

	return filteredIPs, nil
}
