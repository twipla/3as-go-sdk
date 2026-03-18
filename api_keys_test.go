package twipla3as_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	twipla3as "github.com/twipla/3as-go-sdk"
)

func TestApiKeys(t *testing.T) {
	sdk := websiteSubSDK

	intpc := randINTPC(twipla3as.SubscriptionTypeWebsite)
	_, err := sdk.CreateINTPC(t.Context(), intpc)
	assert.NoError(t, err)
	t.Cleanup(func() {
		websiteSubSDK.DeleteINTPC(context.Background(), intpc.ExternalCustomerID)
	})

	t.Run("api key can be created successfully", func(t *testing.T) {
		args := twipla3as.CreateApiKeyArgs{
			ExternalWebsiteID: intpc.ExternalWebsiteID,
			Name:              fmt.Sprintf("go-sdk-api-key-%d", rand.Int()),
		}
		res, err := sdk.CreateWebsiteApiKey(t.Context(), args)
		assert.NoError(t, err)

		assert.Equal(t, args.ExternalWebsiteID, res.IntpWebsiteId)
		assert.NotEmpty(t, res.ApiKey)
		assert.NotEmpty(t, res.Id)
		assert.NotEmpty(t, res.ExpiresAt)
		assert.NotEmpty(t, res.CreatedAt)
		err = sdk.DeleteWebsiteApiKey(t.Context(), intpc.ExternalWebsiteID, res.Id)
		require.NoError(t, err)
	})

	t.Run("list api keys for website should work", func(t *testing.T) {
		createdKeys := []string{}
		keysCount := 3

		for range keysCount {
			res, err := sdk.CreateWebsiteApiKey(t.Context(), twipla3as.CreateApiKeyArgs{
				ExternalWebsiteID: intpc.ExternalWebsiteID,
				Name:              fmt.Sprintf("go-sdk-api-key-%d", rand.Int()),
			})
			assert.NoError(t, err)
			createdKeys = append(createdKeys, res.Id)

		}
		defer func() {
			for _, id := range createdKeys {
				err := sdk.DeleteWebsiteApiKey(t.Context(), intpc.ExternalWebsiteID, id)
				assert.NoError(t, err)
			}
		}()

		res, err := sdk.ListWebsiteApiKeys(t.Context(), intpc.ExternalWebsiteID)
		assert.NoError(t, err)

		assert.Equal(t, len(res), keysCount)
		for _, item := range res {
			assert.Empty(t, item.ApiKey)

			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.ExpiresAt)
			assert.NotEmpty(t, item.CreatedAt)
			assert.Equal(t, intpc.ExternalWebsiteID, item.IntpWebsiteId)
		}
	})

}
