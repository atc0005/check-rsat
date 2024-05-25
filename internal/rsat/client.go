// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package rsat

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/atc0005/check-rsat/internal/netutils"
	"github.com/rs/zerolog"
)

// APILimits represents the settings used to comply with the limits set by an
// API endpoint.
type APILimits struct {
	PerPage int
}

// APIClient represents a customized HTTP client for interacting with Red
// Hat Satellite API endpoints.
type APIClient struct {
	*http.Client
	AuthInfo APIAuthInfo
	Logger   zerolog.Logger
	Limits   APILimits
	// APIResponseCache CachedAPIResponses
}

// CachedAPIResponses represents specific API responses which are cached to
// reduce overhead of frequent access.
//
// type CachedAPIResponses struct {
// 	orgs OrganizationsResponse
// }

func getCustomTLSConfig(apiAuthInfo APIAuthInfo) *tls.Config {
	// https://www.golinuxcloud.com/golang-http/#Create_HTTPS_client
	// https://www.golinuxcloud.com/golang-http/#Create_TLS_Config
	var tlsConfig *tls.Config

	// Apply minimal relaxation of TLS renegotiation settings if sysadmin
	// requested that the Red Hat Satellite be permitted to request
	// renegotiation.
	tlsRenegotiation := func() tls.RenegotiationSupport {
		if apiAuthInfo.PermitTLSRenegotiation {
			return tls.RenegotiateOnceAsClient
		}
		return tls.RenegotiateNever
	}()

	switch {
	case apiAuthInfo.CACert != nil:
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(apiAuthInfo.CACert)

		tlsConfig = &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: apiAuthInfo.TrustCert, // nolint:gosec
			Renegotiation:      tlsRenegotiation,
		}

	default:
		tlsConfig = &tls.Config{
			InsecureSkipVerify: apiAuthInfo.TrustCert, // nolint:gosec
			Renegotiation:      tlsRenegotiation,
		}
	}

	return tlsConfig
}

// NewAPIClient uses the provided API Auth details to construct a custom HTTP
// client used to interact with
func NewAPIClient(apiAuthInfo APIAuthInfo, apiLimits APILimits, logger zerolog.Logger) *APIClient {
	tlsConfig := getCustomTLSConfig(apiAuthInfo)

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		MaxIdleConns:    1,                // TODO: Allow adjusting this via config package
		IdleConnTimeout: 30 * time.Second, // TODO: Allow adjusting this via config package
		DialContext: netutils.DialerWithContext(
			apiAuthInfo.NetworkType,
			logger,
		),
	}

	c := &http.Client{
		Transport: transport,
	}

	return &APIClient{
		Client:   c,
		AuthInfo: apiAuthInfo,
		Logger:   logger,
		Limits:   apiLimits,
	}
}

// submitAPIQueryRequest is a helper function used to submit a request to an
// API endpoint and perform basic validation of the results.
//
// TODO: Refactor to be an APIClient method
func submitAPIQueryRequest(
	ctx context.Context,
	client *APIClient,
	apiURL string,
	apiURLQueryParams map[string]string,
	logger zerolog.Logger,
) (*http.Response, error) {

	logger.Debug().Msg("Preparing request for API query")
	request, reqErr := prepareRequest(ctx, client, apiURL, apiURLQueryParams)
	if reqErr != nil {
		return nil, reqErr
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
			logger.Error().Err(closeErr).Msg("error closing response body")
		}
	}()

	// Evaluate the response
	validateErr := validateResponse(ctx, response, logger, client.AuthInfo.ReadLimit)
	if validateErr != nil {
		return nil, validateErr
	}

	logger.Debug().Msg("Successfully validated HTTP response")

	return response, nil
}
