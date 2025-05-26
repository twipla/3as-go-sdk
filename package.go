package twipla3as

import (
	"cmp"
	"context"
	"net/http"
	"path"
	"slices"
	"time"
)

type Currency string

const (
	CurrencyEUR Currency = "EUR"
	CurrencyRON Currency = "RON"
	CurrencyUSD Currency = "USD"
)

type Period string

const (
	PeriodMonthly Period = "monthly"
	PeriodYearly  Period = "yearly"
)

type Package struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Currency    Currency  `json:"currency"`
	Period      Period    `json:"period"`
	Recommended bool      `json:"recommended"`
	IntpID      string    `json:"intpId"`
	Touchpoints float64   `json:"touchpoints"`
}

func (sdk *TwiplaSDK) Packages(ctx context.Context) ([]Package, error) {
	resp, err := parseResponse[[]Package](sdk.apiCall(ctx, http.MethodGet, "/v2/3as/packages", nil))
	if err != nil {
		return nil, err
	}
	slices.SortStableFunc(resp.Payload, func(a, b Package) int {
		return cmp.Compare(a.Touchpoints, b.Touchpoints)
	})
	return resp.Payload, nil
}

func (sdk *TwiplaSDK) Package(ctx context.Context, packageID string) (Package, error) {
	resp, err := parseResponse[Package](sdk.apiCall(ctx, http.MethodGet, path.Join("/v2/3as/packages", packageID), nil))
	if err != nil {
		return Package{}, err
	}
	return resp.Payload, nil
}

type CreatePackageArgs struct {
	Name        string   `json:"name"`
	Touchpoints float64  `json:"touchpoints"`
	Price       float64  `json:"price"`
	Currency    Currency `json:"currency"`
	Period      Period   `json:"period"`
}

func (sdk *TwiplaSDK) CreatePackage(ctx context.Context, args CreatePackageArgs) (Package, error) {
	resp, err := parseResponse[Package](sdk.apiCall(ctx, http.MethodPost, "/v2/3as/packages", args))
	if err != nil {
		return Package{}, err
	}
	return resp.Payload, nil
}

type UpdatePackageArgs struct {
	Name string `json:"name"`
}

func (sdk *TwiplaSDK) UpdatePackage(ctx context.Context, packageID string, args UpdatePackageArgs) (Package, error) {
	resp, err := parseResponse[Package](sdk.apiCall(ctx, http.MethodPatch, path.Join("/v2/3as/packages", packageID), args))
	if err != nil {
		return Package{}, err
	}
	return resp.Payload, nil
}
