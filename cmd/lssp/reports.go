// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"io"

	"github.com/atc0005/check-rsat/internal/config"
	"github.com/atc0005/check-rsat/internal/reports"
	"github.com/atc0005/check-rsat/internal/rsat"
	"github.com/rs/zerolog"
)

func generateReport(w io.Writer, orgs rsat.Organizations, cfg *config.Config, logger zerolog.Logger) {
	logger.Info().Msg("Generating sync plans report")

	switch cfg.InspectorOutputFormat {
	case config.InspectorOutputFormatOverview:
		fmt.Fprintln(w, reports.SyncPlansOverviewReport(orgs, cfg, logger))

	case config.InspectorOutputFormatSimpleTable:
		fmt.Fprintln(w, reports.SyncPlansSimpleTableReport(orgs, cfg, logger))

	case config.InspectorOutputFormatPrettyTable:
		fmt.Fprintln(w, reports.SyncPlansPrettyTableReport(orgs, cfg, logger))

	case config.InspectorOutputFormatVerbose:
		fmt.Fprintln(w, reports.SyncPlansVerboseReport(orgs, cfg, logger))
	}

}
