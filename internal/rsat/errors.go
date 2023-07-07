// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package rsat

import (
	"errors"
	"fmt"
)

// FIXME: Should we consistently use the PrepError type instead of using these
// sentinel errors?
var (
	// ErrMissingValue indicates that an expected value was missing.
	ErrMissingValue = errors.New("missing expected value")

	// ErrHTTPResponseOutsideRange indicates that a response was received
	// which falls outside of an acceptable range.
	ErrHTTPResponseOutsideRange = errors.New("response is outside acceptable range")

	// ErrJSONUnexpectedObjectCount indicates that a response was received
	// with more provided JSON objects than expected.
	ErrJSONUnexpectedObjectCount = errors.New("unexpected JSON object count")

	// ErrJSONDecodeFailure = errors.New("")

	// ErrOrgsRetrievalFailed = errors.New("failed to retrieve organizations")
)

// PrepError represents a class of errors encountered while performing tasks
// related to preparing a components Set.
type PrepError struct {

	// Step indicates the specific prep task which failed.
	//
	// NOTE: Constants should be used to make comparisons more reliable.
	Task string

	// Message provides additional (brief) context describing why the error
	// occurred.
	//
	// e.g., "error parsing URL" or "error preparing request for URL"
	Message string

	// Source associated with the prep task.
	//
	// e.g.,
	// "https://rsat.example.com/katello/api/v2/organizations/27/subscriptions"
	Source string

	// Cause is the underlying error which occurred while performing a task as
	// part of preparing a components set. This error is "bundled" for later
	// evaluation.
	Cause error
}

// Error provides a human readable explanation for a components Set
// preparation task failure.
func (s *PrepError) Error() string {
	return fmt.Sprintf(
		"task: %q: %s: source: %s cause: %v",
		s.Task,
		s.Message,
		s.Source,
		s.Cause,
	)
}

// Is supports error wrapping by indicating whether a given error matches the
// specific failed task associated with this error.
func (s *PrepError) Is(target error) bool {
	t, ok := target.(*PrepError)
	if !ok {
		return false
	}

	return t.Task == s.Task
}

// Unwrap supports error wrapping by returning the enclosed error associated
// with the specific failed task  was encountered as part of preparing a components Set.
func (s *PrepError) Unwrap() error {
	return s.Cause
}
