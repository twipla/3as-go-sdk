package twipla3as

import (
	"context"
	"net/http"
	"path"
)

func (sdk *TwiplaSDK) AddWebsiteWhitelistedDomain(ctx context.Context, websiteID string, domain string) error {
	_, err := parseResponse[any](sdk.apiCall(ctx, http.MethodPost, path.Join("/v2/3as/websites", websiteID, "whitelisted-domains"), map[string]string{"domain": domain}))
	return err
}

func (sdk *TwiplaSDK) RemoveWebsiteWhitelistedDomain(ctx context.Context, websiteID string, domain string) error {
	_, err := parseResponse[any](sdk.apiCall(ctx, http.MethodPatch, path.Join("/v2/3as/websites", websiteID, "whitelisted-domains"), map[string]string{"domain": domain}))
	return err
}

func (sdk *TwiplaSDK) WhitelistedDomains(ctx context.Context, websiteID string) ([]string, error) {
	resp, err := parseResponse[[]string](sdk.apiCall(ctx, http.MethodGet, path.Join("/v2/3as/websites", websiteID, "whitelisted-domains"), nil))
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
