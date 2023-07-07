// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package rsat

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/atc0005/check-rsat/internal/netutils"
	"github.com/rs/zerolog"
)

// APIClient represents a customized HTTP client for interacting with Red
// Hat Satellite API endpoints.
type APIClient struct {
	*http.Client
	AuthInfo APIAuthInfo
	Logger   zerolog.Logger
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
func NewAPIClient(apiAuthInfo APIAuthInfo, logger zerolog.Logger) *APIClient {
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
	}
}
