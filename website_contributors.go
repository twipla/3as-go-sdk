package twipla3as

import (
	"context"
	"fmt"
	"net/http"
)

type ContributorRole string

const (
	ContributorRoleEditor    ContributorRole = "editor"
	ContributorRoleWatcher   ContributorRole = "watcher"
	ContributorRoleDashboard ContributorRole = "dashboard"
)

type CreateWebsiteContributorArgs struct {
	// Represents the website where contributor access will be granted.
	// ID of the website within the integration partner’s system.
	ExternalWebsiteId string

	// Represents the user who is being added as a contributor.
	// ID of the customer within the integration partner’s system.
	ExternalCustomerId string

	/**
	Defines the role assigned to the contributor.

	Available roles:

	editor — Full edit access, including content updates and structural changes.
	watcher — View-only access to website data. Cannot make any edits.
	dashboard — View-only access to custom dashboards explicitly shared with them. No access to any other platform content or settings. Access to specific dashboards is granted by the website owner from within the dashboard.
	*/
	Role ContributorRole
}

func (sdk *TwiplaSDK) AddWebsiteContributor(ctx context.Context, args CreateWebsiteContributorArgs) error {
	url := fmt.Sprintf("/v3/3as/websites/%v/contributors", args.ExternalWebsiteId)
	body := map[string]any{
		"intpCustomerId": args.ExternalCustomerId,
		"role":           args.Role,
	}
	_, err := parseResponse[any](sdk.apiCall(ctx, http.MethodPost, url, body))
	return err
}

type DeleteWebsiteContributorArgs struct {
	// Represents the website from which the contributor will be removed.
	// ID of the website within the integration partner’s system.
	ExternalWebsiteId string

	// Represents the user who is being removed as a contributor.
	// ID of the customer within the integration partner’s system.
	ExternalCustomerId string
}

func (sdk *TwiplaSDK) DeleteWebsiteContributor(ctx context.Context, args DeleteWebsiteContributorArgs) error {
	url := fmt.Sprintf("/v3/3as/websites/%v/contributors/%v", args.ExternalWebsiteId, args.ExternalCustomerId)
	_, err := parseResponse[any](sdk.apiCall(ctx, http.MethodDelete, url, nil))
	return err
}

type ListWebsiteContributorsArgs struct {
	// ID of the website within the integration partner’s system
	ExternalWebsiteId string
}

type ContributorInfo struct {
	ExternalCustomerId string `json:"intpCustomerId"`
	Email              string `json:"email"`
}

type ListWebsiteContributorsResponse struct {
	Owner        ContributorInfo                       `json:"owner"`
	Contributors map[ContributorRole][]ContributorInfo `json:"contributors"`
}

func (sdk *TwiplaSDK) ListWebsiteContributors(ctx context.Context, args ListWebsiteContributorsArgs) (*ListWebsiteContributorsResponse, error) {
	url := fmt.Sprintf("/v3/3as/websites/%v/contributors", args.ExternalWebsiteId)
	res, err := parseResponse[ListWebsiteContributorsResponse](sdk.apiCall(ctx, http.MethodGet, url, nil))
	if err != nil {
		return nil, err
	}
	return &res.Payload, nil
}
