// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package rsat

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"
)

// syncTimeGraceMinutes indicates how much "grace" time should be applied
// between the next scheduled time a sync plan should run and the current
// time. Other tasks may conflict with the sync plan's execution and place it
// in a pending state for longer than expected. This time is intended to
// offset that delay and help avoid false positive reports of stuck sync
// plans.
const syncTimeGraceMinutes float64 = 5

// SyncPlansResponse represents the API response from a request of all sync
// plans for a specific organization.
type SyncPlansResponse struct {
	Error     NullString  `json:"error"`
	Search    NullString  `json:"search"`
	SyncPlans SyncPlans   `json:"results"`
	Sort      SortOptions `json:"sort"`
	Subtotal  int         `json:"subtotal"`
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	PerPage   int         `json:"per_page"`
}

// SyncPlan represents a Red Hat Satellite sync plan. Sync plans are used to
// schedule execution of content synchronization.
type SyncPlan struct {
	OriginalSyncDate  SyncTime            `json:"sync_date"`
	NextSync          SyncTime            `json:"next_sync"`
	UpdatedAt         StandardAPITime     `json:"updated_at"`
	CreatedAt         StandardAPITime     `json:"created_at"`
	Products          Products            `json:"products"`
	CronExpression    NullString          `json:"cron_expression"`
	Description       NullString          `json:"description"`
	Interval          string              `json:"interval"`
	Name              string              `json:"name"`
	OrganizationName  string              `json:"-"`
	OrganizationLabel string              `json:"-"`
	OrganizationTitle string              `json:"-"`
	RecurringLogicID  int                 `json:"foreman_tasks_recurring_logic_id"`
	ID                int                 `json:"id"`
	OrganizationID    int                 `json:"organization_id"`
	Permissions       SyncPlanPermissions `json:"permissions"`
	Enabled           bool                `json:"enabled"`
}

// SyncPlanPermissions is the collection of permissions that a user querying
// the Red Hat Satellite API has for interacting with sync plans.
type SyncPlanPermissions struct {
	DestroySyncPlans bool `json:"destroy_sync_plans"`
	EditSyncPlans    bool `json:"edit_sync_plans"`
	ViewSyncPlans    bool `json:"view_sync_plans"`
}

// Product is a collection of content repositories used to group custom
// repositories.
type Product struct {
	LastSync        StandardAPITime `json:"last_sync"`
	Description     NullString      `json:"description"`
	CpID            string          `json:"cp_id"`
	Label           string          `json:"label"`
	LastSyncText    string          `json:"last_sync_words"`
	Name            string          `json:"name"`
	SyncState       string          `json:"sync_state"`
	ID              int             `json:"id"`
	RepositoryCount int             `json:"repository_count"`
}

// Products is a collection of product values associated with a Red Hat
// Satellite sync plan.
type Products []Product

// SyncPlans is a collection of Red Hat Satellite sync plans.
type SyncPlans []SyncPlan

// GetSyncPlans uses the provided APIClient to retrieve all sync plans for
// each specified Red Hat Satellite organization. If no organizations are
// specified then an attempt will be made to retrieve sync plans from all RSAT
// organizations.
func GetSyncPlans(ctx context.Context, client *APIClient, orgs ...Organization) (SyncPlans, error) {
	funcTimeStart := time.Now()

	if client == nil {
		return nil, fmt.Errorf(
			"required API client was not provided: %w",
			ErrMissingValue,
		)
	}

	logger := client.Logger

	if len(orgs) == 0 {
		var orgsErr error
		orgs, orgsErr = GetOrganizations(ctx, client)
		if orgsErr != nil {
			return nil, orgsErr
		}
	}

	// We'll assume a default set of 3 sync plans per Org as a preallocation
	// starting point.
	allSyncPlans := make([]SyncPlan, 0, len(orgs)*3)

	reqsCounter := newRequestsCounter(len(orgs))

	for _, org := range orgs {

		subLogger := logger.With().
			Int("org_id", org.ID).
			Str("org_name", org.Name).
			Logger()

		retrievalStart := time.Now()

		subLogger.Debug().Msg("Retrieving sync plans for organization")

		syncPlans, err := getOrgSyncPlans(ctx, client, org)
		if err != nil {
			return nil, err
		}

		requestNum, requestsRemaining := reqsCounter()

		// If we are processing in bulk use the requests counter to provide
		// additional debugging context, otherwise keep the messages simple as
		// this function may be used by the caller to process bulk items and
		// may prefer to build a tally there.
		switch {
		case len(orgs) > 1:
			subLogger.Debug().
				Int("retrieved_plans", len(syncPlans)).
				Int("request", requestNum).
				Int("requests_remaining", requestsRemaining).
				Str("runtime_request", time.Since(retrievalStart).String()).
				Str("runtime_elapsed", time.Since(funcTimeStart).String()).
				Msg("Finished sync plans retrieval for this organization")
		default:
			subLogger.Debug().
				Int("retrieved_plans", len(syncPlans)).
				Msg("Finished sync plans retrieval for this organization")
		}

		allSyncPlans = append(allSyncPlans, syncPlans...)
	}

	logger.Debug().
		Str("runtime_total", time.Since(funcTimeStart).String()).
		Msg("Completed sync plans retrieval for all requested organizations")

	return allSyncPlans, nil
}

