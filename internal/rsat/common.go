// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package rsat

// requestsCounterFunc is a helper function used to track the current request
// number and the requests remaining for a collection.
type requestsCounterFunc func() (int, int)

// newRequestsCounter creates a new requests counter function using the given
// value as the starting value. When called, the requests counter function
// will return the requests issued thus far and the requests remaining.
//
// For example, if you call newRequestsCounter(20) you will get back a
// function that returns two values. The first time you call this function it
// will return the values 1 and 19.
func newRequestsCounter(start int) requestsCounterFunc {
	remaining := start
	issued := 0

	return func() (int, int) {
		if remaining > 0 {
			remaining--
			issued++
		}

		return issued, remaining
	}
}
