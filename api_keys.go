package twipla3as

import (
	"context"
	"net/http"
	"path"
	"time"
)

type ApiKey struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	ApiKey         *string   `json:"apiKey"`
	Comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"createdAt"`
	ExpiresAt      time.Time `json:"expiresAt"`
	IntpWebsiteId  string    `json:"intpWebsiteId"`
	IntpCustomerId string    `json:"intpCustomerId"`
}

type CreateApiKeyArgs struct {
	// ExternalID is the websiteID to use
	ExternalWebsiteID string
	// Name to identify the API key
	Name string
	// Optional description or notes
	Comment *string
	// // Expiration timestamp; unlimited if omitted
	ExpiresAt *time.Time
}

func (sdk *TwiplaSDK) CreateWebsiteApiKey(ctx context.Context, args CreateApiKeyArgs) (*ApiKey, error) {
	reqBody := map[string]interface{}{
		"name": args.Name,
	}
	if args.Comment != nil {
		reqBody["comment"] = args.Comment
	}

	if args.ExpiresAt != nil {
		reqBody["expiresAt"] = args.ExpiresAt.Format(time.RFC3339)
	}

	res, err := parseResponse[ApiKey](sdk.apiCall(ctx, http.MethodPost, path.Join("/v2/3as/websites", args.ExternalWebsiteID, "api-keys"), reqBody))
	if err != nil {
		return nil, err
	}
	return &res.Payload, err
}

func (sdk *TwiplaSDK) ListWebsiteApiKeys(ctx context.Context, externalWebsiteId string) ([]ApiKey, error) {
	res, err := parseResponse[[]ApiKey](sdk.apiCall(ctx, http.MethodGet, path.Join("/v2/3as/websites", externalWebsiteId, "api-keys"), nil))
	if err != nil {
		return nil, err
	}
	return res.Payload, err
}

func (sdk *TwiplaSDK) DeleteWebsiteApiKey(ctx context.Context, externalWebsiteId string, apiKeyId string) error {
	_, err := parseResponse[any](sdk.apiCall(ctx, http.MethodDelete, path.Join("/v2/3as/websites", externalWebsiteId, "api-keys", apiKeyId), nil))

	return err
}
