// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"strings"

	"github.com/atc0005/check-rsat/internal/config"
	"github.com/atc0005/check-rsat/internal/rsat"
	"github.com/atc0005/go-nagios"
)

// setPluginOutput is a helper function used to set plugin output and state
// values.
func setPluginOutput(
	stateLabel string,
	message string,
	extendedMessage string,
	err error,
	orgs rsat.Organizations,
	cfg *config.Config,
	plugin *nagios.Plugin,
) {
	if err != nil {
		plugin.AddError(err)
	}

	plugin.ExitStatusCode = nagios.StateLabelToExitCode(stateLabel)

	plugin.ServiceOutput = fmt.Sprintf(
		"%s: %s",
		strings.ToUpper(stateLabel),
		message,
	)

	if cfg != nil {
		setLongServiceOutput(extendedMessage, orgs, cfg, plugin)
	}

}

func setLongServiceOutput(report string, _ rsat.Organizations, cfg *config.Config, plugin *nagios.Plugin) {
	var output strings.Builder

	// If provided, put the report content first.
	if report != "" {
		fmt.Fprintf(
			&output,
			"%s%s",
			report,
			nagios.CheckOutputEOL,
		)
	}

	if cfg.ShowVerbose {
		fmt.Fprintf(&output, "%s", nagios.CheckOutputEOL)

		fmt.Fprintf(
			&output,
			"%s------%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			&output,
			"Configuration settings: %s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			&output,
			"* Server: %v%s",
			cfg.Server,
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			&output,
			"* Port: %v%s",
			cfg.TCPPort,
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			&output,
			"* Username: %v%s",
			cfg.Username,
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			&output,
			"* NetworkType: %v%s",
			cfg.NetworkType,
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			&output,
			"* Timeout: %v%s",
			cfg.Timeout(),
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			&output,
			"* UserAgent: %v%s",
			cfg.UserAgent(),
			nagios.CheckOutputEOL,
		)
	}

	plugin.LongServiceOutput = output.String()
}
