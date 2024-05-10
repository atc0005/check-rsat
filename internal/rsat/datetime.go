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
	"time"
)

// References:
//
// - https://romangaranin.net/posts/2021-02-19-json-time-and-golang/
// - https://pkg.go.dev/time#pkg-constants

// Time layouts observed "in the wild" for various versions of Red Hat
// Satellite.
const (
	// StandardAPITimeLayoutWithTimezone is the time layout format as used by
	// the Red Hat Satellite API for the majority of the date/time properties
	// when the user has their Satellite account timezone setting configured
	// as `(GMT+00:00) UTC`.
	//
	// Examples (from JSON API response) for Satellite 6.15:
	//
	// "created_at": "2024-05-09 21:14:51 UTC",
	// "updated_at": "2024-05-09 21:14:51 UTC",
	//
	// See also https://rsat.example.com/apidoc/v2/sync_plans/index.html
	StandardAPITimeLayoutWithTimezone string = "2006-01-02 15:04:05 UTC"

	// StandardAPITimeLayoutWithOffset is the time layout format as used by
	// the Red Hat Satellite API for the majority of the date/time properties
	// when the user has their Satellite account timezone setting configured
	// as `Browser timezone`.
	//
	// Examples (from JSON API response):
	//
	// "created_at": "2024-05-09 16:14:51 -0500",
	// "updated_at": "2024-05-09 16:14:51 -0500",
	//
	// See also https://rsat.example.com/apidoc/v2/sync_plans/index.html
	StandardAPITimeLayoutWithOffset string = "2006-01-02 15:04:05 -0700"

	// SyncTimeLayoutWithTimezone is the time layout format as used by current
	// versions of the Red Hat Satellite Sync Plans API for the next_sync
	// property in current versions of Red Hat Satellite when the user has
	// their Satellite account timezone setting configured as `(GMT+00:00)
	// UTC`.
	//
	// Example: "next_sync": "2024-05-10 17:14:00 UTC",
	//
	// See also https://rsat.example.com/apidoc/v2/sync_plans/index.html
	SyncTimeLayoutWithTimezone string = "2006-01-02 15:04:05 UTC"

	// SyncTimeLayoutWithOffset is the time layout format as used by current
	// versions of the Red Hat Satellite Sync Plans API for the next_sync
	// property when the user has their Satellite account timezone setting
	// configured as `Browser timezone`.
	//
	// Example: "next_sync": "2024/05/10 15:16:00 -0500",
	//
	// See also https://rsat.example.com/apidoc/v2/sync_plans/index.html
	SyncTimeLayoutWithOffset string = "2006-01-02 15:04:05 -0700"

	// LegacySyncTimeLayout is the time layout format as used by legacy
	// versions of the Red Hat Satellite Sync Plans API for the next_sync
	// property (e.g., Satellite 6.5).
	//
	// Example(account Timezone property is set to `(GMT+00:00) UTC`):
	//
	// "next_sync": "2024/05/10 20:16:00 +0000",
	//
	// Example(account Timezone property is set to `Browser timezone`):
	//
	// "next_sync": "2024/05/10 15:16:00 -0500",
	//
	// This layout works equally well for both.
	//
	// See also https://rsat.example.com/apidoc/v2/sync_plans/index.html
	LegacySyncTimeLayout string = "2006/01/02 15:04:05 -0700"
)

// StandardAPITime is time value as represented in the Red Hat Satellite API
// for the majority of the date/time properties. It uses the
// StandardAPITimeLayout format.
type StandardAPITime time.Time

// SyncTime is time value as represented in the Red Hat Satellite Sync Plans
// API for the next_sync and sync_date properties. It uses the SyncTimeLayout
// format.
type SyncTime time.Time

// String implements the fmt.Stringer interface as a convenience method.
func (dt StandardAPITime) String() string {
	return dt.Format(StandardAPITimeLayoutWithOffset)
}

// String implements the fmt.Stringer interface as a convenience method.
func (dt SyncTime) String() string {
	// return dt.Format(StandardAPITimeLayout)
	switch {
	case time.Time(dt).IsZero():
		return "Not scheduled"
	default:
		return time.Time(dt).Local().Format(StandardAPITimeLayoutWithOffset)
	}
}

// Format calls (time.Time).Format as a convenience for the caller.
func (dt StandardAPITime) Format(layout string) string {
	return time.Time(dt).Format(layout)
}

// Format calls (time.Time).Format as a convenience for the caller.
func (dt SyncTime) Format(layout string) string {
	return time.Time(dt).Format(layout)
}

// MarshalJSON implements the json.Marshaler interface. This compliments the
// custom Unmarshaler implementation to handle conversion of a native Go
// time.Time format to a time value in a format that matches JSON API
// expectations.
func (dt StandardAPITime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dt).Format(StandardAPITimeLayoutWithOffset))
}

// MarshalJSON implements the json.Marshaler interface. This compliments the
// custom Unmarshaler implementation to handle conversion of a native Go
// time.Time format to a time value in a format that matches JSON API
// expectations.
func (dt SyncTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dt).Format(SyncTimeLayoutWithOffset))
}

// UnmarshalJSON implements the json.Unmarshaler interface to handle
// converting a time string from the JSON API (most time properties) to a
// native Go time.Time value using a supported auto-detected layout.
func (dt *StandardAPITime) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), `"`) // get rid of "
	if value == "" || value == JSONNullKeyword {

		// Per json.Unmarshaler convention we treat "null" value as a no-op.
		return nil
	}

	t, err := parseDate(value)
	if err != nil {
		return err
	}

	*dt = StandardAPITime(t) // set result using the pointer

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface to handle
// converting a time string from the next_sync property in the JSON API to a
// native Go time.Time value using the SyncTimeLayout format.
func (dt *SyncTime) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), `"`) // get rid of "
	if value == "" || value == JSONNullKeyword {

		// Per json.Unmarshaler convention we treat "null" value as a no-op.
		return nil
	}

	t, err := parseDate(value)
	if err != nil {
		return err
	}

	*dt = SyncTime(t) // set result using the pointer

	return nil
}

// parseDate is a helper function that attempts to handle all known datetime
// formats for legacy and current Red Hat Satellite APIs. An error is returned
// if the given datetime string does not match a known layout.
func parseDate(datetime string) (time.Time, error) {
	knownLayouts := []string{
		StandardAPITimeLayoutWithTimezone,
		StandardAPITimeLayoutWithOffset,
		SyncTimeLayoutWithTimezone,
		SyncTimeLayoutWithOffset,
		LegacySyncTimeLayout,
	}

	var err error
	for _, layout := range knownLayouts {
		result, err := time.Parse(layout, datetime)
		if err == nil {
			return result, nil
		}
	}

	return time.Time{}, err
}
