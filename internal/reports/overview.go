// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package reports

import (
	"fmt"
	"strings"

	"github.com/atc0005/check-rsat/internal/config"
	"github.com/atc0005/check-rsat/internal/rsat"
	"github.com/atc0005/go-nagios"
	"github.com/rs/zerolog"
)

// SyncPlansOverviewReport provides a listing of Red Hat Satellite
// organizations and the overall (high-level) state of sync plans in each
// organization. This report is intentionally light on specifics.
func SyncPlansOverviewReport(orgs rsat.Organizations, _ *config.Config, _ zerolog.Logger) string {
	var output strings.Builder

	addSyncPlansReportLeadIn(&output)

	orgs.Sort()

	for _, org := range orgs {
		fmt.Fprintf(
			&output,
			"* %s (%d problems, %d enabled, %d disabled)%s",
			org.Name,
			org.SyncPlans.NumStuck(),
			org.SyncPlans.NumEnabled(),
			org.SyncPlans.NumDisabled(),
			nagios.CheckOutputEOL,
		)
	}

	return output.String()
}
