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
	"text/tabwriter"

	"github.com/atc0005/check-rsat/internal/config"
	"github.com/atc0005/check-rsat/internal/rsat"
	"github.com/atc0005/go-nagios"
	"github.com/rs/zerolog"
)

// addSimpleTableDataSeparatorRow generates a separator row intended to be
// used between groupings in a "simple table" report. The number of "columns"
// in the generated separator row is of the same length as the provided header
// row.
func simpleTableDataSeparatorRow(headerRow string, headerRowSeparator string) string {
	numFields := len(strings.Split(headerRow, headerRowSeparator))

	return fmt.Sprint(
		strings.Repeat("\t", numFields),
		nagios.CheckOutputEOL,
	)
}

// simpleTableProblemStateToString is a helper function that formats a given
// state (problem present, or not) for use in a "simple table" report as a
// status indicator.
func simpleTableProblemStateToString(v interface{}) string {
	if b, ok := v.(bool); ok {
		return map[bool]string{
			false: "  OK  ",
			true:  "  !!  ",
		}[b]
	}
	return ""
}

// syncPlansSimpleTableReport is a helper function that performs the bulk of
// the "simple table" report output logic.
func syncPlansSimpleTableReport(w io.Writer, cfg *config.Config, headerRow string, dataRowTmpl string, orgs rsat.Organizations) {
	_, _ = fmt.Fprintln(w, headerRow)
	_, _ = fmt.Fprintln(w, simpleTableHeaderSeparatorRow(headerRow, "\t"))

	for i, org := range orgs {
		for _, syncPlan := range org.SyncPlans {
			switch {
			case syncPlan.IsOKState() && cfg.OmitOKSyncPlans:
				continue

			case orgs.NumProblemPlans() > 0:
				_, _ = fmt.Fprintf(
					w,
					dataRowTmpl,
					org.Name,
					syncPlan.Name,
					syncPlan.DaysStuckHR(),
					syncPlan.Interval,
					syncPlan.NextSync.String(),
					simpleTableProblemStateToString(!syncPlan.IsOKState()),
				)

			default:
				_, _ = fmt.Fprintf(
					w,
					dataRowTmpl,
					org.Name,
					syncPlan.Name,
					syncPlan.Interval,
					syncPlan.NextSync.String(),
					simpleTableProblemStateToString(!syncPlan.IsOKState()),
				)
			}
		}

		// Group sync plans visually based on Org.
		if i+1 < len(orgs) {
			_, _ = fmt.Fprint(w, simpleTableDataSeparatorRow(headerRow, "\t"))
		}
	}
}

// addHeaderSeparatorRow generates a separator row intended to be used between
// the header and data rows. Each "column" in the generated separator row
// template is of the same length as the header row column above it.
func simpleTableHeaderSeparatorRow(headerRow string, headerRowSeparator string) string {
	var row strings.Builder

	headerTmplItems := strings.Split(headerRow, headerRowSeparator)

	// Drop the last trailing tab character from the slice.
	if len(headerTmplItems) > 0 {
		headerTmplItems = headerTmplItems[:len(headerTmplItems)-1]
	}

	for _, item := range headerTmplItems {
		_, _ = fmt.Fprint(&row, strings.Repeat("-", len(item)))
		_, _ = fmt.Fprint(&row, (headerRowSeparator))
	}

	return row.String()
}

// SyncPlansSimpleTableReport provides a report of Red Hat Satellite
// organizations in "simple" table format. This table format is intentionally
// simple in an effort for the broadest compatible output.
//
// Each sync plan is listed along with relevant status information.
func SyncPlansSimpleTableReport(orgs rsat.Organizations, cfg *config.Config, logger zerolog.Logger) string {
	var output strings.Builder

	tw := tabwriter.NewWriter(&output, 4, 4, 4, ' ', 0)

	addSyncPlansReportLeadIn(&output)

	// Add some lead-in spacing to better separate any earlier log messages from
	// summary output
	_, _ = fmt.Fprintf(tw, "\n\n")

	orgs.Sort()

	var (
		headerRow   string
		dataRowTmpl string
	)

	// REMINDER: Column cells must be tab-terminated, not tab-separated:
	// non-tab terminated trailing text at the end of a line forms a cell but
	// that cell is not part of an aligned column.
	switch {
	case orgs.NumProblemPlans() > 0:
		headerRow = "Org Name\tPlan Name\tDays Stuck\tInterval\tNext Sync\tStatus\t"
		dataRowTmpl = "%s\t%s\t%s\t%s\t%s\t%s\t\n"
	default:
		headerRow = "Org Name\tPlan Name\tInterval\tNext Sync\tStatus\t"
		dataRowTmpl = "%s\t%s\t%s\t%s\t%s\t\n"
	}

	syncPlansSimpleTableReport(tw, cfg, headerRow, dataRowTmpl, orgs)

	_, _ = fmt.Fprintln(tw)

	if err := tw.Flush(); err != nil {
		logger.Error().Err(err).Msg("Error flushing tabwriter")
	}

	return output.String()
}
