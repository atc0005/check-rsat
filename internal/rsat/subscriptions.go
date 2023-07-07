// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/check-rsat
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package rsat

// SubscriptionsResponse represents the API response from a request of all
// subscriptions for a specific organization.
type SubscriptionsResponse struct {
	Error         NullString     `json:"error"`
	Organization  struct{}       `json:"organization"` // I have only encountered: "organization": {},
	Page          int            `json:"page"`
	PerPage       int            `json:"per_page"`
	Subscriptions []Subscription `json:"results"`
	Search        NullString     `json:"search"`
	Sort          SortOptions    `json:"sort"`
	Subtotal      int            `json:"subtotal"`
	Total         int            `json:"total"`
}

// Subscription represents an entitlement for receiving content and service
// from Red Hat. Subscription allocations are applied/managed separately
// within each Red Hat Satellite organization.
type Subscription struct {
	Hypervisor         Hypervisor      `json:"hypervisor,omitempty"`
	StartDate          StandardAPITime `json:"start_date"`
	EndDate            StandardAPITime `json:"end_date"`
	Cores              interface{}     `json:"cores"`             // null is the only value I've encountered
	MultiEntitlement   *bool           `json:"multi_entitlement"` // null or true/false
	AccountNumber      *int            `json:"account_number"`    // null or integer
	Available          int             `json:"available"`
	Consumed           int             `json:"consumed"`
	Quantity           int             `json:"quantity"`
	SubscriptionID     int             `json:"subscription_id"`
	ID                 int             `json:"id"`
	InstanceMultiplier int             `json:"instance_multiplier"`
	RAM                NullString      `json:"ram"`
	Sockets            NullString      `json:"sockets"`
	StackingID         NullString      `json:"stacking_id"`
	SupportLevel       NullString      `json:"support_level"`
	UpstreamPoolID     NullString      `json:"upstream_pool_id"`
	CpID               string          `json:"cp_id"`
	Name               string          `json:"name"`
	ProductID          string          `json:"product_id"`
	ProductName        string          `json:"product_name"`
	Type               string          `json:"type"`
	UnmappedGuest      bool            `json:"unmapped_guest"`
	Upstream           bool            `json:"upstream"`
	VirtOnly           bool            `json:"virt_only"`
	VirtWho            bool            `json:"virt_who"`
}

// Hypervisor represents the hypervisor associated with a specific
// subscription. Not all subscriptions are associated with a hypervisor;
// subscriptions  associated with a hypervisor require that a virtual guest be
// running on the indicated hypervisor (tracked via virt-who) to be eligible
// for that subscription (e.g., Red Hat Enterprise Linux Extended Life Cycle
// Support).
type Hypervisor struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
