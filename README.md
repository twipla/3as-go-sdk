# TWIPLA 3AS Go SDK

[![GoDoc](https://pkg.go.dev/badge/github.com/twipla/3as-go-sdk)](https://pkg.go.dev/github.com/twipla/3as-go-sdk)

A simple API wrapper for integrating the Analytics as a Service (3AS) APIs provided by TWIPLA.

## Getting started

1. [Create an RSA Key Pair (PEM format)](#creating-an-rsa-key-pair)
2. Send the resulting public key (`jwtRS256.key.pub`) to the TWIPLA Dev Team
3. [Install the library](#installation)
4. [Use the SDK instance](#how-to-use-the-library) to interact with the API

## Installation
```sh
go get github.com/twipla/3as-go-sdk
```

## How to use the library

Please refer to the example on [pkg.go.dev](https://pkg.go.dev/github.com/twipla/3as-go-sdk).

For example, here is how you can generate a dashboard URL to embed in an iframe (or redirect the customer):

```go
package main

import (
	_ "embed"
	"fmt"
	twipla3as "github.com/twipla/3as-go-sdk"
	"log"
)

// Suppose that the RSA private key you use for signing is stored in jwtRS256.key
//
//go:embed jwtRS256.key
var privateKey string

const (
	// intpID is the integration provider UUID received.
	intpID = "2f8b7fd2-f958-4c10-b9d7-6aa0213ae299"
	// intpcID is your internal customer identifier
	intpcID = "first_test_01"
	// websiteID is your website identifier associated with the customer with intpcID
	websiteID = "first_test_01_website_01"
)

func main() {
	// Create a new SDK instance.
	sdk, err := twipla3as.NewSDK(&twipla3as.TwiplaConfig{
		IntpID:      intpID,
		PrivateKey:  privateKey,
		Environment: twipla3as.EnvironmentProduction,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Generate the dashboard URL for your specific (intpc, website) pair.
	url, err := sdk.GenerateIframeURL(intpcID, websiteID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url)
}
```


## Creating an RSA Key pair

1. Create the keypair: `ssh-keygen -t rsa -b 2048 -m PEM -f jwtRS256.key`
2. Convert the public key to PEM: `openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub`

## Concepts

### Terms

- **INTP (Integration Partner)**\
  The company that is integrating the analytics as a service solution (3AS)
- **STPs (Server Touchpoints)**\
  Credits used to measure data usage for a given website
- **INTPC (INTPC integration partner customer)**\
  One user of the INTP, can have many websites
- **Website**\
  The website where data will be tracked.
  When the website is created a tracking code snippet is returned that must be embedded within the website's HTML.
- **Package**\
  A package has a price and contains a certain number of STPs. They are used when upgrading/downgrading the subscription of a website.
- **Subscription**\
  A subscription has a package with a certain limit of STPs. This subscription can be upgraded or downgraded. There are two types of subscriptions:
    - **Website Subscription**\
      If the INTP is configured with Website Subscriptions, each website has a distinct subscription and their own limits.
    - **INTPC Subscription**\
      If the INTP is configured with INTPC Subscriptions, each customer has a subscription and all of their websites pool their consumed STPs together.

### General

Most endpoints that deal with customers or websites support some form of an ID which can be provided and then used for all following requests.

For example, creating a new customer with a website requires an `intpCustomerId`|`intpcId` and an `intpWebsiteId`. These must be provided by the INTP and are intended to make integrations easier because there is no need to save any external IDs. Then when getting data about a customer the request is done using the same `intpCustomerId` provided on creation.

Paginated endpoints (website and INTPC listing) have pages that are 0-indexed. Alongside the page's results, metadata such as the total number of pages and how many items there are in total are returned as well as the second return value.

### Subscription types

There are currently **two types of subscription** available:

#### 1. `Website` Subscription

- Applies to a **single website**.
- Created using an **`intp` package**, which defines the subscription plan.
- Can be billed **monthly** or **yearly**.
- Each `website` subscription is tied to an **`intpc`**, which is the entity responsible for creating the website.

#### 2. `Intpc` Subscription

- Covers **one or more websites** under a single subscription.
- Created using an **`intp` package**, which defines the subscription plan.
- Can be billed **monthly** or **yearly**.
- The **touchpoint limit** defined by the package is **shared across all associated websites**.
- The `intpc` can **monitor individual usage** per website, providing detailed insights into how each site consumes touchpoints.
- Ideal for managing multiple websites with a **centralized billing**.


### Example implementation flow

1. Create a new intpc with a website
2. Inject the resulting tracking code in the website's HTML
3. Use the SDK's [generate iframe url](#generate-the-visitoranalytics-dashboard-iframe-url) method to create an url
4. Show an iframe to the user with the url created previously
5. Show a modal to the user to upgrade his subscription
6. Display all the available packages using the SDK
7. After the payment is complete, use the SDK to upgrade the subscription of the website


## Available APIs

- [INTPCs](#intpcs-api)
- [INTPC](#intpc-api)
- [Package](#package-api)
- [Packages](#packages-api)
- [Website](#website-api)
- [Websites](#websites-api)
- [Utils](#utils-api)

### INTPCs API

Integration partners (INTP) are able to get data about their customers (INTPc).

#### List all available customers

```go
intpcs, pagination, err := sdk.INTPCs(ctx, twipla3as.Pagination{
    Page:     0, // Pages are 0-indexed
    PageSize: 10,
})
```

#### Get a single customer by its INTP given id

```go
intpc, err := sdk.INTPC(ctx, "INTP_CUSTOMER_ID")
```

#### Register and start an INTPc level subscription. This will allow subsequently added websites to consume from the same `touchpoint` pool provided by the `package` used during setup.

```go
intpc, err := sdk.CreateINTPC(ctx, twipla3as.CreateINTPCArgs{
    ExternalCustomerID: "INTP_CUSTOMER_ID",
    Email:              "INTP_CUSTOMER_EMAIL",
    SubscriptionType:   twipla3as.SubscriptionTypeINTPC,
    PackageID:          "PACKAGE_UUID"
    BillingDate:        time.Now(), // (optional, defaults to current time)
    ExternalWebsiteID:  "INTP_WEBSITE_ID",
    Domain:             "INTP_WEBSITE_DOMAIN_URI",
})
```

#### Register an INTPc and start a website level subscription. Each added website will have its own subscription.

```go
intpc, err := sdk.CreateINTPC(ctx, twipla3as.CreateINTPCArgs{
    ExternalCustomerID: "INTP_CUSTOMER_ID",
    Email:              "INTP_CUSTOMER_EMAIL",
    SubscriptionType:   twipla3as.SubscriptionTypeWebsite,
    PackageID:          "PACKAGE_UUID"
    BillingDate:        time.Now(), // (optional, defaults to current time)
    ExternalWebsiteID:  "INTP_WEBSITE_ID",
    Domain:             "INTP_WEBSITE_DOMAIN_URI",
})
```

### INTPC API

Integration partners (INTP) are able to get data about their customers (INTPc).

#### List all websites belonging to an INTP Customer

```go
websites, pagination, err := sdk.IntpcWebsites(ctx, "INTP_CUSTOMER_ID", twipla3as.Pagination{Page: 0, PageSize: 15})
```

#### Delete a Customer belonging to an INTP

```go
intpc, err := sdk.DeleteINTPC(ctx, "INTP_CUSTOMER_ID")
```

#### Generate the VisitorAnalytics Dashboard IFrame Url

```go
url, err := sdk.GenerateIframeURL("INTP_CUSTOMER_ID", "INTP_WEBSITE_ID")
```

### Packages API

An Integration Partner (INTP) is able to get data about their packages

#### List all available packages

```go
packages, err := sdk.Packages(ctx)
```

#### Get a single package by ID

```go
pkg, err := sdk.Package(ctx, "PACKAGE_UUID")
```

#### Create a package

```go
pkg, err := sdk.CreatePackage(ctx, twipla3as.CreatePackageArgs{
    Name:        "PACKAGE_NAME",
    Touchpoints: TOUCHPOINT_LIMIT,
    Price:       float64(PRICE),
    Period:      twipla3as.PeriodMonthly, // or twipla3as.PeriodYearly
    Currency:    twipla3as.CurrencyEUR, // or twipla3as.CurrencyUSD / twipla3as.CurrencyRON
})
```

### Package API

#### An INTP can update its packages

```go
pkg, err := sdk.UpdatePackage(ctx, "PACKAGE_UUID", twipla3as.UpdatePackageArgs{
    Name: "UPDATED_PACKAGE_NAME",
})
```

### Websites API

#### List all websites

```go
websites, pagination, err := sdk.Websites(ctx, twipla3as.Pagination{
    Page:     0,
    PageSize: 15,
})
```

#### Get a single website by its INTP given id

```go
website, err := sdk.Website(ctx, "INTP_WEBSITE_ID")
```

#### Create a website with its own subscription and attach it to an existing INTPc

```go
err := sdk.CreateWebsite(ctx, twipla3as.CreateWebsiteArgs{
    ExternalID: "INTP_WEBSITE_ID",
    IntpcID:    "INTP_CUSTOMER_ID",
    Domain:     "INTP_WEBSITE_DOMAIN",
    PackageID:  "PACKAGE_UUID",
    BillingDate: time.Now(), // (optional, defaults to current time)
})
```


#### Create a website and attach it to an existing INTPc subscription. This website, alongside other pre-existing website will consume `touchpoints` from the same pool.

```go
err := sdk.CreateWebsite(ctx, twipla3as.CreateWebsiteArgs{
    ExternalID: "INTP_WEBSITE_ID",
    IntpcID:    "INTP_CUSTOMER_ID",
    Domain:     "INTP_WEBSITE_DOMAIN",
})
```

#### Create a website with its own `30 day, unlimited free trial` subscription and attach it to an INTPc. After the 30 day free trial ends, the subscription will be downgraded to the `free` package.

```go
err := sdk.CreateWebsite(ctx, twipla3as.CreateWebsiteArgs{
    ExternalID: "INTP_WEBSITE_ID",
    IntpcID:    "INTP_CUSTOMER_ID",
    Domain:     "INTP_WEBSITE_DOMAIN",
    UFT: true,
})
```


### Website API

#### Delete a website by its INTP given id

```go
err := sdk.DeleteWebsite(ctx, "INTP_WEBSITE_ID")
```

#### Add a whitelisted domain

```go
err := sdk.AddWebsiteWhitelistedDomain(ctx, "INTP_WEBSITE_ID", "google.com")
```

#### Delete a whitelisted domain

```go
err := sdk.RemoveWebsiteWhitelistedDomain(ctx, "INTP_WEBSITE_ID", "google.com")
```

#### List all whitelisted domains

```go
domains, err := sdk.WhitelistedDomains(ctx, "INTP_WEBSITE_ID")
```

### API for managing a subscription of type `website`

#### Upgrade - immediately applies a higher stp count package to the subscription

```go
err := sdk.UpgradeWebsiteSubscription(ctx, twipla3as.UpgradeWebsiteSubscriptionArgs{
    WebsiteID: "INTP_WEBSITE_ID",
    PackageID: "PACKAGE_UUID",
    Trial:     false,
    Prorate:   false,
})
```

#### Downgrade - auto-renew the subscription at the end of the current billing interval to a new lower stp count package

```go
err := sdk.DowngradeWebsiteSubscription(ctx, twipla3as.DowngradeWebsiteSubscriptionArgs{
    WebsiteID: "INTP_WEBSITE_ID",
    PackageID: "PACKAGE_UUID",
})
```

#### Cancel - disable the subscription auto-renewal at the end of the current billing interval

```go
err := sdk.CancelWebsiteSubscription(ctx, twipla3as.CancelWebsiteSubscriptionArgs{
    WebsiteID: "INTP_WEBSITE_ID",
})
```

#### Resume - re-enable the subscription auto-renewal at the end of the current billing interval

```go
err := sdk.ResumeWebsiteSubscription(ctx, twipla3as.ResumeWebsiteSubscriptionArgs{
    WebsiteID: "INTP_WEBSITE_ID",
})
```

#### Deactivate - immediately disables the subscription, reversible by an upgrade

```go
err := sdk.DeactivateWebsiteSubscription(ctx, twipla3as.DeactivateWebsiteSubscriptionArgs{
    WebsiteID: "INTP_WEBSITE_ID",
})
```

### API for managing a subscription of type `intpc`

#### Upgrade - immediately applies a higher stp count package to the subscription

```go
err := sdk.UpgradeINTPCSubscription(ctx, twipla3as.UpgradeINTPCSubscriptionArgs{
    IntpcID:   "INTP_CUSTOMER_ID",
    PackageID: "PACKAGE_UUID",
    Trial:     false,
    Prorate:   false,
})
```

#### Downgrade - auto-renew the subscription at the end of the current billing interval to a new lower stp count package

```go
err := sdk.DowngradeINTPCSubscription(ctx, twipla3as.DowngradeINTPCSubscriptionArgs{
    IntpcID:   "INTP_CUSTOMER_ID",
    PackageID: "PACKAGE_UUID",
})
```

#### Cancel - disable the subscription auto-renewal at the end of the current billing interval

```go
err := sdk.CancelINTPCSubscription(ctx, twipla3as.CancelINTPCSubscriptionArgs{
    IntpcID: "INTP_CUSTOMER_ID",
})
```

#### Resume - re-enable the subscription auto-renewal at the end of the current billing interval

```go
err := sdk.ResumeINTPCSubscription(ctx, twipla3as.ResumeINTPCSubscriptionArgs{
    IntpcID: "INTP_CUSTOMER_ID",
})
```

#### Deactivate - immediately disables the subscription, reversible by an upgrade

```go
err := sdk.DeactivateINTPCSubscription(ctx, twipla3as.DeactivateINTPCSubscriptionArgs{
    IntpcID: "INTP_CUSTOMER_ID",
})
```

### Utils API

#### Generate a valid access token for the current INTP configuration.

```go
token, err := sdk.IntpAccessToken()
```

#### Generate a valid access token for the current INTPc configuration.

```go
token, err := sdk.IntpcAccessToken("INTP_CUSTOMER_ID")
```

## Dashboard IFrame

The IFrame is one of the main ways a user can interract with the data gathered for his website. The URL of the IFrame is [generated using the SDK](#generate-the-visitoranalytics-dashboard-iframe-url)

The resulting URL can be further enhanced with query parameters:

1. `allowUpgrade=true` - Show upgrade CTAs

Upgrade buttons will be added to the Dashboard for all features that require a certain minimum package.
Once the upgrade button is clicked, the iframe posts a message to the parent frame, containing the following payload:

```json5
{
    "type": "UPGRADE_BUTTON_CLICKED",
    "data": {
        "intpWebsiteId": "", // string; external website id
        "intpCustomerId": "", // string; customer id
        "packageName": "", // string; current package name
        "packageId": "", // string; current package id
        "inTrial": false, // boolean;
        "expiresAt": "", // string; expiry date in ISO 8601 format
        "billingInterval": "monthly" // "monthly"|"yearly";
    }
}
```

