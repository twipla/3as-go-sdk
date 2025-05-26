package twipla3as_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	twipla3as "github.com/twipla/3as-go-sdk"
	"math/rand/v2"
	"testing"
	"time"
)

func TestIntpcSubscriptions(t *testing.T) {
	if intpcSubSDK == nil {
		t.Skip("No INTPC Subscription SDK set")
	}

	packages, err := intpcSubSDK.Packages(t.Context())
	assert.NoError(t, err)
	var pkg twipla3as.Package
	for _, p := range packages {
		if p.Touchpoints > 1000 {
			pkg = p
			break
		}
	}
	var superiorPkg twipla3as.Package
	for _, p := range packages {
		if p.Touchpoints > pkg.Touchpoints && p.Period == pkg.Period {
			superiorPkg = p
			break
		}
	}

	intpcName := fmt.Sprintf("go-sdk-intpc-%d", rand.Int())
	websiteName := fmt.Sprintf("go-sdk-website-%d", rand.Int())
	rndEmail := fmt.Sprintf("%d@twipla.com", rand.Int())
	rndDomain := fmt.Sprintf("%d.twiplatest.com", rand.Int())
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

	defer intpcSubSDK.DeleteINTPC(t.Context(), intpcName)

	t.Run("Upgrade", func(t *testing.T) {
		assert.NoError(t, intpcSubSDK.UpgradeINTPCSubscription(t.Context(), twipla3as.UpgradeINTPCSubscriptionArgs{
			IntpcID:   intpcName,
			PackageID: superiorPkg.ID,
			Trial:     false,
			Prorate:   false,
		}))
	})

	t.Run("Downgrade", func(t *testing.T) {
		// It was upgraded to a superior package before, so we can downgrade to the original package
		assert.NoError(t, intpcSubSDK.DowngradeINTPCSubscription(t.Context(), twipla3as.DowngradeINTPCSubscriptionArgs{
			IntpcID:   intpcName,
			PackageID: pkg.ID,
		}))
	})

	t.Run("Cancel", func(t *testing.T) {
		assert.NoError(t, intpcSubSDK.CancelINTPCSubscription(t.Context(), twipla3as.CancelINTPCSubscriptionArgs{
			IntpcID: intpcName,
		}))
	})

	t.Run("Resume", func(t *testing.T) {
		assert.NoError(t, intpcSubSDK.ResumeINTPCSubscription(t.Context(), twipla3as.ResumeINTPCSubscriptionArgs{
			IntpcID: intpcName,
		}))
	})

	t.Run("Deactivate", func(t *testing.T) {
		assert.NoError(t, intpcSubSDK.DeactivateINTPCSubscription(t.Context(), twipla3as.DeactivateINTPCSubscriptionArgs{
			IntpcID: intpcName,
		}))
	})

	t.Run("Reactivate through upgrade", func(t *testing.T) {
		t.Skip("Broken in aaas-api. Issue raised")
		assert.NoError(t, intpcSubSDK.UpgradeINTPCSubscription(t.Context(), twipla3as.UpgradeINTPCSubscriptionArgs{
			IntpcID:   intpcName,
			PackageID: pkg.ID,
			Trial:     false,
			Prorate:   false,
		}))
	})
}

func TestWebsiteSubscriptions(t *testing.T) {
	if websiteSubSDK == nil {
		t.Skip("No Website Subscription SDK set")
	}

	packages, err := websiteSubSDK.Packages(t.Context())
	assert.NoError(t, err)
	var pkg twipla3as.Package
	for _, p := range packages {
		if p.Touchpoints > 1000 {
			pkg = p
			break
		}
	}
	var superiorPkg twipla3as.Package
	for _, p := range packages {
		if p.Touchpoints > pkg.Touchpoints && p.Period == pkg.Period {
			superiorPkg = p
			break
		}
	}

	intpcName := fmt.Sprintf("go-sdk-intpc-%d", rand.Int())
	websiteName := fmt.Sprintf("go-sdk-website-%d", rand.Int())
	rndEmail := fmt.Sprintf("%d@twipla.com", rand.Int())
	rndDomain := fmt.Sprintf("%d.twiplatest.com", rand.Int())
	intpc, err := websiteSubSDK.CreateINTPC(t.Context(), twipla3as.CreateINTPCArgs{
		ExternalCustomerID: intpcName,
		Email:              rndEmail,
		SubscriptionType:   twipla3as.SubscriptionTypeWebsite,
		PackageID:          pkg.ID,
		BillingDate:        time.Now(),
		ExternalWebsiteID:  websiteName,
		Domain:             rndDomain,
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.NotNil(t, intpc)
	assert.NotEmpty(t, intpc.ID)
	assert.Equal(t, intpcName, intpc.IntpCustomerID)
	assert.Equal(t, rndEmail, intpc.Email)

	defer func() {
		_, err := websiteSubSDK.DeleteINTPC(t.Context(), intpcName)
		assert.NoError(t, err)
	}()

	t.Run("Upgrade", func(t *testing.T) {
		assert.NoError(t, websiteSubSDK.UpgradeWebsiteSubscription(t.Context(), twipla3as.UpgradeWebsiteSubscriptionArgs{
			WebsiteID: websiteName,
			PackageID: superiorPkg.ID,
			Trial:     false,
			Prorate:   false,
		}))
	})

	t.Run("Downgrade", func(t *testing.T) {
		// It was upgraded to a superior package before, so we can downgrade to the original package
		assert.NoError(t, websiteSubSDK.DowngradeWebsiteSubscription(t.Context(), twipla3as.DowngradeWebsiteSubscriptionArgs{
			WebsiteID: websiteName,
			PackageID: pkg.ID,
		}))
	})

	t.Run("Cancel", func(t *testing.T) {
		assert.NoError(t, websiteSubSDK.CancelWebsiteSubscription(t.Context(), twipla3as.CancelWebsiteSubscriptionArgs{
			WebsiteID: websiteName,
		}))
	})

	t.Run("Resume", func(t *testing.T) {
		assert.NoError(t, websiteSubSDK.ResumeWebsiteSubscription(t.Context(), twipla3as.ResumeWebsiteSubscriptionArgs{
			WebsiteID: websiteName,
		}))
	})

	t.Run("Deactivate", func(t *testing.T) {
		assert.NoError(t, websiteSubSDK.DeactivateWebsiteSubscription(t.Context(), twipla3as.DeactivateWebsiteSubscriptionArgs{
			WebsiteID: websiteName,
		}))
	})

	t.Run("Reactivate through upgrade", func(t *testing.T) {
		t.Skip("Broken in aaas-api. Issue raised")
		assert.NoError(t, websiteSubSDK.UpgradeWebsiteSubscription(t.Context(), twipla3as.UpgradeWebsiteSubscriptionArgs{
			WebsiteID: websiteName,
			PackageID: pkg.ID,
			Trial:     false,
			Prorate:   false,
		}))
	})
}
