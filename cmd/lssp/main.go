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
	"github.com/atc0005/check-rsat/internal/rsat"

	"github.com/rs/zerolog"
)

func main() {
	// Setup configuration by parsing user-provided flags.
	cfg, cfgErr := config.New(config.AppType{Inspector: true})

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
		os.Exit(config.ExitCodeCatchall)
	}

	// Emulate returning exit code from main function by "queuing up" a
	// default exit code that matches expectations, but allow explicitly
	// setting the exit code in such a way that is compatible with using
	// deferred function calls throughout the application.
	var appExitCode int
	defer func(code *int) {
		var exitCode int
		if code != nil {
			exitCode = *code
		}
		os.Exit(exitCode)
	}(&appExitCode)

	// Set context deadline equal to user-specified timeout value for
	// runtime/execution.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout())
	defer cancel()

	logger := setupLogger(cfg)

	authInfo, authErr := getAuthInfo(cfg, logger)
	if authErr != nil {
		logger.Error().
			Err(authErr).
			Msg("Error preparing auth info for Red Hat Satellite instance")

		appExitCode = config.ExitCodeCatchall

		return
	}

	apiLimits := rsat.APILimits{
		PerPage: cfg.PerPageLimit,
	}

	client := rsat.NewAPIClient(authInfo, apiLimits, logger)

	logger.Info().
		Str("timeout", cfg.Timeout().String()).
		Msg("Retrieving Red Hat Satellite sync plans (this may take a while)")

	orgs, orgsFetchErr := rsat.GetOrgsWithSyncPlans(ctx, client)
	if orgsFetchErr != nil {
		logger.Error().
			Err(orgsFetchErr).
			Msg("Error retrieving Red Hat Satellite sync plans")

		appExitCode = config.ExitCodeCatchall

		return
	}

	logger.Info().
		Int("organizations", orgs.NumOrgs()).
		Int("sync_plans", orgs.NumPlans()).
		Msg("Retrieved sync plans")

	logger.Info().Msg("Evaluating sync plans")

	switch {
	case !orgs.IsOKState():
		logger.Warn().
			Int("total", orgs.NumPlans()).
			Int("enabled", orgs.NumPlansEnabled()).
			Int("disabled", orgs.NumPlansDisabled()).
			Int("problematic", orgs.NumProblemPlans()).
			Msg("Problem sync plans detected")

		generateReport(os.Stdout, orgs, cfg, logger)

	default:
		logger.Info().Msg("No problems detected")

		generateReport(os.Stdout, orgs, cfg, logger)
	}

}
