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
	"sort"
	"time"

	"github.com/atc0005/go-nagios"
)

// OrganizationsResponse represents the API response from a request for all
// organizations in the Red Hat Satellite server.
type OrganizationsResponse struct {
	Organizations []Organization `json:"results"`
	Search        NullString     `json:"search"`
	Sort          SortOptions    `json:"sort"`
	Subtotal      int            `json:"subtotal"`
	Total         int            `json:"total"`
	Page          int            `json:"page"`
	PerPage       int            `json:"per_page"`
}

// Organization is an isolated collection of systems, content, and other
// functionality within a Red Hat Satellite deployment.
type Organization struct {
	CreatedAt   StandardAPITime `json:"created_at"`
	UpdatedAt   StandardAPITime `json:"updated_at"`
	Description NullString      `json:"description"`
	Label       string          `json:"label"`
	Name        string          `json:"name"`
	Title       string          `json:"title"`
	SyncPlans   SyncPlans       `json:"-"`
	// Products    Products        `json:"-"`
	// Hosts       Hosts           `json:"-"`
	ID int `json:"id"`
}

// Organizations is a collection of Red Hat Satellite organizations.
type Organizations []Organization

// GetOrganizations uses the given client to retrieve all Red Hat Satellite
// organizations.
func GetOrganizations(ctx context.Context, client *APIClient) ([]Organization, error) {
	funcTimeStart := time.Now()

	if client == nil {
		return nil, fmt.Errorf(
			"required API client was not provided: %w",
			ErrMissingValue,
		)
	}

	logger := client.Logger

	apiURL := fmt.Sprintf(
		OrganizationsAPIEndPointURLTemplate,
		client.AuthInfo.Server,
		client.AuthInfo.Port,
	)

	request, err := prepareRequest(ctx, client, apiURL)
	if err != nil {
		return nil, err
	}

	logger.Debug().Msg("Submitting HTTP request")

	response, respErr := client.Do(request)
	if respErr != nil {
		return nil, respErr
	}

	logger.Debug().Msg("Successfully submitted HTTP request")

	// Make sure that we close the response body once we're done with it
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msgf("error closing response body")
		}
	}()

	// Evaluate the response
	validateErr := validateResponse(ctx, response, logger, client.AuthInfo.ReadLimit)
	if validateErr != nil {
		return nil, err
	}

	logger.Debug().Msgf(
		"Decoding JSON data from %q using a limit of %d bytes",
		apiURL,
		client.AuthInfo.ReadLimit,
	)

	var orgsQueryResp OrganizationsResponse
	decodeErr := decode(&orgsQueryResp, response.Body, logger, apiURL, client.AuthInfo.ReadLimit)
	if decodeErr != nil {
		return nil, decodeErr
	}

	logger.Debug().
		Str("api_endpoint", apiURL).
		Msg("Successfully decoded JSON data")

	logger.Debug().
		Str("runtime_total", time.Since(funcTimeStart).String()).
		Msg("Completed retrieval of all organizations")

	return orgsQueryResp.Organizations, nil
}

// Sort sorts the Organizations in the collection by name.
func (orgs Organizations) Sort() {
	sort.SliceStable(orgs, func(i int, j int) bool {
		return orgs[i].Name < orgs[j].Name
	})
}

// GetOrgsWithSyncPlans uses the provided API client to retrieve all Red Hat
// Satellite organizations along with their sync plans.
func GetOrgsWithSyncPlans(ctx context.Context, client *APIClient) (Organizations, error) {
	funcTimeStart := time.Now()

	if client == nil {
		return nil, fmt.Errorf(
			"required API client was not provided: %w",
			ErrMissingValue,
		)
	}

	logger := client.Logger

	logger.Debug().Msg("Retrieving organizations")

	orgs, orgsErr := GetOrganizations(ctx, client)
	if orgsErr != nil {
		logger.Error().Err(orgsErr).Msg("Failed to retrieve organizations")
		return nil, fmt.Errorf(
			"failed to retrieve organizations: %w",
			orgsErr,
		)
	}

	logger.Debug().Msg("Successfully retrieved organizations")

	reqsCounter := newRequestsCounter(len(orgs))

	// Update all organizations with retrieved sync plans.
	for i := range orgs {

		subLogger := logger.With().
			Int("org_id", orgs[i].ID).
			Str("org_name", orgs[i].Name).
			Stack().Logger()

		retrievalStart := time.Now()

		subLogger.Debug().Msg("Retrieving sync plans for organization")

		syncPlans, syncPlansErr := GetSyncPlans(ctx, client, orgs[i])
		if syncPlansErr != nil {
			subLogger.Error().Err(syncPlansErr).Msg("Failed to retrieve sync plans")
			return nil, fmt.Errorf(
				"failed to retrieve sync plans for organization"+
					" (name: %s, id: %d) %w",
				orgs[i].Name,
				orgs[i].ID,
				syncPlansErr,
			)
		}

		requestNum, requestsRemaining := reqsCounter()

		subLogger.Debug().
			Int("retrieved_plans", len(syncPlans)).
			Int("request", requestNum).
			Int("requests_remaining", requestsRemaining).
			Str("runtime_request", time.Since(retrievalStart).String()).
			Str("runtime_elapsed", time.Since(funcTimeStart).String()).
			Msg("Finished sync plans retrieval for this organization")

		orgs[i].SyncPlans = syncPlans
	}

	logger.Debug().Msg("Successfully retrieved sync plans for all organizations")

	return orgs, nil
}

