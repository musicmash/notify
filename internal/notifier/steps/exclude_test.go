package steps

import (
	"testing"
	"time"

	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/musicmash/notify/internal/db"
	"github.com/musicmash/notify/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNotifierSteps_Exclude(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.CreateNotification(&db.Notification{
		Date: time.Now(), UserName: testutil.UserBot, ReleaseID: 1,
	}))
	assert.NoError(t, db.DbMgr.CreateNotification(&db.Notification{
		Date: time.Now(), UserName: testutil.UserObjque, ReleaseID: 2,
	}))
	items := []*Item{{
		ArtistID:    testutil.StoreIDQ,
		Releases:    []*releases.Release{{ID: 1}},
		Subscribers: []string{testutil.UserObjque, testutil.UserBot},
	}}
	step := NewExcludeSubscribersStep()

	// action
	items, err := step.Do(time.Now(), items)

	// assert
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, int64(testutil.StoreIDQ), items[0].ArtistID)
	assert.Len(t, items[0].Subscribers, 1)
	assert.Equal(t, testutil.UserObjque, items[0].Subscribers[0])
	assert.Len(t, items[0].Releases, 1)
	assert.Equal(t, uint64(1), items[0].Releases[0].ID)
	assert.Empty(t, items[0].Chats)
	assert.Empty(t, items[0].ArtistName)
}
