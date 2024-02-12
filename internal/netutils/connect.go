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
	"time"

	"github.com/rs/zerolog"
)

// HTTPTransportDialContextFunc represents a function that is compatible with
// the http.Transport DialContext field.
type HTTPTransportDialContextFunc func(ctx context.Context, network, addr string) (net.Conn, error)

// DialerWithContext returns a function for use with the http.Transport
// DialContext field. Use of this function allows the caller to override the
// default "auto" network type selection behavior used by the net.Dial
// function when opening a network connection to the specified address/port.
func DialerWithContext(networkType string, logger zerolog.Logger) HTTPTransportDialContextFunc {

	// This function is provided with an address value in host:port format.
	return func(ctx context.Context, network string, address string) (net.Conn, error) {
		logger = logger.With().
			Str("address", address).
			Str("net_type", networkType).
			Logger()

		logger.Debug().Msg("resolving hostname")

		host, port, splitErr := net.SplitHostPort(address)
		if splitErr != nil {
			return nil, fmt.Errorf(
				"failed to split given pattern %q into host and port pair: %w",
				address,
				splitErr,
			)
		}

		addrs, resolveErr := resolveIPAddresses(ctx, host, networkType, logger)
		if resolveErr != nil {
			return nil, fmt.Errorf(
				"resolve hostname %s to %s IPs: %w",
				host,
				networkTypeToIPTypeStr(networkType),
				resolveErr,
			)
		}

		conn, connectErr := openConnection(
			ctx,
			addrs,
			port,
			networkType,
			logger,
		)

		if connectErr != nil {
			return nil, fmt.Errorf(
				"failed to create client connection to %s (port %s): %w",
				host,
				port,
				connectErr,
			)
		}

		return conn, nil

	}
}

// func DialContext() HTTPTransportDialContextFunc {
//
// }

// openConnection receives a list of IP Addresses and returns a net.Conn value
// for the first successful connection attempt. An error is returned instead
// if one occurs.
func openConnection(ctx context.Context, addrs []string, port string, netType string, logger zerolog.Logger) (net.Conn, error) {
	if len(addrs) < 1 {
		logger.Error().Msg("empty list of IP Addresses received")

		return nil, fmt.Errorf(
			"empty list of IP Addresses received: %w",
			ErrMissingValue,
		)
	}

	var (
		c          net.Conn
		connectErr error
	)

	for _, addr := range addrs {
		logger.Debug().
			Str("ip_address", addr).
			Msg("Connecting to server")

		if err := ctx.Err(); err != nil {
			logger.Debug().Msg("context has expired")

			return nil, fmt.Errorf("failed to open connection: %w", err)
		}

		s := net.JoinHostPort(addr, port)

		// Unless sysadmin explicitly requested one of IPv4 or IPv6 network
		// types we fall back to default behavior.
		switch strings.ToLower(netType) {
		case NetTypeTCP4:
		case NetTypeTCP6:
		default:
			netType = NetTypeTCPAuto
		}

		// Ensure that dialer has required KeepAlive and Timeout values to
		// prevent connections from hanging indefinitely.
		//
		// TODO: Research & confirm whether this is still true. For now, play
		// it safe and use the suggested settings to enable reasonable network
		// timeout behavior.
		//
		// https://joshrendek.com/2015/09/using-a-custom-http-dialer-in-go/
		// https://pkg.go.dev/net#Dialer
		dialer := &net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 2 * time.Second,
		}

		// Attempt to connect to the given IP Address.
		c, connectErr = dialer.Dial(netType, s)

		if connectErr != nil {
			logger.Debug().
				Err(connectErr).
				Str("ip_address", addr).
				Msg("error connecting to server")

			continue
		}

		// If no connection errors were received, we can consider the
		// connection attempt a success and skip further attempts to connect
		// to any remaining IP Addresses for the specified server name.
		logger.Debug().
			Str("ip_address", addr).
			Msg("Connected to server")

		return c, nil
	}

	// If all connection attempts failed, report the last connection error.
	// Log all failed IP Addresses for review.
	if connectErr != nil {
		errMsg := fmt.Sprintf(
			"failed to connect to server using any of %d IP Addresses (%s)",
			len(addrs),
			strings.Join(addrs, ", "),
		)
		logger.Debug().
			Err(connectErr).
			Str("failed_ip_addresses", strings.Join(addrs, ", ")).
			Msg(errMsg)

		return nil, fmt.Errorf(
			"%s; last error: %v: %w",
			errMsg,
			connectErr,
			ErrNetworkConnectionFailed,
		)
	}

	return c, nil
}
