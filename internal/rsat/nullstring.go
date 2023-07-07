// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package rsat

import (
	"encoding/json"
	"strings"
)

// NullString represents a string value that may potentially be null in the
// input JSON feed.
type NullString string

// MarshalJSON implements the json.Marshaler interface. This compliments the
// custom Unmarshaler implementation to handle potentially null string field
// values in place of using an empty interface.
func (ns NullString) MarshalJSON() ([]byte, error) {

	if len(string(ns)) == 0 {
		return []byte(JSONNullKeyword), nil
	}

	// NOTE: If we fail to convert the type, an infinite loop will occur.
	return json.Marshal(string(ns))

}

// UnmarshalJSON implements the json.Unmarshaler interface to handle
// potentially null string field values in place of using an empty interface.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == JSONNullKeyword {
		*ns = ""
		return nil
	}

	*ns = NullString(strings.Trim(str, "\""))

	return nil
}
