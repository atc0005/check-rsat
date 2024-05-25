// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package rsat

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/atc0005/go-nagios"
)

// OrganizationsResponse represents the API response from a request for all
// organizations in the Red Hat Satellite server.
//
// https://access.redhat.com/documentation/en-us/red_hat_satellite/6.5/html-single/api_guide/index#sect-API_Guide-Understanding_the_JSON_Response_Format
// https://access.redhat.com/documentation/en-us/red_hat_satellite/6.15/html-single/api_guide/index#sect-API_Guide-Understanding_the_JSON_Response_Format
type OrganizationsResponse struct {
	// Organizations is the collection of Organizations returned in the API
	// query response.
	Organizations []Organization `json:"results"`

	// Search is the search string based on scoped_scoped syntax.
	Search NullString `json:"search"`

	// Sort is the optional sorting criteria for API query responses.
	Sort SortOptions `json:"sort"`

	// Subtotal is the number of objects returned with the given search
	// parameters. If there is no search, then subtotal is equal to total.
	Subtotal int `json:"subtotal"`

	// Total is the total number of objects without any search parameters.
	Total int `json:"total"`

	// Page is the page number for the current query response results.
	//
	// NOTE: In practice, this value has been found to be  returned as an
	// integer in the first response and as a string value for each additional
	// page of results. The json.Number type accepts either format when
	// decoding the response.
	Page json.Number `json:"page"`

	// PerPage is the pagination limit applied to API query results. If not
	// specified by the client this is the default value set by the API.
	PerPage int `json:"per_page"`
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

	allOrgs := make([]Organization, 0, client.Limits.PerPage*2)

	apiURLQueryParams := make(map[string]string)
	apiURLQueryParams[APIEndpointURLQueryParamFullResultKey] = APIEndpointURLQueryParamFullResultDefaultValue
	apiURLQueryParams[APIEndpointURLQueryParamPerPageKey] = strconv.Itoa(client.Limits.PerPage)

	var nextPage int
	remainingOrgs := true

	for remainingOrgs {
		logger.Debug().
			Msg("Collecting organizations from the API")

		nextPage++
		apiURLQueryParams[APIEndpointURLQueryParamPageKey] = strconv.Itoa(nextPage)

		response, respErr := submitAPIQueryRequest(ctx, client, apiURL, apiURLQueryParams, logger)
		if respErr != nil {
			return nil, respErr
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

		// Close the response body once we're done with it. We explicitly
		// close here vs deferring via closure to prevent accumulating client
		// connections to the API if we need to perform multiple paged
		// requests.
		if closeErr := response.Body.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("error closing response body")
		}

		allOrgs = append(allOrgs, orgsQueryResp.Organizations...)

		numNewOrgs := len(orgsQueryResp.Organizations)
		numCollectedOrgs := len(allOrgs)
		numOrgsRemaining := orgsQueryResp.Subtotal - numCollectedOrgs

		logger.Debug().
			Str("api_endpoint", apiURL).
			Int("orgs_collected", numCollectedOrgs).
			Int("orgs_new", numNewOrgs).
			Int("orgs_remaining", numOrgsRemaining).
			Msg("Added decoded organizations to collection")

		logger.Debug().
			Msg("Determining if we have collected all organizations from the API")

		remainingOrgs = numOrgsRemaining != 0
	}

	logger.Debug().
		Str("runtime_total", time.Since(funcTimeStart).String()).
		Msg("Completed retrieval of all organizations")

	return allOrgs, nil
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
