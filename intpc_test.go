package twipla3as_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	twipla3as "github.com/twipla/3as-go-sdk"
	"math/rand/v2"
	"testing"
	"time"
)

func TestINTPCs(t *testing.T) {
	t.Run("INTPC subscriptions", func(t *testing.T) {
		if intpcSubSDK == nil {
			t.Skip("No INTPC Subscription SDK set")
		}
		packages, err := intpcSubSDK.Packages(t.Context())
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
		t.Run("Create", func(t *testing.T) {
			intpc, err := intpcSubSDK.CreateINTPC(t.Context(), twipla3as.CreateINTPCArgs{
				ExternalCustomerID: intpcName,
				Email:              rndEmail,
				SubscriptionType:   twipla3as.SubscriptionTypeINTPC,
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
		})

		t.Run("List", func(t *testing.T) {
			intpcs, pagination, err := intpcSubSDK.INTPCs(t.Context(), twipla3as.Pagination{
				Page:     0,
				PageSize: 15,
			})
			assert.NoError(t, err)
			assert.NotEmpty(t, intpcs)
			assert.NotEmpty(t, pagination)
			assert.Equal(t, 15, pagination.PageSize)
		})

		t.Run("Delete", func(t *testing.T) {
			intpc, err := intpcSubSDK.DeleteINTPC(t.Context(), intpcName)
			assert.NoError(t, err)
			assert.NotEmpty(t, intpc)
			assert.Equal(t, intpcName, intpc.IntpCustomerID)
		})
	})

	t.Run("Website subscriptions", func(t *testing.T) {
		if websiteSubSDK == nil {
			t.Skip("No Website Subscription SDK set")
		}
		packages, err := websiteSubSDK.Packages(t.Context())
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
		t.Run("Create", func(t *testing.T) {
			intpc, err := websiteSubSDK.CreateINTPC(t.Context(), twipla3as.CreateINTPCArgs{
				ExternalCustomerID: intpcName,
				Email:              rndEmail,
				SubscriptionType:   twipla3as.SubscriptionTypeWebsite,
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
		})

		t.Run("List", func(t *testing.T) {
			intpcs, pagination, err := websiteSubSDK.INTPCs(t.Context(), twipla3as.Pagination{
				Page:     0,
				PageSize: 15,
			})
			assert.NoError(t, err)
			assert.NotEmpty(t, intpcs)
			assert.NotEmpty(t, pagination)
			assert.Equal(t, 15, pagination.PageSize)
		})

		t.Run("Delete", func(t *testing.T) {
			intpc, err := websiteSubSDK.DeleteINTPC(t.Context(), intpcName)
			assert.NoError(t, err)
			assert.NotEmpty(t, intpc)
			assert.Equal(t, intpcName, intpc.IntpCustomerID)
		})
	})
}
