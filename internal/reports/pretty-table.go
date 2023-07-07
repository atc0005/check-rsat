// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package reports

import (
	"io"
	"strings"

	"github.com/atc0005/check-rsat/internal/config"
	"github.com/atc0005/check-rsat/internal/rsat"
	"github.com/rs/zerolog"
	"zgo.at/acidtab"
)

// SyncPlansPrettyTableReport provides a report of Red Hat Satellite
// organizations in "pretty" table format. This table format uses more visual
// "polish" while attempting to remain compatible with modern terminals.
//
// Each sync plan is listed along with relevant status information.
func SyncPlansPrettyTableReport(orgs rsat.Organizations, cfg *config.Config, _ zerolog.Logger) string {
	var output strings.Builder

	addSyncPlansReportLeadIn(&output)

	orgs.Sort()

	syncPlansPrettyTableReport(&output, cfg, orgs)

	return output.String()
}

// prettyTableFormatColumnHeader is a helper function to format a given column
// header for use in a "pretty table" report.
func prettyTableFormatColumnHeader(s string) string {
	return "\x1b[1m" + s + "\x1b[0m"
}

// prettyTableProblemState is a helper function that formats a given state
// (problem present, or not) for use in a "pretty table" report as a status
// indicator.
func prettyTableProblemState(v interface{}) string {
	if b, ok := v.(bool); ok {
		return map[bool]string{
			false: "\x1b[32m ✔ \x1b[0m",
			true:  "\x1b[31m ✘ \x1b[0m",
		}[b]
	}
	return "\x00"
}

// syncPlansPrettyTableReport is a helper function that performs the bulk of
// the pretty table report output logic.
func syncPlansPrettyTableReport(w io.Writer, cfg *config.Config, orgs rsat.Organizations) {
	var t *acidtab.Table
	switch {
	case orgs.NumProblemPlans() > 0:
		t = acidtab.New(
			prettyTableFormatColumnHeader("Org Name"),
			prettyTableFormatColumnHeader("Plan Name"),
			prettyTableFormatColumnHeader("Days Stuck"),
			prettyTableFormatColumnHeader("Enabled"),
			prettyTableFormatColumnHeader("Interval"),
			prettyTableFormatColumnHeader("Next Sync"),
			prettyTableFormatColumnHeader("Status"),
		).
			Close(acidtab.CloseAll).
			AlignCol(6, acidtab.Center).
			FormatColFunc(6, prettyTableProblemState)

	default:
		t = acidtab.New(
			prettyTableFormatColumnHeader("Org Name"),
			prettyTableFormatColumnHeader("Plan Name"),
			prettyTableFormatColumnHeader("Enabled"),
			prettyTableFormatColumnHeader("Interval"),
			prettyTableFormatColumnHeader("Next Sync"),
			prettyTableFormatColumnHeader("Status"),
		).
			Close(acidtab.CloseAll).
			AlignCol(5, acidtab.Center).
			FormatColFunc(5, prettyTableProblemState)
	}

	for i, org := range orgs {
		for _, syncPlan := range org.SyncPlans {
			switch {
			case syncPlan.IsOKState() && cfg.OmitOKSyncPlans:
				continue

			case orgs.NumProblemPlans() > 0:
				t.Row(
					org.Name,
					syncPlan.Name,
					syncPlan.DaysStuckHR(),
					syncPlan.Enabled,
					syncPlan.Interval,
					syncPlan.NextSync.String(),
					!syncPlan.IsOKState(),
				)

			default:
				t.Row(
					org.Name,
					syncPlan.Name,
					syncPlan.Enabled,
					syncPlan.Interval,
					syncPlan.NextSync.String(),
					!syncPlan.IsOKState(),
				)
			}
		}

		// Group sync plans visually based on Org.
		if i+1 < len(orgs) {
			t.Row()
		}
	}

	t.Horizontal(w)
}
