// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

const myAppName string = "check-rsat"
const myAppURL string = "https://github.com/atc0005/check-rsat"

// ExitCodeCatchall indicates a general or miscellaneous error has occurred.
// This exit code is not directly used by monitoring plugins in this project.
// See https://tldp.org/LDP/abs/html/exitcodes.html for additional details.
const ExitCodeCatchall int = 1

// Shared flags help text.
const (
	helpFlagHelp                   string = "Emit this help text"
	versionFlagHelp                string = "Whether to display application version and then immediately exit application."
	logLevelFlagHelp               string = "Sets log level."
	brandingFlagHelp               string = "Toggles emission of branding details with plugin status details. This output is disabled by default."
	trustCertFlagHelp              string = "Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option."
	serverFlagHelp                 string = "The Red Hat Satellite server FQDN or IP Address."
	usernameFlagHelp               string = "The valid user for the given Red Hat Satellite server."
	passwordFlagHelp               string = "The valid password for the specified user." //nolint:gosec
	tcpPortFlagHelp                string = "The port used by the Red Hat Satellite server API."
	networkTypeFlagHelp            string = "Limits network connections to one of tcp4 (IPv4-only), tcp6 (IPv6-only) or auto (either)."
	caCertificateFlagHelp          string = "CA Certificate used to validate the certificate chain used by the Red Hat Satellite server."
	permitTLSRenegotiationFlagHelp string = "Whether support for accepting renegotiation requests from the Red Hat Satellite server are permitted. This support is disabled by default. Renegotiation is not supported for TLS 1.3."
	omitOKSyncPlansHelp            string = "Whether sync plans listed in plugin output should be limited to just those in a non-OK state."
	verboseFlagHelp                string = "Whether to display verbose details in the final plugin output."
)

// CLI App flags help text.
const (
	cliAppTimeoutFlagHelp         string = "Timeout value in seconds before application execution is abandoned and an error returned."
	inspectorOutputFormatFlagHelp string = "Sets output format."
)

// Plugin flags help text.
const (
	readLimitFlagHelp     string = "Limit in bytes used to help prevent abuse when reading input that could be larger than expected."
	pluginTimeoutFlagHelp string = "Timeout value in seconds before plugin execution is abandoned and an error returned."
)

// Flag names for consistent references. Exported so that they're available
// from tests.
const (
	HelpFlagLong                   string = "help"
	HelpFlagShort                  string = "h"
	VersionFlagLong                string = "version"
	VerboseFlagLong                string = "verbose"
	BrandingFlag                   string = "branding"
	TrustCertFlagLong              string = "trust-cert"
	TimeoutFlagLong                string = "timeout"
	TimeoutFlagShort               string = "t"
	ReadLimitFlagLong              string = "read-limit"
	LogLevelFlagLong               string = "log-level"
	LogLevelFlagShort              string = "ll"
	ServerFlagLong                 string = "server"
	UsernameFlagLong               string = "username"
	PasswordFlagLong               string = "password"
	PortFlagLong                   string = "port"
	NetTypeFlagLong                string = "net-type"
	CACertificateFlagLong          string = "ca-cert"
	PermitTLSRenegotiationFlagLong string = "permit-tls-renegotiation"
	OmitOKSyncPlansFlagLong        string = "omit-ok"
	InspectorOutputFormatFlagLong  string = "output-format"
)

// Default flag settings if not overridden by user input
const (
	defaultHelp                   bool   = false
	defaultLogLevel               string = "info"
	defaultVerbose                bool   = false
	defaultEmitBranding           bool   = false
	defaultDisplayVersionAndExit  bool   = false
	defaultTrustCert              bool   = false
	defaultPermitTLSRenegotiation bool   = false
	defaultOmitOKSyncPlans        bool   = false
	defaultServer                 string = ""
	defaultUsername               string = ""
	defaultPassword               string = ""
	defaultTCPPort                int    = 443
	defaultNetworkType            string = netTypeTCPAuto
	defaultCACertificate          string = ""

	// Red Hat Satellite API response times can be slow, so best to set a
	// generous default timeout.
	defaultCLIAppTimeout int = 300

	// Red Hat Satellite API response times can be slow, so best to set a
	// generous default timeout.
	defaultPluginTimeout int = 240

	// Set a read limit to help prevent abuse from unexpected/overly large
	// input. The limit set here is OVERLY generous and is unlikely to be met
	// unless something is broken.
	defaultReadLimit int64 = 1 * MB

	defaultInspectorOutputFormat string = InspectorOutputFormatPrettyTable
)

const (
	// netTypeTCPAuto is a custom keyword indicating that either of IPv4 or
	// IPv6 is an acceptable network type.
	netTypeTCPAuto string = "auto"

	// netTypeTCP4 indicates that IPv4 network connections are required.
	netTypeTCP4 string = "tcp4"

	// netTypeTCP6 indicates that IPv6 network connections are required
	netTypeTCP6 string = "tcp6"
)

const (
	appTypePlugin    string = "plugin"
	appTypeInspector string = "Inspector"
)

// MB represents 1 Megabyte
const MB int64 = 1048576

// Supported Inspector type application output formats
const (
	InspectorOutputFormatOverview    string = "overview"
	InspectorOutputFormatPrettyTable string = "pretty-table"
	InspectorOutputFormatSimpleTable string = "simple-table"
	InspectorOutputFormatVerbose     string = "verbose"
)
