// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package reports

import (
	"fmt"
	"io"
	"strings"

	"github.com/atc0005/check-rsat/internal/config"
	"github.com/atc0005/check-rsat/internal/rsat"
	"github.com/atc0005/go-nagios"
	"github.com/rs/zerolog"
)

// SyncPlansVerboseReport provides a verbose listing of Red Hat Satellite
// organizations and sync plans. This is useful for providing a detailed
// (while still manageable) report of the status of all sync plans in each
// organization.
//
// NOTE: If no problems are detected the output
func SyncPlansVerboseReport(orgs rsat.Organizations, cfg *config.Config, _ zerolog.Logger) string {
	var output strings.Builder

	addSyncPlansReportLeadIn(&output)

	orgs.Sort()

	syncPlansVerboseReport(&output, cfg, orgs)

	return output.String()
}

// syncPlansVerboseReport is a helper function that performs the bulk of
// the "verbose" report output logic.
func syncPlansVerboseReport(w io.Writer, cfg *config.Config, orgs rsat.Organizations) {
	for _, org := range orgs {
		switch {
		// If no problems to report and user opted to omit OK results we just
		// list the Orgs here with summary details. We will skip listing the
		// sync plans within each org.
		case orgs.NumProblemPlans() > 0:
			fmt.Fprintf(
				w,
				"%s%s (%d stuck, %d enabled, %d disabled)%s",
				nagios.CheckOutputEOL,
				org.Name,
				org.SyncPlans.NumStuck(),
				org.SyncPlans.NumEnabled(),
				org.SyncPlans.NumDisabled(),
				nagios.CheckOutputEOL,
			)
			continue

		default:
			fmt.Fprintf(
				w,
				"* %s (%d enabled, %d disabled)%s",
				org.Name,
				org.SyncPlans.NumEnabled(),
				org.SyncPlans.NumDisabled(),
				nagios.CheckOutputEOL,
			)

		}

		for _, syncPlan := range org.SyncPlans {
			switch {
			case syncPlan.IsOKState() && cfg.OmitOKSyncPlans:
				continue

			// We evaluate the collection as a whole vs just this specific
			// sync plan so that we can have consistency across each "row"; we
			// want to include "days stuck" even if the specific sync plan we
			// are looking at isn't stuck (to contrast against any plans which
			// are stuck).
			case orgs.NumProblemPlans() > 0:
				fmt.Fprintf(
					w,
					"  * [Name: %s, Days Stuck: %s, Interval: %s, Next Sync: %s]%s",
					syncPlan.Name,
					syncPlan.DaysStuckHR(),
					syncPlan.Interval,
					syncPlan.NextSync.String(),
					nagios.CheckOutputEOL,
				)

			default:
				fmt.Fprintf(
					w,
					"  * [Name: %s, Interval: %s, Next Sync: %s]%s",
					syncPlan.Name,
					syncPlan.Interval,
					syncPlan.NextSyncTime(),
					nagios.CheckOutputEOL,
				)
			}
		}

		fmt.Fprint(w, nagios.CheckOutputEOL)
	}
}
