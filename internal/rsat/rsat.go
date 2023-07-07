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
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog"
)

// JSONNullKeyword is the keyword used in JSON to represent null.
const JSONNullKeyword string = "null"

// Red Hat Satellite API endpoints URL templates.
const (
	// OrganizationsAPIEndPointURLTemplate provides a template for a fully
	// qualified API endpoint URL for retrieving Organizations from a Red Hat
	// Satellite instance.
	OrganizationsAPIEndPointURLTemplate string = "https://%s:%d/api/v2/organizations?full_result=1"

	// SubscriptionsAPIEndPointURLTemplate provides a template for a fully
	// qualified API endpoint URL for retrieving Subscriptions associated with
	// a Red Hat Satellite Organization.
	SubscriptionsAPIEndPointURLTemplate string = "https://%s:%d/katello/api/v2/organizations/%d/subscriptions?full_result=1"

	// SyncPlansAPIEndPointURLTemplate provides a template for a fully
	// qualified API endpoint URL for retrieving Sync Plans associated with a
	// Red Hat Satellite Organization.
	SyncPlansAPIEndPointURLTemplate string = "https://%s:%d/katello/api/v2/organizations/%d/sync_plans?full_result=1"

	// ProductsAPIEndPointURLTemplate provides a template for a fully
	// qualified API endpoint URL for retrieving Products associated with a
	// Red Hat Satellite Organization.
	ProductsAPIEndPointURLTemplate string = "https://%s:%d/katello/api/v2/products?organization_id=%d&full_result=1"
)

// Prep tasks for processing of Red Hat Satellite API endpoints.
const (
	PrepTaskParseURL         string = "parse URL"
	PrepTaskPrepareRequest   string = "prepare request"
	PrepTaskDecode           string = "decode JSON data"
	PrepTaskSubmitRequest    string = "submit request"
	PrepTaskValidateResponse string = "validate response"
)

// APIAuthInfo represents the settings needed to access Red Hat Satellite
// server API endpoints.
type APIAuthInfo struct {
	// ReadLimit is the size in bytes used to help prevent abuse when reading
	// input that could be larger than expected.
	ReadLimit int64

	// Port is the TCP/IP port associated with the Red Hat Satellite server's
	// API endpoints.
	Port int

	// Server is the FQDN or IP Address of the Red Hat Satellite server.
	Server string

	// Username is the valid user for the specified Red Hat Satellite server.
	Username string

	// Password is the valid password for the specified Red Hat Satellite
	// Server user account.
	Password string

	// UserAgent is an optional custom user agent string used to override the
	// default Go user agent ("Go-http-client/1.1").
	UserAgent string

	// NetworkType indicates whether an attempt should be made to connect to
	// only IPv4, only IPv6 or SSH servers listening on either of IPv4 or IPv6
	// addresses ("auto").
	NetworkType string

	// CACert is the optional certificate authority certificate used to
	// validate the certificate chain used by the Red Hat Satellite server.
	CACert []byte

	// PermitTLSRenegotiation controls whether the server is allowed to
	// request TLS renegotiation.
	PermitTLSRenegotiation bool

	// TrustCert indicates whether the certificate should be trusted as-is
	// without validation.
	TrustCert bool
}

// SortOptions is the optional sorting criteria for API query responses.
type SortOptions struct {
	By    NullString `json:"by"`
	Order NullString `json:"order"`
}

// decode is a helper function intended to handle the core JSON decoding tasks
// for various JSON sources (file, http body, etc.).
func decode(dst interface{}, reader io.Reader, logger zerolog.Logger, sourceName string, limit int64) error {
	if reader == nil {
		return &PrepError{
			Task:    PrepTaskDecode,
			Message: "failed to decode JSON data",
			Source:  sourceName,
			Cause: fmt.Errorf(
				"required JSON source was not provided: %w",
				ErrMissingValue,
			),
		}
	}

	logger.Debug().Msgf(
		"Setting up JSON decoder for source %s with a limit of %d bytes",
		sourceName,
		limit,
	)
	dec := json.NewDecoder(io.LimitReader(reader, limit))

	// This project does not use all fields from Red Hat Satellite API
	// responses so we do not attempt to assert that we've accounted for all
	// of them.
	logger.Debug().Msg("Allowing unknown JSON feed fields")

	logger.Debug().Msg("Decoding JSON input")

	// Decode the first JSON object.
	if err := dec.Decode(dst); err != nil {
		return &PrepError{
			Task:    PrepTaskDecode,
			Message: "failed to decode JSON data",
			Source:  sourceName,
			Cause:   err,
		}
	}
	logger.Debug().Msg("Successfully decoded JSON input")

	// If there is more than one object, something is off.
	if dec.More() {

		return &PrepError{
			Task:    PrepTaskDecode,
			Message: "failed to decode JSON data",
			Source:  sourceName,
			Cause: fmt.Errorf(
				"source %s contains multiple JSON objects; only one JSON object is supported: %w",
				sourceName,
				ErrJSONUnexpectedObjectCount,
			),
		}
	}

	return nil

}