// IsOKState indicates whether any problems have been identified with this
// sync plan.
func (sp SyncPlan) IsOKState() bool {
	switch {
	case sp.IsStuck():
		return false

	// NOTE: While stuck plans are the current focus we may wish to expand the
	// list of problem "symptoms" (i.e., use additional case statements) to
	// include other attributes in the future.

	default:
		return true
	}
}

// IsStuck indicates whether (after any applied grace time) the sync plan is
// considered to be in a "stuck" state (Next Sync state set to past date/time).
//
// Grace time is applied to help prevent flagging a sync plan that is
// "spinning up" or in a temporary pending status (e.g., on a busy system) as
// problematic.
//
// NOTE: Very busy systems keeping sync plans in a pending state for an
// extended duration are still likely to be flagged as non-OK by current
// logic.
func (sp SyncPlan) IsStuck() bool {
	now := time.Now().UTC()
	nextSync := time.Time(sp.NextSync).UTC()

	switch {
	case sp.Enabled && nextSync.Before(now):
		diff := now.Sub(nextSync).Minutes()

		if diff <= syncTimeGraceMinutes {
			return false
		}

		return true

	default:
		return false
	}
}

// DaysStuck indicates how many days the sync plan has been in a "stuck"
// state.
func (sp SyncPlan) DaysStuck() int {
	switch {
	case !sp.Enabled:
		// Disabled sync plans are not considered "stuck" as they have been
		// turned off a sysadmin.
		return 0

	case time.Time(sp.NextSync).IsZero():

		// Use creation date of the plan instead of the time zero value.
		timeSinceStuck := time.Since(time.Time(sp.OriginalSyncDate)).Hours()

		// Toss remainder so that we only get the whole number of days
		daysStuck := int(math.Trunc(timeSinceStuck / 24))
		if daysStuck < 0 {
			daysStuck = 0
		}

		return daysStuck

	default:
		timeSinceStuck := time.Since(time.Time(sp.NextSync)).Hours()

		// Toss remainder so that we only get the whole number of days
		daysStuck := int(math.Trunc(timeSinceStuck / 24))
		if daysStuck < 0 {
			daysStuck = 0
		}

		return daysStuck
	}
}

// DaysStuckHR provides a human readable indication of how many days in the
// past the sync plan has been in a "stuck" state.
func (sp SyncPlan) DaysStuckHR() string {
	if sp.IsOKState() {
		return "N/A"
	}

	if sp.DaysStuck() == 0 {
		return "<1d"
	}

	return strconv.Itoa(sp.DaysStuck())
}

// NextSyncTime provides a display friendly version of the next scheduled sync
// time for the sync plan.
func (sp SyncPlan) NextSyncTime() string {
	if time.Time(sp.NextSync).IsZero() {
		return "N/A"
	}

	return sp.NextSync.String()
}

// Total provides the number of sync plans in the collection.
func (sps SyncPlans) Total() int {
	return len(sps)
}

// NumEnabled provides the number of sync plans in the collection in an
// enabled state.
func (sps SyncPlans) NumEnabled() int {
	var num int

	for _, syncPlan := range sps {
		if syncPlan.Enabled {
			num++
		}
	}

	return num
}

// NumDisabled provides the number of sync plans in the collection in a
// disabled state.
func (sps SyncPlans) NumDisabled() int {
	var num int

	for _, syncPlan := range sps {
		if !syncPlan.Enabled {
			num++
		}
	}

	return num
}

