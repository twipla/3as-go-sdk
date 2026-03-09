package twipla3as_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	twipla3as "github.com/twipla/3as-go-sdk"
)

func TestWebsiteContributors(t *testing.T) {
	owner := randINTPC(twipla3as.SubscriptionTypeWebsite)
	contributor := randINTPC(twipla3as.SubscriptionTypeWebsite)

	_, err := websiteSubSDK.CreateINTPC(t.Context(), owner)
	assert.NoError(t, err)
	t.Cleanup(func() {
		websiteSubSDK.DeleteINTPC(context.Background(), owner.ExternalCustomerID)
	})

	_, err = websiteSubSDK.CreateINTPC(t.Context(), contributor)
	assert.NoError(t, err)
	t.Cleanup(func() {
		websiteSubSDK.DeleteINTPC(context.Background(), contributor.ExternalCustomerID)
	})

	// Add the contributor to the owner's website with an Editor role.
	err = websiteSubSDK.AddWebsiteContributor(t.Context(), twipla3as.CreateWebsiteContributorArgs{
		ExternalWebsiteId:  owner.ExternalWebsiteID,
		ExternalCustomerId: contributor.ExternalCustomerID,
		Role:               twipla3as.ContributorRoleEditor,
	})
	assert.NoError(t, err)

	// Verify the contributor list reflects the newly added editor.
	contributorsList, err := websiteSubSDK.ListWebsiteContributors(t.Context(), twipla3as.ListWebsiteContributorsArgs{
		ExternalWebsiteId: owner.ExternalWebsiteID,
	})
	assert.NoError(t, err)

	assert.Equal(t, owner.ExternalCustomerID, contributorsList.Owner.ExternalCustomerId)
	assert.Equal(t, owner.Email, contributorsList.Owner.Email)

	assert.Equal(t, 1, len(contributorsList.Contributors[twipla3as.ContributorRoleEditor]))

	got := contributorsList.Contributors[twipla3as.ContributorRoleEditor][0]
	assert.Equal(t, contributor.ExternalCustomerID, got.ExternalCustomerId)
	assert.Equal(t, contributor.Email, got.Email)

	// Remove the contributor and confirm they no longer appear in the list.
	err = websiteSubSDK.DeleteWebsiteContributor(t.Context(), twipla3as.DeleteWebsiteContributorArgs{ExternalWebsiteId: owner.ExternalWebsiteID, ExternalCustomerId: contributor.ExternalCustomerID})
	assert.NoError(t, err)

	contributorsList, err = websiteSubSDK.ListWebsiteContributors(t.Context(), twipla3as.ListWebsiteContributorsArgs{
		ExternalWebsiteId: owner.ExternalWebsiteID,
	})
	assert.NoError(t, err)

	assert.Equal(t, 0, len(contributorsList.Contributors[twipla3as.ContributorRoleEditor]))
}