// validateResponse is a helper function responsible for validating a response
// from an endpoint after submitting a message.
func validateResponse(ctx context.Context, response *http.Response, logger zerolog.Logger, limit int64) error {
	if response == nil {
		return &PrepError{
			Task:    PrepTaskValidateResponse,
			Message: "error validating HTTP request",
			Source:  "missing",
			Cause: fmt.Errorf(
				"required HTTP response was not provided: %w",
				ErrMissingValue,
			),
		}
	}

	feedSource := response.Request.URL.RequestURI()

	if err := ctx.Err(); err != nil {
		logger.Debug().Msg("context has expired")
		return &PrepError{
			Task:    PrepTaskValidateResponse,
			Message: "timeout reached",
			Source:  feedSource,
			Cause:   err,
		}
	}

	switch {
	case response.ContentLength == -1:
		logger.Debug().Msgf("Response indicates unknown length of content from %q", feedSource)
	default:
		logger.Debug().Msgf(
			"Response indicates %d bytes available to be read from %q",
			response.ContentLength,
			feedSource,
		)
	}

	// TODO: Refactor this block
	switch {

	// Successful / expected response.
	case response.StatusCode == http.StatusOK:
		logger.Debug().Msgf("Status code %d received as expected", response.StatusCode)

		return nil

	// Success status range, but not expected value.
	case response.StatusCode > 200 && response.StatusCode <= 299:
		logger.Debug().Msgf(
			"Status code %d (%s) received; expected %d (%s), but received value within success range",
			response.StatusCode,
			http.StatusText(response.StatusCode),
			http.StatusOK,
			http.StatusText(http.StatusOK),
		)

		return nil

	// Everything else is assumed to be an error (outside of success range).
	default:

		// Get the response body, then convert to string for use with extended
		// error messages
		responseData, readErr := io.ReadAll(io.LimitReader(response.Body, limit))
		if readErr != nil {
			return &PrepError{
				Task:    PrepTaskValidateResponse,
				Message: "error reading response data",
				Source:  feedSource,
				Cause:   readErr,
			}
		}
		responseString := string(responseData)

		statusCodeErr := fmt.Errorf(
			"response %v (%s) from API: %w",
			response.Status,
			responseString,
			ErrHTTPResponseOutsideRange,
		)

		return &PrepError{
			Task:    PrepTaskValidateResponse,
			Message: "unexpected response",
			Source:  feedSource,
			Cause:   statusCodeErr,
		}

	}

}

// prepareRequest is a helper function that prepares a http.Request (including
// all desired headers) for submission to an endpoint.
func prepareRequest(ctx context.Context, client *APIClient, apiURL string) (*http.Request, error) {
	if client == nil {
		return nil, &PrepError{
			Task:    PrepTaskPrepareRequest,
			Message: "error preparing HTTP request",
			Source:  apiURL,
			Cause: fmt.Errorf(
				"required API client was not provided: %w",
				ErrMissingValue,
			),
		}
	}

	if apiURL == "" {
		return nil, &PrepError{
			Task:    PrepTaskPrepareRequest,
			Message: "error preparing HTTP request",
			Source:  apiURL,
			Cause: fmt.Errorf(
				"required API URL was not provided: %w",
				ErrMissingValue,
			),
		}
	}

	logger := client.Logger

	logger.Debug().Msgf("Parsing %q as URL", apiURL)
	parsedURL, parseErr := url.Parse(apiURL)
	if parseErr != nil {
		return nil, &PrepError{
			Task:    PrepTaskParseURL,
			Message: "error parsing URL",
			Source:  apiURL,
			Cause:   parseErr,
		}
	}
	logger.Debug().Msgf("Successfully parsed %q as URL", apiURL)

	logger.Debug().Msg("Preparing HTTP request")
	request, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if reqErr != nil {
		return nil, &PrepError{
			Task:    PrepTaskPrepareRequest,
			Source:  parsedURL.String(),
			Message: "error preparing request for URL",
			Cause:   reqErr,
		}
	}

	// Explicitly note that we want JSON content.
	request.Header.Add("Content-Type", "application/json;charset=utf-8")

	// Provide API authentication credentials.
	// https://stackoverflow.com/questions/16673766/basic-http-auth-in-go
	request.SetBasicAuth(client.AuthInfo.Username, client.AuthInfo.Password)

	// If provided, override the default Go user agent ("Go-http-client/1.1")
	// with custom value.
	if client.AuthInfo.UserAgent != "" {
		logger.Debug().Msg("Setting custom user agent")
		request.Header.Set("User-Agent", client.AuthInfo.UserAgent)
	}

	return request, nil
}