// NumOrgs returns the number of organizations in the collection.
func (orgs Organizations) NumOrgs() int {
	return len(orgs)
}

// NumPlans returns the number of sync plans for all organizations in the
// collection.
func (orgs Organizations) NumPlans() int {
	var num int
	for _, org := range orgs {
		num += len(org.SyncPlans)
	}

	return num
}

// NumPlansEnabled returns the total number of sync plans for all
// organizations in the collection with enabled state.
func (orgs Organizations) NumPlansEnabled() int {
	var num int

	for _, org := range orgs {
		for _, syncPlan := range org.SyncPlans {
			if syncPlan.Enabled {
				num++
			}
		}
	}

	return num
}

// NumPlansStuck returns the total number of sync plans for all organizations
// in the collection with Next Sync state set to past date/time.
func (orgs Organizations) NumPlansStuck() int {
	var num int

	for _, org := range orgs {
		num += org.SyncPlans.NumStuck()
	}

	return num
}

// NumPlansDisabled returns the total number of sync plans for all
// organizations in the collection with disabled state.
func (orgs Organizations) NumPlansDisabled() int {
	var num int

	for _, org := range orgs {
		num += org.SyncPlans.NumDisabled()
	}

	return num
}

// NumProblemPlans returns the total number of sync plans for all
// organizations in the collection with a non-OK state.
func (orgs Organizations) NumProblemPlans() int {
	// NOTE: While stuck plans are the current focus we may wish to expand the
	// list of problem "symptoms" to include other attributes in the future.
	// This method provides a more generic "are there any problems" status
	// check to cover that possibility.
	return orgs.NumPlansStuck()
}

// IsOKState indicates whether all items in the collection were evaluated to
// an OK state.
func (orgs Organizations) IsOKState() bool {
	// return orgs.NumProblemPlans() == 0

	// The scope is a higher level than just whether there are problematic
	// sync plans (e.g., the Org might have problematic subscriptions that we
	// can alert on in the future).
	return !orgs.HasWarningState() && !orgs.HasCriticalState()
}

// HasCriticalState indicates whether any items in the collection were
// evaluated to a CRITICAL state.
func (orgs Organizations) HasCriticalState() bool {
	// TODO: Add support for performing threshold check to determine how many
	// days in the past a sync plan has been stuck. If greater than given
	// threshold indicate CRITICAL state.
	return false
}

// HasWarningState indicates whether any items in the collection were
// evaluated to a WARNING state.
func (orgs Organizations) HasWarningState() bool {
	return !orgs.HasCriticalState() && orgs.NumProblemPlans() > 0
}

// ServiceState returns the appropriate Service Check Status label and exit
// code for the collection's evaluation results.
func (orgs Organizations) ServiceState() nagios.ServiceState {
	var stateLabel string
	var stateExitCode int

	switch {
	case orgs.HasCriticalState():
		stateLabel = nagios.StateCRITICALLabel
		stateExitCode = nagios.StateCRITICALExitCode
	case orgs.HasWarningState():
		stateLabel = nagios.StateWARNINGLabel
		stateExitCode = nagios.StateWARNINGExitCode
	case orgs.IsOKState():
		stateLabel = nagios.StateOKLabel
		stateExitCode = nagios.StateOKExitCode
	default:
		stateLabel = nagios.StateUNKNOWNLabel
		stateExitCode = nagios.StateUNKNOWNExitCode
	}

	return nagios.ServiceState{
		Label:    stateLabel,
		ExitCode: stateExitCode,
	}
}
