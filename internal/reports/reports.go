// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package reports

import (
	"fmt"
	"io"

	"github.com/atc0005/go-nagios"
)

func addSyncPlansReportLeadIn(w io.Writer) {
	fmt.Fprintf(
		w,
		"%sSYNC PLANS OVERVIEW%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

}
