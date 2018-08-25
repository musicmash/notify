package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDB_Releases_EnsureExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	err := DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
	})

	// assert
	assert.NoError(t, err)
	release, err := DbMgr.FindRelease("skrillex", 182821355)
	assert.NoError(t, err)
	assert.Equal(t, "skrillex", release.ArtistName)
	assert.Equal(t, uint64(182821355), release.StoreID)
}

func TestDB_Releases_IsExists(t *testing.T) {
	setup()
	defer teardown()

	// action
	const storeID = uint64(182821355)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    storeID,
	}))

	// action and assert
	assert.True(t, DbMgr.IsReleaseExists(storeID))
}

func TestDB_Releases_List(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
	}))
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreID:    213551828,
	}))

	// action
	releases, err := DbMgr.GetAllReleases()

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, "skrillex", releases[0].ArtistName)
	assert.Equal(t, "S.P.Y", releases[1].ArtistName)
}

func TestDB_Releases_GetReleasesForUserFilterByPeriod(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	const userID = "objque@me"
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
		Date:       time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreID:    213551828,
		Date:       time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserID:     userID,
		ArtistName: "skrillex",
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserID:     userID,
		ArtistName: "S.P.Y",
	}))
	// add stores
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore("itunes", "4nDoR", 1))
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore("yandex", "4nDoR", 1))
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore("itunes", "wC5Bh", 2))

	// action
	since := time.Now().UTC().Add(-time.Hour * 48)
	till := time.Now().UTC()
	releases, err := DbMgr.GetReleasesForUserFilterByPeriod(userID, since, till)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "skrillex", releases[0].ArtistName)
	assert.Len(t, releases[0].Stores, 2)
}

func TestDB_Releases_GetReleasesForUserFilterByPeriod_WithFuture(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	const userID = "objque@me"
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
		Date:       time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreID:    213551828,
		Date:       time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserID:     userID,
		ArtistName: "skrillex",
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserID:     userID,
		ArtistName: "S.P.Y",
	}))
	// add stores
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore("itunes", "4nDoR", 1))
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore("yandex", "4nDoR", 1))

	// action
	since := time.Now().UTC().Add(-time.Hour * 48)
	till := time.Now().UTC().Add(time.Hour * 48)
	releases, err := DbMgr.GetReleasesForUserFilterByPeriod(userID, since, till)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	assert.Equal(t, "skrillex", releases[0].ArtistName)
	assert.Len(t, releases[0].Stores, 2)
	assert.Len(t, releases[1].Stores, 0)
}

func TestDB_Releases_GetReleasesForUserSince(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// released release
	const userID = "objque@me"
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
		Date:       time.Now().UTC().Add(-time.Hour * 24),
	}))
	// announced release
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreID:    213551828,
		Date:       time.Now().UTC().Add(time.Hour * 24),
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserID:     userID,
		ArtistName: "skrillex",
	}))
	assert.NoError(t, DbMgr.EnsureSubscriptionExists(&Subscription{
		UserID:     userID,
		ArtistName: "S.P.Y",
	}))
	// add stores
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore("itunes", "4nDoR", 1))

	// action
	releases, err := DbMgr.GetReleasesForUserSince(userID, time.Now())

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "S.P.Y", releases[0].ArtistName)
	assert.Len(t, releases[0].Stores, 0)
}

func TestDB_Releases_FindReleasesWithFilter(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const userID = "objque@me"
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "skrillex",
		StoreID:    182821355,
		Date:       time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  time.Now().UTC().Add(-time.Hour * 48),
	}))
	date := time.Now().UTC().Add(-time.Hour * 24)
	assert.NoError(t, DbMgr.EnsureReleaseExists(&Release{
		ArtistName: "S.P.Y",
		StoreID:    213551828,
		Date:       time.Now().UTC().Add(time.Hour * 48),
		CreatedAt:  date,
	}))
	// add stores
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore("itunes", "4nDoR", 1))
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore("yandex", "4nDoR", 1))
	assert.NoError(t, DbMgr.EnsureReleaseExistsInStore("itunes", "wC5Bh", 2))

	// action
	releases, err := DbMgr.FindNewReleases(date)

	// assert
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "S.P.Y", releases[0].ArtistName)
	assert.Len(t, releases[0].Stores, 1)
	assert.Equal(t, releases[0].Stores[0].StoreType, "itunes")
	assert.Equal(t, releases[0].Stores[0].StoreID, "wC5Bh")
}
