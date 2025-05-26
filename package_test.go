package twipla3as_test

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	twipla3as "github.com/twipla/3as-go-sdk"
	"math/rand/v2"
	"testing"
)

func TestPackages(t *testing.T) {
	var packageID string
	var recommended bool
	t.Run("list", func(t *testing.T) {
		packages, err := mainSDK.Packages(t.Context())
		assert.NoError(t, err)
		assert.NotEmpty(t, packages)
		for _, p := range packages {
			if p.Touchpoints > 100 && !p.Recommended {
				packageID = p.ID
				break
			}
		}
		if packageID == "" {
			recommended = true
			for _, p := range packages {
				if p.Touchpoints > 100 {
					packageID = p.ID
					break
				}
			}
		}
		spew.Dump(packages)
	})

	t.Run("get", func(t *testing.T) {
		pkg, err := mainSDK.Package(t.Context(), packageID)
		assert.NoError(t, err)
		assert.Equal(t, packageID, pkg.ID)
	})

	t.Run("create", func(t *testing.T) {
		t.Skip("cannot remove packages after creation, skipping to not cause clutter")

		pkg, err := mainSDK.CreatePackage(t.Context(), twipla3as.CreatePackageArgs{
			Name:        "New Package",
			Touchpoints: 1234,
			Price:       123.45,
			Period:      twipla3as.PeriodMonthly,
			Currency:    twipla3as.CurrencyEUR,
		})
		assert.NoError(t, err)
		assert.NotNil(t, pkg)
		assert.NotEmpty(t, pkg.ID)
	})

	t.Run("update", func(t *testing.T) {
		if recommended == true {
			t.Skip("All packages are recommended, which can't be updated. Update test is skipped")
		}
		newName := fmt.Sprintf("New Package - %d", rand.Int())

		pkg, err := mainSDK.UpdatePackage(t.Context(), packageID, twipla3as.UpdatePackageArgs{
			Name: newName,
		})
		assert.NoError(t, err)
		assert.NotNil(t, pkg)
		assert.Equal(t, newName, pkg.Name)
		assert.Equal(t, packageID, pkg.ID)

		pkg2, err := mainSDK.Package(t.Context(), packageID)
		assert.NoError(t, err)
		assert.NotNil(t, pkg2)
		assert.Equal(t, newName, pkg2.Name)
	})
}