// NumStuck indicates the number of sync plans in the collection are in a
// "stuck" state.
func (sps SyncPlans) NumStuck() int {
	var num int

	for _, syncPlan := range sps {
		if syncPlan.IsStuck() {
			num++
		}
	}

	return num
}

// NumProblemPlans returns the total number of sync plans with a non-OK state.
func (sps SyncPlans) NumProblemPlans() int {
	// NOTE: While stuck plans are the current focus we may wish to expand the
	// list of problem "symptoms" to include other attributes in the future.
	// This method provides a more generic "are there any problems" status
	// check to cover that possibility.
	return sps.NumStuck()
}

// IsOKState indicates whether any problems have been identified with the sync
// plans in this collection.
func (sps SyncPlans) IsOKState() bool {
	for _, syncPlan := range sps {
		if !syncPlan.IsOKState() {
			return false
		}
	}

	return true
}

// Enabled returns a new collection containing all sync plans from the
// original collection which are in an enabled state.
func (sps SyncPlans) Enabled() SyncPlans {
	matches := make(SyncPlans, 0, sps.NumEnabled())

	for _, syncPlan := range sps {
		if syncPlan.Enabled {
			matches = append(matches, syncPlan)
		}
	}

	return matches
}

// Disabled returns a new collection containing all sync plans from the
// original collection which are not in an enabled state.
func (sps SyncPlans) Disabled() SyncPlans {
	matches := make(SyncPlans, 0, sps.NumDisabled())

	for _, syncPlan := range sps {
		if !syncPlan.Enabled {
			matches = append(matches, syncPlan)
		}
	}

	return matches
}

// Stuck returns a new collection containing all sync plans from the original
// collection which are in a "stuck" state.
func (sps SyncPlans) Stuck() SyncPlans {
	matches := make(SyncPlans, 0, sps.NumStuck())
	now := time.Now()

	for _, syncPlan := range sps {
		if syncPlan.Enabled && time.Time(syncPlan.NextSync).Before(now) {
			matches = append(matches, syncPlan)
		}
	}

	return matches
}

// getOrgSyncPlans retrieves all sync plans for the given organization.
func getOrgSyncPlans(ctx context.Context, client *APIClient, org Organization) (SyncPlans, error) {
	subLogger := client.Logger.With().
		Int("org_id", org.ID).
		Str("org_name", org.Name).
		Logger()

	apiURL := fmt.Sprintf(
		SyncPlansAPIEndPointURLTemplate,
		client.AuthInfo.Server,
		client.AuthInfo.Port,
		org.ID,
		client.Limits.PerPage,
	)

	subLogger.Debug().Msg("Preparing request to retrieve sync plans")
	request, reqErr := prepareRequest(ctx, client, apiURL)
	if reqErr != nil {
		return nil, reqErr
	}

	subLogger.Debug().Msg("Submitting HTTP request")
	response, respErr := client.Do(request)
	if respErr != nil {
		return nil, respErr
	}
	subLogger.Debug().Msg("Successfully submitted HTTP request")

	// Make sure that we close the response body once we're done with it
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			subLogger.Error().Err(closeErr).Msg("error closing response body")
		}
	}()

	// Evaluate the response
	validateErr := validateResponse(ctx, response, subLogger, client.AuthInfo.ReadLimit)
	if validateErr != nil {
		return nil, validateErr
	}

	subLogger.Debug().Msg("Successfully validated HTTP response")

	subLogger.Debug().Msgf(
		"Decoding JSON data from %q using a limit of %d bytes",
		apiURL,
		client.AuthInfo.ReadLimit,
	)

	var syncPlansQueryResp SyncPlansResponse
	decodeErr := decode(&syncPlansQueryResp, response.Body, subLogger, apiURL, client.AuthInfo.ReadLimit)
	if decodeErr != nil {
		return nil, decodeErr
	}

	// Annotate Sync Plans with specific Org values for convenience.
	for i := range syncPlansQueryResp.SyncPlans {
		syncPlansQueryResp.SyncPlans[i].OrganizationName = org.Name
		syncPlansQueryResp.SyncPlans[i].OrganizationLabel = org.Label
		syncPlansQueryResp.SyncPlans[i].OrganizationTitle = org.Title
	}

	subLogger.Debug().
		Str("api_endpoint", apiURL).
		Msg("Successfully decoded JSON data")

	return syncPlansQueryResp.SyncPlans, nil

}
