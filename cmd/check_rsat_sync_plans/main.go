// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

//go:generate go-winres make --product-version=git-tag --file-version=git-tag

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/atc0005/check-rsat/internal/config"
	"github.com/atc0005/check-rsat/internal/reports"
	"github.com/atc0005/check-rsat/internal/rsat"

	"github.com/atc0005/go-nagios"
	"github.com/rs/zerolog"
)

func main() {
	plugin := nagios.NewPlugin()

	// defer this from the start so it is the last deferred function to run
	defer plugin.ReturnCheckResults()

	// Setup configuration by parsing user-provided flags.
	cfg, cfgErr := config.New(config.AppType{Plugin: true})

	switch {
	case errors.Is(cfgErr, config.ErrVersionRequested):
		fmt.Println(config.Version())

		return

	case errors.Is(cfgErr, config.ErrHelpRequested):
		fmt.Println(cfg.Help())

		return

	case cfgErr != nil:
		// We make some assumptions when setting up our logger as we do not
		// have a working configuration based on sysadmin-specified choices.
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}
		logger := zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()

		logger.Err(cfgErr).Msg("Error initializing application")

		setPluginOutput(
			nagios.StateUNKNOWNLabel,
			"Error initializing application",
			"",
			cfgErr,
			nil,
			cfg,
			plugin,
		)

		return
	}

	// Annotate all errors (if any) with remediation advice just before ending
	// plugin execution.
	defer annotateErrors(plugin)

	// Set context deadline equal to user-specified timeout value for
	// runtime/execution.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout())
	defer cancel()

	if cfg.EmitBranding {
		// If enabled, show application details at end of notification
		plugin.BrandingCallback = config.Branding("Notification generated by ")
	}

	logger := cfg.Log.With().
		Str("server", cfg.Server).
		Str("user", cfg.Username).
		Int("port", cfg.TCPPort).
		Str("net_type", cfg.NetworkType).
		Str("timeout", cfg.Timeout().String()).
		Bool("cert-validation-disabled", cfg.TrustCert).
		Bool("ca-cert-specified", cfg.CACertificate != "").
		Bool("permit-tls-renegotiation", cfg.PermitTLSRenegotiation).
		Logger()

	logger.Debug().Msg("Beginning plugin execution")

	// If specified, attempt to load the CA certificate associated with the
	// Red Hat Satellite server's certificate chain.
	var caCert []byte
	if cfg.CACertificate != "" {
		logger.Debug().Msg("CA Cert specified: attempting to load CA cert")

		var readErr error
		caCert, readErr = os.ReadFile(cfg.CACertificate)
		if readErr != nil {
			setPluginOutput(
				nagios.StateUNKNOWNLabel,
				"Error loading CA certificate for Red Hat Satellite instance",
				"",
				readErr,
				nil,
				cfg,
				plugin,
			)

			return
		}

		logger.Debug().Msg("Successfully loaded CA cert")
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

	apiLimits := rsat.APILimits{
		PerPage: cfg.PerPageLimit,
	}

	client := rsat.NewAPIClient(authInfo, apiLimits, logger)

	orgs, orgsFetchErr := rsat.GetOrgsWithSyncPlans(ctx, client)
	if orgsFetchErr != nil {
		setPluginOutput(
			nagios.StateCRITICALLabel,
			"Error retrieving Red Hat Satellite sync plans",
			"",
			orgsFetchErr,
			orgs,
			cfg,
			plugin,
		)

		return
	}

	logger.Debug().
		Int("orgs", orgs.NumOrgs()).
		Int("sync_plans", orgs.NumPlans()).
		Msg("Retrieved sync plans")

	pd := getPerfData(orgs)
	if err := plugin.AddPerfData(false, pd...); err != nil {
		setPluginOutput(
			nagios.StateUNKNOWNLabel,
			"Failed to process performance data metrics",
			"",
			err,
			orgs,
			cfg,
			plugin,
		)

		return
	}

	switch {
	case !orgs.IsOKState():
		logger.Debug().Msg("Problem sync plans detected")

		setPluginOutput(
			orgs.ServiceState().Label,
			fmt.Sprintf(
				"%d problem sync plans detected for %s (evaluated %d orgs, %d sync plans)",
				orgs.NumProblemPlans(),
				cfg.Server,
				orgs.NumOrgs(),
				orgs.NumPlans(),
			),
			reports.SyncPlansVerboseReport(orgs, cfg, logger),
			nil,
			orgs,
			cfg,
			plugin,
		)

	default:
		logger.Debug().Msg("No problems detected")

		setPluginOutput(
			nagios.StateOKLabel,
			fmt.Sprintf(
				"No sync plans with non-OK status detected for %s (evaluated %d orgs, %d sync plans)",
				cfg.Server,
				orgs.NumOrgs(),
				orgs.NumPlans(),
			),
			reports.SyncPlansVerboseReport(orgs, cfg, logger),
			nil,
			orgs,
			cfg,
			plugin,
		)
	}

}
