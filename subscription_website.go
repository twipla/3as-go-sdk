package twipla3as

import (
	"context"
	"net/http"
)

type websiteSubscription struct {
	WebsiteID string `json:"intpWebsiteId"`
	PackageID string `json:"packageId"`
}

type UpgradeWebsiteSubscriptionArgs struct {
	// WebsiteID is the ID of the website whose subscription to upgrade.
	WebsiteID string `json:"intpWebsiteId"`
	// PackageID is the new package ID to use for the subscription.
	// Note that upgrades need to be strictly better (old touchpoints < new touchpoints)
	// Extra limitations may apply.
	PackageID string `json:"packageId"`

	// Trial marks whether the upgrade should be considered trial or not.
	Trial bool `json:"trial"`

	// Prorate is used in the billing system to determine pricing.
	Prorate bool `json:"proRate"`
}

// UpgradeWebsiteSubscription upgrades a Website subscription to a new package immediately.
func (sdk *TwiplaSDK) UpgradeWebsiteSubscription(ctx context.Context, args UpgradeWebsiteSubscriptionArgs) error {
	_, err := parseResponse[websiteSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/website-subscriptions/upgrade", args))
	return err
}

type DowngradeWebsiteSubscriptionArgs struct {
	// WebsiteID is the ID of the website whose subscription to downgrade.
	WebsiteID string `json:"intpWebsiteId"`
	// PackageID is the new package ID to use for the subscription.
	// Note that downgrades need to be inferior (current touchpoints > future touchpoints)
	// Extra limitations may apply.
	PackageID string `json:"packageId"`
}

// DowngradeWebsiteSubscription schedules a Website subscription downgrade to a lesser package at the beginning of the next billing period.
func (sdk *TwiplaSDK) DowngradeWebsiteSubscription(ctx context.Context, args DowngradeWebsiteSubscriptionArgs) error {
	_, err := parseResponse[websiteSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/website-subscriptions/downgrade", args))
	return err
}

type ResumeWebsiteSubscriptionArgs struct {
	// WebsiteID is the ID of the website whose subscription to resume.
	WebsiteID string `json:"intpWebsiteId"`
}

// ResumeWebsiteSubscription resumes a Website subscription.
func (sdk *TwiplaSDK) ResumeWebsiteSubscription(ctx context.Context, args ResumeWebsiteSubscriptionArgs) error {
	_, err := parseResponse[websiteSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/website-subscriptions/resume", args))
	return err
}

type DeactivateWebsiteSubscriptionArgs struct {
	// WebsiteID is the ID of the website whose subscription to deactivate.
	WebsiteID string `json:"intpWebsiteId"`
}

// DeactivateWebsiteSubscription deactivates a Website subscription immediately.
func (sdk *TwiplaSDK) DeactivateWebsiteSubscription(ctx context.Context, args DeactivateWebsiteSubscriptionArgs) error {
	_, err := parseResponse[websiteSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/website-subscriptions/deactivate", args))
	return err
}

type CancelWebsiteSubscriptionArgs struct {
	// WebsiteID is the ID of the website whose subscription to cancel.
	WebsiteID string `json:"intpWebsiteId"`
}

// CancelWebsiteSubscription cancels a Website subscription after the end of the billing period.
func (sdk *TwiplaSDK) CancelWebsiteSubscription(ctx context.Context, args CancelWebsiteSubscriptionArgs) error {
	_, err := parseResponse[websiteSubscription](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/website-subscriptions/cancel", args))
	return err
}
