package twipla3as_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	twipla3as "github.com/twipla/3as-go-sdk"
	"math/rand/v2"
	"testing"
	"time"
)

func TestWebsites(t *testing.T) {
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

	t.Run("List all websites", func(t *testing.T) {
		websites, pagination, err := mainSDK.Websites(t.Context(), twipla3as.Pagination{
			Page:     0,
			PageSize: 15,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, websites)
		assert.NotEmpty(t, pagination)
		assert.Equal(t, 15, pagination.PageSize)
	})

	secondWebsiteName := websiteName + "-2"
	t.Run("Create", func(t *testing.T) {
		args := twipla3as.CreateWebsiteArgs{
			ExternalID: secondWebsiteName,
			IntpcID:    intpcName,
			Domain:     "2-" + rndDomain,
			UFT:        false,
		}
		if subType == twipla3as.SubscriptionTypeWebsite {
			args.PackageID = pkg.ID
			args.BillingDate = time.Now()
		}
		err := mainSDK.CreateWebsite(t.Context(), args)
		assert.NoError(t, err)
	})

	t.Run("List INTPC websites", func(t *testing.T) {
		websites, pagination, err := mainSDK.IntpcWebsites(t.Context(), intpcName, twipla3as.Pagination{PageSize: 15})
		assert.NoError(t, err)
		assert.NotNil(t, websites)
		assert.NotNil(t, pagination)
		assert.NotEmpty(t, websites)
		assert.Equal(t, 15, pagination.PageSize)
		assert.Equal(t, 2, pagination.Total)
		if assert.Equal(t, 2, len(websites)) {
			assert.ElementsMatch(t, []string{websites[0].Domain, websites[1].Domain}, []string{rndDomain, "2-" + rndDomain})
		}
	})

	t.Run("Get one website", func(t *testing.T) {
		website, err := mainSDK.Website(t.Context(), secondWebsiteName)
		assert.NoError(t, err)
		assert.NotNil(t, website)
		assert.NotEmpty(t, website.ID)
		assert.Equal(t, website.ExternalWebsiteID, secondWebsiteName)
		assert.Equal(t, website.IntpCustomerID, intpcName)
	})

	t.Run("Delete", func(t *testing.T) {
		assert.NoError(t, mainSDK.DeleteWebsite(t.Context(), secondWebsiteName))
	})
}
