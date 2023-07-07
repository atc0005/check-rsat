// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import "github.com/atc0005/go-nagios"

// annotateError is a helper function used to add additional human-readable
// explanation for errors encountered during plugin execution. We first apply
// common advice for more general errors then apply advice specific to errors
// routinely encountered by this specific project.
func annotateErrors(plugin *nagios.Plugin) {
	// If nothing to process, skip setup/processing steps.
	if len(plugin.Errors) == 0 {
		return
	}

	// Start off with the default advice collection.
	errorAdviceMap := nagios.DefaultErrorAnnotationMappings()

	// FIXME: Annotate errors related to TLS renegotiation not being enabled
	// for plugin but requested for server.

	// Override specific error with project-specific feedback.
	// errorAdviceMap[syscall.ECONNRESET] = connectionResetByPeerAdvice

	// Apply error advice annotations.
	plugin.AnnotateRecordedErrors(errorAdviceMap)
}
