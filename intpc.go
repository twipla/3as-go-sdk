package twipla3as

import (
	"context"
	"net/http"
	"path"
	"time"
)

type SubscriptionType string

const (
	SubscriptionTypeNone    SubscriptionType = ""
	SubscriptionTypeWebsite SubscriptionType = "website"
	SubscriptionTypeINTPC   SubscriptionType = "intpc"
)

type INTPC struct {
	// ID is the internal 3AS UUID for the customer.
	ID string `json:"id"`
	// IntpCustomerID is the INTP's id value for the customer.
	IntpCustomerID string `json:"intpCustomerID"`
	// VisaID is the internal TWIPLA UUID for the customer.
	// It should be the same as ID, but it is not guaranteed.
	VisaID string `json:"visaId"`
	// Email is the given email for the customer.
	Email string `json:"email"`
	// IntpID is the id associated with the INTP.
	IntpID string `json:"intpId"`
	// CreatedAt is the customer's creation time with the TWIPLA PI.
	CreatedAt time.Time `json:"createdAt"`
}

type CreateINTPCArgs struct {
	// ExternalCustomerID is the intpcID to use when interacting with the SDK.
	// It should be the INTP's internal representation of the customer ID.
	ExternalCustomerID string
	// Email is the customer's associated email address.
	Email string

	// SubscriptionType defines which type of subscription to create.
	// The field must match the INTP subscription type configured for the integration.
	//  - If SubscriptionType is SubscriptionTypeWebsite, a website subscription will be created for the new website.
	//  - If SubscriptionType is SubscriptionTypeINTPC, an INTPC subscription and an associated pool of credits will be created. Further websites created must NOT specify subscription data.
	// SubscriptionType should hold one of these two values.
	SubscriptionType SubscriptionType
	// PackageID holds the package ID for the newly created subscription, regardless of type.
	PackageID string
	// BillingDate is an optional field that sets the createdAt of the subscription.
	// If empty (zero value), it defaults to the current time.
	BillingDate time.Time

	// ExternalWebsiteID is the websiteID to use when interacting with the SDK.
	// It should be the INTP's internal representation of the website ID.
	ExternalWebsiteID string
	// Domain is the host part of the website URL (example: `mail.google.com`, `twipla.com`)
	Domain string
}

func (sdk *TwiplaSDK) CreateINTPC(ctx context.Context, args CreateINTPCArgs) (INTPC, error) {
	if args.BillingDate.IsZero() {
		args.BillingDate = time.Now()
	}
	var apiArgs createIntpcAPIArgs
	apiArgs.IntpCustomerID = args.ExternalCustomerID
	apiArgs.Email = args.Email
	switch args.SubscriptionType {
	case SubscriptionTypeWebsite:
		apiArgs.Website.PackageID = args.PackageID
		apiArgs.Website.BillingDate = args.BillingDate.UTC().Format(time.RFC3339)
	case SubscriptionTypeINTPC:
		apiArgs.PackageID = args.PackageID
		apiArgs.BillingDate = args.BillingDate.UTC().Format(time.RFC3339)
	default:
		return INTPC{}, ErrInvalidSubscriptionType
	}
	apiArgs.Website.IntpWebsiteID = args.ExternalWebsiteID
	apiArgs.Website.Domain = args.Domain
	resp, err := parseResponse[INTPC](sdk.apiCall(ctx, http.MethodPost, "/v2/3as/customers", apiArgs))
	if err != nil {
		return INTPC{}, err
	}
	return resp.Payload, nil
}

func (sdk *TwiplaSDK) INTPCs(ctx context.Context, pagination Pagination) ([]INTPC, PaginationMetadata, error) {
	query := pagination.buildQuery()
	resp, err := parseResponse[[]INTPC](sdk.apiCall(ctx, http.MethodGet, "/v2/3as/customers", query))
	if err != nil {
		return nil, PaginationMetadata{}, err
	}

	return resp.Payload, resp.Metadata, nil
}

// INTPC gets an INTPC/customer based on the INTP's own customer ID.
func (sdk *TwiplaSDK) INTPC(ctx context.Context, intpcID string) (INTPC, error) {
	resp, err := parseResponse[INTPC](sdk.apiCall(ctx, http.MethodGet, path.Join("/v2/3as/customers", intpcID), nil))
	if err != nil {
		return INTPC{}, err
	}
	return resp.Payload, nil
}

// DeleteINTPC removes an INTPC and its linked websites.
func (sdk *TwiplaSDK) DeleteINTPC(ctx context.Context, intpcID string) (INTPC, error) {
	resp, err := parseResponse[INTPC](sdk.apiCall(ctx, http.MethodDelete, path.Join("/v2/3as/customers", intpcID), nil))
	if err != nil {
		return INTPC{}, err
	}
	return resp.Payload, nil
}

type createIntpcAPIArgs struct {
	IntpCustomerID string `json:"intpCustomerId"`
	Email          string `json:"email"`
	PackageID      string `json:"packageId,omitempty"`
	BillingDate    string `json:"billingDate,omitempty"`
	Website        struct {
		IntpWebsiteID string `json:"intpWebsiteId"`
		Domain        string `json:"domain"`
		PackageID     string `json:"packageId,omitempty"`
		BillingDate   string `json:"billingDate,omitempty"`
	} `json:"website"`
}
