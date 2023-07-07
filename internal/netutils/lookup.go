// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package netutils

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/rs/zerolog"
)

func lookupIPs(ctx context.Context, server string, logger zerolog.Logger) ([]string, error) {
	if err := ctx.Err(); err != nil {
		logger.Debug().Msg("context has expired")

		return nil, fmt.Errorf("failed to lookup IPs: %w", err)
	}

	logger.Debug().Str("host", server).Msg("Performing name resolution")

	resolver := &net.Resolver{}
	lookupResults, lookupErr := resolver.LookupHost(ctx, server)
	if lookupErr != nil {
		logger.Error().
			Err(lookupErr).
			Str("server", server).
			Msg("error resolving hostname")

		return nil, fmt.Errorf(
			"error resolving hostname %s: %v: %w",
			server,
			lookupErr,
			ErrDNSLookupFailed,
		)
	}

	// FIXME: Is this length check really needed? Presumably if there were
	// zero results returned an error would have also been returned?
	switch {
	case len(lookupResults) < 1:
		errMsg := fmt.Sprintf(
			"failed to resolve hostname %s to IP Addresses",
			server,
		)

		logger.Error().
			Str("server", server).
			Msg(errMsg)

		return nil, fmt.Errorf(
			"%s: %w",
			errMsg,
			ErrDNSLookupFailed,
		)

	default:
		logger.Debug().
			Int("count", len(lookupResults)).
			Str("ips", strings.Join(lookupResults, ", ")).
			Str("server", server).
			Msg("successfully resolved IP Addresses for hostname")
	}

	return lookupResults, nil
}

func resolveIPAddresses(ctx context.Context, server string, networkType string, logger zerolog.Logger) ([]string, error) {
	if err := ctx.Err(); err != nil {
		logger.Debug().Msg("context has expired")

		return nil, fmt.Errorf("failed to resolve IPs: %w", err)
	}

	lookupResults, lookupErr := lookupIPs(ctx, server, logger)
	if lookupErr != nil {
		return nil, lookupErr
	}

	netIPs, ipConvertErr := ipStringsToNetIPs(lookupResults, logger)
	if ipConvertErr != nil {
		return nil, ipConvertErr
	}

	filteredNetIPs, filterIPsErr := filterNetIPsToNetworkType(netIPs, networkType, logger)
	if filterIPsErr != nil {
		return nil, filterIPsErr
	}

	ipStrings := netIPsToIPStrings(filteredNetIPs)

	return ipStrings, nil
}
