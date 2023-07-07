// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"os"
	"path/filepath"

	"github.com/atc0005/check-rsat/internal/config"
	"github.com/atc0005/check-rsat/internal/rsat"
	"github.com/rs/zerolog"
)

func getAuthInfo(cfg *config.Config, logger zerolog.Logger) (rsat.APIAuthInfo, error) {

	// If specified, attempt to load the CA certificate associated with the
	// Red Hat Satellite server's certificate chain.
	var caCert []byte
	if cfg.CACertificate != "" {
		logger.Info().
			Str("ca-cert", cfg.CACertificate).
			Msg("Attempting to load specified CA cert")

		var readErr error
		caCert, readErr = os.ReadFile(filepath.Clean(cfg.CACertificate))
		if readErr != nil {
			logger.Error().
				Err(readErr).
				Msg("Error loading CA certificate for Red Hat Satellite instance")
			return rsat.APIAuthInfo{}, readErr
		}

		logger.Info().Msg("Successfully loaded CA cert")
	}

	authInfo := rsat.APIAuthInfo{
		Server:                 cfg.Server,
		Port:                   cfg.TCPPort,
		NetworkType:            cfg.NetworkType,
		ReadLimit:              cfg.ReadLimit,
		Username:               cfg.Username,
		Password:               cfg.Password,
		UserAgent:              cfg.UserAgent(),
		TrustCert:              cfg.TrustCert,
		PermitTLSRenegotiation: cfg.PermitTLSRenegotiation,
		CACert:                 caCert,
	}

	return authInfo, nil
}
