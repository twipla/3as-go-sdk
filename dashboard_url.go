package twipla3as

import (
	"fmt"
	"net/url"
)

// GenerateIframeURL generates a URL that can be used to embed the 3as dashboard in an iframe.
// intpcID and websiteID are the INTP's internal IDs for the customer and the website.
func (sdk *TwiplaSDK) GenerateIframeURL(intpcID string, websiteID string) (string, error) {
	var baseURL string
	switch sdk.env {
	case EnvironmentDevelop:
		baseURL = "https://dev-dashboard-3as.va-endpoint.com/"
	case EnvironmentStage:
		baseURL = "https://stage-dashboard-3as.va-endpoint.com/"
	case EnvironmentProduction:
		baseURL = "https://app-3as.visitor-analytics.io/"
	default:
		return "", fmt.Errorf("unsupported env: %s", sdk.env)
	}

	token, err := sdk.signer.IntpcToken(intpcID)
	if err != nil {
		return "", fmt.Errorf("could not generate intpc token: %w", err)
	}
	query := url.Values{}
	query.Set("intpc_token", token)
	query.Set("externalWebsiteId", websiteID)

	return baseURL + "?" + query.Encode(), nil
}
