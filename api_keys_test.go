package twipla3as_test

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	twipla3as "github.com/twipla/3as-go-sdk"
)

func TestApiKeys(t *testing.T) {
	intpcName := fmt.Sprintf("go-sdk-intpc-%d", rand.Int())
	websiteId := fmt.Sprintf("go-sdk-website-%d", rand.Int())
	rndEmail := fmt.Sprintf("%d@twipla.com", rand.Int())
	rndDomain := fmt.Sprintf("%d.twiplatest.com", rand.Int())
	_, err := mainSDK.CreateINTPC(t.Context(), twipla3as.CreateINTPCArgs{
		ExternalCustomerID: intpcName,
		Email:              rndEmail,
		SubscriptionType:   twipla3as.SubscriptionTypeWebsite,
		BillingDate:        time.Now(),
		ExternalWebsiteID:  websiteId,
		Domain:             rndDomain,
	})
	assert.NoError(t, err)

	sdk := websiteSubSDK

	t.Run("api key can be created successfully", func(t *testing.T) {
		args := twipla3as.CreateApiKeyArgs{
			ExternalWebsiteID: websiteId,
			Name:              fmt.Sprintf("go-sdk-api-key-%d", rand.Int()),
		}
		res, err := sdk.CreateWebsiteApiKey(t.Context(), args)
		assert.NoError(t, err)

		assert.Equal(t, args.ExternalWebsiteID, res.IntpWebsiteId)
		assert.NotEmpty(t, res.ApiKey)
		assert.NotEmpty(t, res.Id)
		assert.NotEmpty(t, res.ExpiresAt)
		assert.NotEmpty(t, res.CreatedAt)
		err = sdk.DeleteWebsiteApiKey(t.Context(), websiteId, res.Id)
		require.NoError(t, err)
	})

	t.Run("list api keys for website should work", func(t *testing.T) {
		createdKeys := []string{}
		keysCount := 3

		for range keysCount {
			res, err := sdk.CreateWebsiteApiKey(t.Context(), twipla3as.CreateApiKeyArgs{
				ExternalWebsiteID: websiteId,
				Name:              fmt.Sprintf("go-sdk-api-key-%d", rand.Int()),
			})
			assert.NoError(t, err)
			createdKeys = append(createdKeys, res.Id)

		}
		defer func() {
			for _, id := range createdKeys {
				err := sdk.DeleteWebsiteApiKey(t.Context(), websiteId, id)
				assert.NoError(t, err)
			}
		}()

		res, err := sdk.ListWebsiteApiKeys(t.Context(), websiteId)
		assert.NoError(t, err)

		assert.Equal(t, len(res), keysCount)
		for _, item := range res {
			assert.Empty(t, item.ApiKey)

			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.ExpiresAt)
			assert.NotEmpty(t, item.CreatedAt)
			assert.Equal(t, websiteId, item.IntpWebsiteId)
		}
	})

}
