// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"github.com/atc0005/check-rsat/internal/config"
	"github.com/rs/zerolog"
)

func setupLogger(cfg *config.Config) zerolog.Logger {
	logger := cfg.Log.With().Logger()

	loggerVerbose := cfg.Log.With().Caller().
		Str("server", cfg.Server).
		Str("user", cfg.Username).
		Int("port", cfg.TCPPort).
		Str("net_type", cfg.NetworkType).
		Str("timeout", cfg.Timeout().String()).
		Bool("cert-validation-disabled", cfg.TrustCert).
		Bool("ca-cert-specified", cfg.CACertificate != "").
		Bool("permit-tls-renegotiation", cfg.PermitTLSRenegotiation).
		Str("version", config.Version()).
		Logger()

	if zerolog.GlobalLevel() == zerolog.DebugLevel ||
		zerolog.GlobalLevel() == zerolog.TraceLevel {

		logger = loggerVerbose
	}

	return logger
}
