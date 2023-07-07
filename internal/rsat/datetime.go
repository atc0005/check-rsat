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

const (
	// StandardAPITimeLayout is the time layout format as used by the Red Hat
	// Satellite API for the majority of the date/time properties.
	//
	// Examples:
	//
	// - "created_at": "2020-12-03 15:05:00 UTC",
	// - "updated_at": "2020-12-03 15:05:00 UTC",
	StandardAPITimeLayout string = "2006-01-02 15:04:05 MST"

	// SyncTimeLayout is the time layout format as used by the Red Hat
	// Satellite Sync Plans API for the next_sync property.
	//
	// Example: "next_sync": "2022/03/28 20:05:00 +0000"
	//
	// See also https://rsat.example.com/apidoc/v2/sync_plans/index.html
	SyncTimeLayout string = "2006/01/02 15:04:05 -0700"
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
	return dt.Format(StandardAPITimeLayout)
}

// String implements the fmt.Stringer interface as a convenience method.
func (dt SyncTime) String() string {
	// return dt.Format(StandardAPITimeLayout)
	switch {
	case time.Time(dt).IsZero():
		return "Not scheduled"
	default:
		return time.Time(dt).Local().Format(StandardAPITimeLayout)
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
// time.Time format to the JSON API expectations of a time value in the
// StandardAPITimeLayout format.
func (dt StandardAPITime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dt).Format(StandardAPITimeLayout))
}

// MarshalJSON implements the json.Marshaler interface. This compliments the
// custom Unmarshaler implementation to handle conversion of a native Go
// time.Time format to the JSON API expectations of a time value in the
// SyncTimeLayout format.
func (dt SyncTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dt).Format(SyncTimeLayout))
}

// UnmarshalJSON implements the json.Unmarshaler interface to handle
// converting a time string from the JSON API (most time properties) to a
// native Go time.Time value using the StandardAPITimeLayout format.
func (dt *StandardAPITime) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), `"`) // get rid of "
	if value == "" || value == JSONNullKeyword {

		// Per json.Unmarshaler convention we treat "null" value as a no-op.
		return nil
	}

	// Parse time, explicitly setting UTC location (even though the JSON API
	// already indicates this). We do this for consistency with the next_sync
	// property.
	loc, _ := time.LoadLocation("UTC")
	t, err := time.ParseInLocation(StandardAPITimeLayout, value, loc)
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

	// Parse time, forcing the location to UTC. We do this to match the same
	// timezone used in other time properties.
	//
	// Example of the next_sync JSON property & value:
	//
	// "next_sync": "2022/03/28 20:05:00 +0000"
	loc, _ := time.LoadLocation("UTC")
	t, err := time.ParseInLocation(SyncTimeLayout, value, loc)
	if err != nil {
		return err
	}

	*dt = SyncTime(t) // set result using the pointer

	return nil
}
