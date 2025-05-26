package twipla3as

import (
	"context"
	"net/http"
)

type intpcSubscription struct {
	IntpcID   string `json:"intpcId"`
	PackageID string `json:"packageId"`
}

type UpgradeINTPCSubscriptionArgs struct {
	// IntpcID is the ID of the INTPC to upgrade.
	IntpcID string `json:"intpcId"`
	// PackageID is the new package ID to use for the subscription.
	// Note that upgrades need to be strictly better (old touchpoints < new touchpoints)
	// Extra limitations may apply.
	PackageID string `json:"packageId"`

	// Trial marks whether the upgrade should be considered trial or not.
	Trial bool `json:"trial"`

	// Prorate is used in the billing system to determine pricing.
	Prorate bool `json:"proRate"`
}

// UpgradeINTPCSubscription upgrades an INTPC subscription to a new package immediately.
func (sdk *TwiplaSDK) UpgradeINTPCSubscription(ctx context.Context, args UpgradeINTPCSubscriptionArgs) error {
	_, err := parseResponse[intpcSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/intpc-subscriptions/upgrade", args))
	return err
}

type DowngradeINTPCSubscriptionArgs struct {
	// IntpcID is the ID of the INTPC whose subscription to downgrade.
	IntpcID string `json:"intpcId"`
	// PackageID is the new package ID to use for the subscription.
	// Note that downgrades need to be inferior (current touchpoints > future touchpoints)
	// Extra limitations may apply.
	PackageID string `json:"packageId"`
}

// DowngradeINTPCSubscription schedules an INTPC subscription downgrade to a lesser package at the beginning of the next billing period.
func (sdk *TwiplaSDK) DowngradeINTPCSubscription(ctx context.Context, args DowngradeINTPCSubscriptionArgs) error {
	_, err := parseResponse[intpcSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/intpc-subscriptions/downgrade", args))
	return err
}

type ResumeINTPCSubscriptionArgs struct {
	// IntpcID is the ID of the INTPC whose subscription to resume.
	IntpcID string `json:"intpcId"`
}

// ResumeINTPCSubscription resumes an INTPC subscription.
func (sdk *TwiplaSDK) ResumeINTPCSubscription(ctx context.Context, args ResumeINTPCSubscriptionArgs) error {
	_, err := parseResponse[intpcSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/intpc-subscriptions/resume", args))
	return err
}

type DeactivateINTPCSubscriptionArgs struct {
	// IntpcID is the ID of the INTPC whose subscription to deactivate.
	IntpcID string `json:"intpcId"`
}

// DeactivateINTPCSubscription deactivates an INTPC subscription immediately.
func (sdk *TwiplaSDK) DeactivateINTPCSubscription(ctx context.Context, args DeactivateINTPCSubscriptionArgs) error {
	_, err := parseResponse[intpcSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/intpc-subscriptions/deactivate", args))
	return err
}

type CancelINTPCSubscriptionArgs struct {
	// IntpcID is the ID of the INTPC whose subscription to cancel.
	IntpcID string `json:"intpcId"`
}

// CancelINTPCSubscription cancels an INTPC subscription after the end of the billing period.
func (sdk *TwiplaSDK) CancelINTPCSubscription(ctx context.Context, args CancelINTPCSubscriptionArgs) error {
	_, err := parseResponse[intpcSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/intpc-subscriptions/cancel", args))
	return err
}
