package twipla3as_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	twipla3as "github.com/twipla/3as-go-sdk"
	"math/rand/v2"
	"testing"
	"time"
)

func TestWhitelistedDomains(t *testing.T) {

	var subType twipla3as.SubscriptionType
	if mainSDK == websiteSubSDK {
		subType = twipla3as.SubscriptionTypeWebsite
	} else if mainSDK == intpcSubSDK {
		subType = twipla3as.SubscriptionTypeINTPC
	} else {
		t.Skip("Unknown main SDK type")
	}

	packages, err := mainSDK.Packages(t.Context())
	assert.NoError(t, err)
	var pkg twipla3as.Package
	for _, p := range packages {
		if p.Touchpoints > 0 {
			pkg = p
			break
		}
	}

	intpcName := fmt.Sprintf("go-sdk-intpc-%d", rand.Int())
	websiteName := fmt.Sprintf("go-sdk-website-%d", rand.Int())
	rndEmail := fmt.Sprintf("%d@twipla.com", rand.Int())
	rndDomain := fmt.Sprintf("%d.twiplatest.com", rand.Int())
	intpc, err := mainSDK.CreateINTPC(t.Context(), twipla3as.CreateINTPCArgs{
		ExternalCustomerID: intpcName,
		Email:              rndEmail,
		SubscriptionType:   subType,
		PackageID:          pkg.ID,
		BillingDate:        time.Now(),
		ExternalWebsiteID:  websiteName,
		Domain:             rndDomain,
	})
	assert.NoError(t, err)
	assert.NotNil(t, intpc)
	assert.NotEmpty(t, intpc.ID)
	assert.Equal(t, intpcName, intpc.IntpCustomerID)
	assert.Equal(t, rndEmail, intpc.Email)

	defer func() {
		_, err := mainSDK.DeleteINTPC(t.Context(), intpcName)
		assert.NoError(t, err)
	}()

	t.Run("Add", func(t *testing.T) {
		assert.NoError(t, mainSDK.AddWebsiteWhitelistedDomain(t.Context(), websiteName, "google.com"))
	})
	t.Run("List", func(t *testing.T) {
		domains, err := mainSDK.WhitelistedDomains(t.Context(), websiteName)
		assert.NoError(t, err)
		assert.NotEmpty(t, domains)
		assert.Contains(t, domains, "google.com")
	})
	t.Run("Delete", func(t *testing.T) {
		assert.NoError(t, mainSDK.RemoveWebsiteWhitelistedDomain(t.Context(), websiteName, "google.com"))
		assert.NoError(t, mainSDK.RemoveWebsiteWhitelistedDomain(t.Context(), websiteName, "notactuallywhitelisted.com"))
	})
}
