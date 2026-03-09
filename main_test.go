// Package twipla3as_test holds the internal integration tests for the 3AS SDK.
// To run them, you need a private key for an INTP with INTPC subscriptions and another one with INTP subscriptions.
package twipla3as_test

import (
	"cmp"
	_ "embed"
	"fmt"
	"log"
	"math/rand/v2"
	"testing"
	"time"

	twipla3as "github.com/twipla/3as-go-sdk"
)

var (
	//go:embed jwtRS256_intpc.key
	privateKeyIntpc string
	intpcSubSDK     *twipla3as.TwiplaSDK

	//go:embed jwtRS256_website.key
	privateKeyWebsite string
	websiteSubSDK     *twipla3as.TwiplaSDK

	mainSDK *twipla3as.TwiplaSDK
)

func TestMain(m *testing.M) {
	var err error
	intpcSubSDK, err = twipla3as.NewSDK(&twipla3as.TwiplaConfig{
		IntpID:      "d71ff306-f21e-4124-bf54-9b75148516e5",
		PrivateKey:  privateKeyIntpc,
		Environment: twipla3as.EnvironmentDevelop,
	})
	if err != nil {
		log.Printf("will skip intpc-sub intp tests: %v", err)
		intpcSubSDK = nil
	}

	websiteSubSDK, err = twipla3as.NewSDK(&twipla3as.TwiplaConfig{
		IntpID:      "14257a0c-f436-4ef2-af91-2ee7ca30ff72",
		PrivateKey:  privateKeyWebsite,
		Environment: twipla3as.EnvironmentDevelop,
	})
	if err != nil {
		log.Printf("will skip website-sub intp tests: %v", err)
		websiteSubSDK = nil
	}

	if websiteSubSDK == nil && intpcSubSDK == nil {
		log.Fatal("missing private keys for testing")
	}
	mainSDK = cmp.Or(websiteSubSDK, intpcSubSDK)

	m.Run()
}

func randIntpcId() string {
	return fmt.Sprintf("go-sdk-intpc-%d", rand.Int())
}

func randEmail() string {
	return fmt.Sprintf("%d@twipla.com", rand.Int())
}

func randWebsiteId() string {
	return fmt.Sprintf("go-sdk-website-%d", rand.Int())
}

func randDomain() string {
	return fmt.Sprintf("%d.twiplatest.com", rand.Int())
}

func randINTPC(subscriptionType twipla3as.SubscriptionType) twipla3as.CreateINTPCArgs {
	return twipla3as.CreateINTPCArgs{
		ExternalCustomerID: fmt.Sprintf("go-sdk-intpc-%d", rand.Int()),
		Email:              fmt.Sprintf("%d@twipla.com", rand.Int()),
		SubscriptionType:   subscriptionType,
		BillingDate:        time.Now(),
		ExternalWebsiteID:  fmt.Sprintf("go-sdk-website-%d", rand.Int()),
		Domain:             fmt.Sprintf("%d.twiplatest.com", rand.Int()),
	}
}
