package twipla3as

import (
	"cmp"
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type SubscriptionState string

const (
	SubscriptionStateActive    SubscriptionState = "active"
	SubscriptionStateCancelled SubscriptionState = "canceled" // sic
	SubscriptionStateInactive  SubscriptionState = "inactive"
)

type Website struct {
	// The ID field is TWIPLA's UUID of the website.
	ID string `json:"id"`
	// ExternalWebsiteID is the ID of the website as supplied by the INTP.
	ExternalWebsiteID string `json:"intpWebsiteId"`

	// VisaCustomerID is TWIPLA's UUID of the customer (INTPC).
	VisaCustomerID string `json:"visaCustomerId"`
	// IntpCustomerID is the ID of the customer as supplied by the INTP.
	IntpCustomerID string `json:"intpCustomerId"`

	// IntpID is the given ID of the INTP.
	IntpID string `json:"intpId"`
	// Status represents the current subscription state of the website.
	Status SubscriptionState `json:"status"`

	// Domain holds the set domain of the website.
	Domain string `json:"domain"`
	// PackageID is the UUID of the website's package subscription.
	PackageID string `json:"packageId"`
	// PackageID is the given name of the website's package subscription.
	PackageName string `json:"packageName"`
	// BillingInterval is the active package's period.
	BillingInterval Period `json:"billingInterval"`
	// BillingMode is the INTP's configured billing mode. Usually "company-managed".
	BillingMode string `json:"billingMode"`

	LastPackageChangeAt time.Time `json:"lastPackageChangeAt"`
	InTrial             bool      `json:"inTrial"`
	HadTrial            bool      `json:"hadTrial"`
	CreatedAt           time.Time `json:"createdAt"`
	ExpiresAt           time.Time `json:"expiresAt"`
	// StpResetAt is the timestamp at which the website's credit quota is reset.
	StpResetAt time.Time `json:"stpResetAt"`
}

type CreateWebsiteArgs struct {
	// ExternalID is the websiteID to use
	ExternalID string
	// IntpcID is the INTP's internal ID for the customer to link the website to.
	IntpcID string
	// Domain is the host part of the website URL (example: `mail.google.com`, `twipla.com`)
	Domain string

	// PackageID holds the package ID for the website if the INTP was configured for website subscriptions.
	// If the INTP uses INTPC subscriptions, this field must be omitted.
	PackageID string
	// BillingDate is an optional field that sets the createdAt of the website subscription.
	// If empty (zero value), it defaults to the current time.
	// It does not do anything if PackageID is empty
	BillingDate time.Time

	// UFT indicates whether the Unlimited Free Trial option is set on website creation.
	// Recommended to remain false unless you know what you're doing.
	UFT bool
}

func (sdk *TwiplaSDK) CreateWebsite(ctx context.Context, args CreateWebsiteArgs) error {
	if args.BillingDate.IsZero() {
		args.BillingDate = time.Now()
	}
	var apiArgs createWebsiteAPIArgs
	apiArgs.Website.ID = args.ExternalID
	apiArgs.Website.Domain = args.Domain
	apiArgs.Website.Package.ID = args.PackageID
	apiArgs.Website.Package.BillingDate = args.BillingDate.UTC().Format(time.RFC3339)
	apiArgs.Intpc.ID = args.IntpcID
	apiArgs.Opts.UFT = args.UFT
	_, err := parseResponse[any](sdk.apiCall(ctx, http.MethodPost, "/v3/3as/websites", apiArgs))
	return err
}

func (sdk *TwiplaSDK) Websites(ctx context.Context, pagination Pagination) ([]Website, PaginationMetadata, error) {
	return sdk.websites(ctx, "", pagination)
}

func (sdk *TwiplaSDK) IntpcWebsites(ctx context.Context, intpcID string, pagination Pagination) ([]Website, PaginationMetadata, error) {
	return sdk.websites(ctx, intpcID, pagination)
}

// Website gets a website based on the INTP's own website ID.
func (sdk *TwiplaSDK) Website(ctx context.Context, websiteID string) (Website, error) {
	resp, err := parseResponse[Website](sdk.apiCall(ctx, http.MethodGet, path.Join("/v2/3as/websites", websiteID), nil))
	if err != nil {
		return Website{}, err
	}
	return resp.Payload, nil
}

// DeleteWebsite gets a website based on the INTP's own website ID.
func (sdk *TwiplaSDK) DeleteWebsite(ctx context.Context, websiteID string) error {
	_, err := parseResponse[any](sdk.apiCall(ctx, http.MethodDelete, path.Join("/v2/3as/websites", websiteID), nil))
	return err
}

func (sdk *TwiplaSDK) websites(ctx context.Context, intpcID string, pagination Pagination) ([]Website, PaginationMetadata, error) {
	query := pagination.buildQuery()
	if len(intpcID) > 0 {
		query.Set("externalCustomerId", intpcID)
	}
	resp, err := parseResponse[[]Website](sdk.apiCall(ctx, http.MethodGet, "/v2/3as/websites", query))
	if err != nil {
		return nil, PaginationMetadata{}, err
	}

	return resp.Payload, resp.Metadata, nil
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (p Pagination) buildQuery() url.Values {
	vals := url.Values{}
	vals.Set("page", strconv.Itoa(p.Page))
	vals.Set("pageSize", strconv.Itoa(cmp.Or(p.PageSize, 10)))
	return vals
}

type createWebsiteAPIArgs struct {
	Website struct {
		ID      string `json:"id"`
		Domain  string `json:"domain"`
		Package struct {
			ID          string `json:"id"`
			BillingDate string `json:"billingDate,omitempty,omitzero"`
		} `json:"package"`
	} `json:"website"`
	Intpc struct {
		ID string `json:"id"`
	} `json:"intpc"`
	Opts struct {
		UFT bool `json:"uft"`
	} `json:"opts"`
}
