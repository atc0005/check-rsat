// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"

	"github.com/atc0005/check-rsat/internal/rsat"
	"github.com/atc0005/go-nagios"
)

// getPerfData gathers performance data metrics that we wish to report.
func getPerfData(orgs rsat.Organizations) []nagios.PerformanceData {
	switch {
	case len(orgs) == 0:
		return []nagios.PerformanceData{}

	default:
		return []nagios.PerformanceData{
			// The `time` (runtime) metric is appended at plugin exit, so do not
			// duplicate it here.
			{
				Label: "organizations",
				Value: fmt.Sprintf("%d", orgs.NumOrgs()),
			},
			{
				Label: "sync_plans_total",
				Value: fmt.Sprintf("%d", orgs.NumPlans()),
			},
			{
				Label: "sync_plans_enabled",
				Value: fmt.Sprintf("%d", orgs.NumPlansEnabled()),
			},
			{
				Label: "sync_plans_disabled",
				Value: fmt.Sprintf("%d", orgs.NumPlansDisabled()),
			},
			{
				Label: "sync_plans_stuck",
				Value: fmt.Sprintf("%d", orgs.NumPlansStuck()),
			},
			{
				Label: "sync_plans_problems",
				Value: fmt.Sprintf("%d", orgs.NumProblemPlans()),
			},
		}
	}

}
